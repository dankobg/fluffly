package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/ptr"
	"github.com/oapi-codegen/nullable"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) CreateCountry(ctx context.Context, request api.CreateCountryRequestObject) (api.CreateCountryResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Countries",
			Object:    "countries",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.CreateCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "invalid permission")}, nil
	}

	countrySetter := dbtype.CountrySetter{
		Name:       nullable.NewNullableWithValue(request.Body.Name),
		IsoAlpha2:  nullable.NewNullableWithValue(request.Body.IsoAlpha2),
		IsoAlpha3:  nullable.NewNullableWithValue(request.Body.IsoAlpha3),
		IsoNumeric: nullable.NewNullableWithValue(request.Body.IsoNumeric),
	}
	country, err := a.persistor.Country().CreateCountry(ctx, countrySetter)
	if err != nil {
		msg := "could not create a country"
		var reason string
		var e1 postgres.ErrCountryUniqueViolation
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
			Object:    AuthzCountryID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.UpdateCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "invalid permission")}, nil
	}

	countrySetter := dbtype.CountrySetter{}
	if request.Body.Name != nil {
		countrySetter.Name = nullable.NewNullableWithValue(*request.Body.Name)
	}
	if request.Body.IsoAlpha2 != nil {
		countrySetter.IsoAlpha2 = nullable.NewNullableWithValue(*request.Body.IsoAlpha2)
	}
	if request.Body.IsoAlpha3 != nil {
		countrySetter.IsoAlpha3 = nullable.NewNullableWithValue(*request.Body.IsoAlpha3)
	}
	if request.Body.IsoNumeric != nil {
		countrySetter.IsoNumeric = nullable.NewNullableWithValue(*request.Body.IsoNumeric)
	}
	country, err := a.persistor.Country().UpdateCountry(ctx, request.ID, countrySetter)
	if err != nil {
		msg := "could not update a country"
		var reason string
		var e1 postgres.ErrCountryUniqueViolation
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
			Object:    AuthzCountryID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DeleteCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "invalid permission")}, nil
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

func (a *ApiHandler) ListCountries(ctx context.Context, request api.ListCountriesRequestObject) (api.ListCountriesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Countries",
			Object:    "countries",
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ListCountries403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "invalid permission")}, nil
	}

	var filters dbtype.CountryFilters
	filters.Pagination = ptr.Of(getPaginationParams(request.Params.Page, request.Params.PageSize))
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
			Object:    AuthzCountryID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetCountry403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("country_permission", "invalid permission")}, nil
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
					Object:    AuthzCountryID(countryID),
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
					Object:    AuthzCountryID(countryID),
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
