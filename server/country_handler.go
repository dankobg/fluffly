package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) CreateCountry(ctx context.Context, request api.CreateCountryRequestObject) (api.CreateCountryResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Countries",
			Object:    "countries",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "permission denied")}, nil
	}

	countrySetter := models.CountrySetter{
		Name:       omit.From(request.Body.Name),
		IsoAlpha2:  omit.From(request.Body.IsoAlpha2),
		IsoAlpha3:  omit.From(request.Body.IsoAlpha3),
		IsoNumeric: omit.From(request.Body.IsoNumeric),
	}

	country, err := a.persistor.Country().CreateCountry(ctx, countrySetter)
	if err != nil {
		msg := "could not create a country"

		var (
			reason string
			e1     postgres.ErrCountryUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateCountry400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "country_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "country integrity error"
			return api.CreateCountry400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "country_save", msg, reason)}, nil
		}

		return api.CreateCountrydefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "country_save", msg, reason)}, nil
	}

	resp := api.CreateCountry201JSONResponse(dto.CountryToResponse(country))

	// @TODO: use outbox pattern
	if err := createCountryRelationTuples(ctx, a.Keto, resp.ID); err != nil {
		a.Log.Error("failed to insert country relation-tuple", slog.Int64("country_id", resp.ID), slog.Any("error", err))
		return api.CreateCountrydefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "country_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) UpdateCountry(ctx context.Context, request api.UpdateCountryRequestObject) (api.UpdateCountryResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Country",
			Object:    shared.AuthzCountryID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "permission denied")}, nil
	}

	countrySetter := models.CountrySetter{
		Name:       omit.FromPtr(request.Body.Name),
		IsoAlpha2:  omit.FromPtr(request.Body.IsoAlpha2),
		IsoAlpha3:  omit.FromPtr(request.Body.IsoAlpha3),
		IsoNumeric: omit.FromPtr(request.Body.IsoNumeric),
	}

	country, err := a.persistor.Country().UpdateCountry(ctx, request.ID, countrySetter)
	if err != nil {
		msg := "could not update a country"

		var (
			reason string
			e1     postgres.ErrCountryUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateCountry400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "country_edit", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "country integrity error"
			return api.UpdateCountry400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "country_edit", msg, reason)}, nil
		}

		return api.UpdateCountrydefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "country_edit", msg, reason)}, nil
	}

	resp := api.UpdateCountry201JSONResponse(dto.CountryToResponse(country))

	return resp, nil
}

func (a *ApiHandler) DeleteCountry(ctx context.Context, request api.DeleteCountryRequestObject) (api.DeleteCountryResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Country",
			Object:    shared.AuthzCountryID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "permission denied")}, nil
	}

	_, err := a.persistor.Country().DeleteCountryByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete a country by id: %w", err)
	}

	resp := api.DeleteCountry204Response{}

	// @TODO: use outbox pattern
	if err := deleteCountryRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete country relation-tuple", slog.Int64("country_id", request.ID), slog.Any("error", err))
		return api.DeleteCountrydefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "country_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) DeleteCountries(ctx context.Context, request api.DeleteCountriesRequestObject) (api.DeleteCountriesResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteCountries204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Country",
			Object:    shared.AuthzCountryID(request.Body.Ids[0]), // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteCountries403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("countries_permission", "permission denied")}, nil
	}

	if err := a.persistor.Country().DeleteCountries(ctx, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a countries by ids: %w", err)
	}

	resp := api.DeleteCountries204Response{}

	// @TODO: use outbox pattern
	for _, id := range request.Body.Ids {
		if err := deleteCountryRelationTuples(ctx, a.Keto, id); err != nil {
			a.Log.Error("failed to delete countries relation-tuple", slog.Int64("country_id", id), slog.Any("error", err))
			return api.DeleteCountriesdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "countries_permissions", "failed to delete permissions")}, nil
		}
	}

	return resp, nil
}

func (a *ApiHandler) ListCountries(ctx context.Context, request api.ListCountriesRequestObject) (api.ListCountriesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Countries",
			Object:    "countries",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListCountries403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "permission denied")}, nil
	}

	filters := dbtype.ListCountriesFilters{ListCountriesParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	countries, err := a.persistor.Country().ListCountries(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list countries: %w", err)
	}

	countriesData := make([]api.Country, len(countries.Data))
	for i, country := range countries.Data {
		countriesData[i] = dto.CountryToResponse(country)
	}

	resp := api.ListCountries200JSONResponse{
		Data: countriesData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, countries.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetCountry(ctx context.Context, request api.GetCountryRequestObject) (api.GetCountryResponseObject, error) {
	sess := GetSession(ctx)

	country, err := a.persistor.Country().GetCountryByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrCountryNotFound) {
			return api.GetCountry404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("country_not_found", "country not found")}, nil
		}

		return nil, fmt.Errorf("failed to get a country by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Country",
			Object:    shared.AuthzCountryID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "permission denied")}, nil
	}

	resp := api.GetCountry200JSONResponse(dto.CountryToResponse(country))

	return resp, nil
}

func createCountryRelationTuples(ctx context.Context, c *keto.Client, countryID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Country",
					Object:    shared.AuthzCountryID(countryID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Countries", "countries", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert country relation tuples: %w", err)
	}

	return nil
}

func deleteCountryRelationTuples(ctx context.Context, c *keto.Client, countryID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Country",
					Object:    shared.AuthzCountryID(countryID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Countries", "countries", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to delete country relation tuples: %w", err)
	}

	return nil
}
