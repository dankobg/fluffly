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

func (a *ApiHandler) ListAnimalSpecies(ctx context.Context, request api.ListAnimalSpeciesRequestObject) (api.ListAnimalSpeciesResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalSpecies",
			Object:    "animal_species",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAnimalSpecies403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_specie_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalSpeciesFilters{ListAnimalSpeciesParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	animalSpecies, err := a.persistor.Animal().ListAnimalSpecies(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animalSpecies: %w", err)
	}

	animalSpeciesData := make([]api.AnimalSpecie, len(animalSpecies.Data))
	for i, animalSpecie := range animalSpecies.Data {
		animalSpeciesData[i] = dto.AnimalSpeciesToResponse(animalSpecie)
	}

	resp := api.ListAnimalSpecies200JSONResponse{
		Data: animalSpeciesData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animalSpecies.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetAnimalSpecie(ctx context.Context, request api.GetAnimalSpecieRequestObject) (api.GetAnimalSpecieResponseObject, error) {
	animalSpecie, err := a.persistor.Animal().GetAnimalSpecieByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalSpeciesNotFound) {
			return api.GetAnimalSpecie404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_specie_not_found", "animalSpecie not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an animalSpecie by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalSpecie",
			Object:    shared.AuthzAnimalSpecieID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnimalSpecie403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_specie_permission", "permission denied")}, nil
	}

	resp := api.GetAnimalSpecie200JSONResponse(dto.AnimalSpeciesToResponse(animalSpecie))

	return resp, nil
}

func (a *ApiHandler) CreateAnimalSpecie(ctx context.Context, request api.CreateAnimalSpecieRequestObject) (api.CreateAnimalSpecieResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalSpecies",
			Object:    "animal_species",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateAnimalSpecie403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_specie_permission", "permission denied")}, nil
	}

	animalSpecieSetter := models.AnimalSpecieSetter{
		AnimalTypeID: omit.From(request.Body.AnimalTypeID),
		Name:         omit.From(request.Body.Name),
	}

	animalSpecie, err := a.persistor.Animal().CreateAnimalSpecie(ctx, animalSpecieSetter)
	if err != nil {
		msg := "could not create an animal species"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateAnimalSpecie400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_specie_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animalSpecie integrity error"
			return api.CreateAnimalSpecie400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_specie_save", msg, reason)}, nil
		}

		return api.CreateAnimalSpeciedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_specie_save", msg, reason)}, nil
	}

	resp := api.CreateAnimalSpecie201JSONResponse(dto.AnimalSpeciesToResponse(animalSpecie))

	// @TODO: use outbox pattern
	if err := createAnimalSpecieRelationTuples(ctx, a.Keto, resp.ID); err != nil {
		a.Log.Error("failed to insert animalSpecie relation-tuple", slog.Int64("animal_specie_id", resp.ID), slog.Any("error", err))
		return api.CreateAnimalSpeciedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_specie_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) UpdateAnimalSpecie(ctx context.Context, request api.UpdateAnimalSpecieRequestObject) (api.UpdateAnimalSpecieResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalSpecie",
			Object:    shared.AuthzAnimalSpecieID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimalSpecie403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_specie_permission", "permission denied")}, nil
	}

	animalSpecieSetter := models.AnimalSpecieSetter{
		Name: omit.FromPtr(request.Body.Name),
	}

	animalSpecie, err := a.persistor.Animal().UpdateAnimalSpecie(ctx, request.ID, animalSpecieSetter)
	if err != nil {
		msg := "could not update an animal species"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimalSpecie400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_specie_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animalSpecie integrity error"
			return api.UpdateAnimalSpecie400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_specie_update", msg, reason)}, nil
		}

		return api.UpdateAnimalSpeciedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_specie_update", msg, reason)}, nil
	}

	resp := api.UpdateAnimalSpecie201JSONResponse(dto.AnimalSpeciesToResponse(animalSpecie))

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalSpecie(ctx context.Context, request api.DeleteAnimalSpecieRequestObject) (api.DeleteAnimalSpecieResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalSpecie",
			Object:    shared.AuthzAnimalSpecieID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalSpecie403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_specie_permission", "permission denied")}, nil
	}

	_, err := a.persistor.Animal().DeleteAnimalSpecieByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an animalSpecie by id: %w", err)
	}

	resp := api.DeleteAnimalSpecie204Response{}

	// @TODO: use outbox pattern
	if err := deleteAnimalSpecieRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete animalSpecie relation-tuple", slog.Int64("animal_specie_id", request.ID), slog.Any("error", err))
		return api.DeleteAnimalSpeciedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_specie_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalSpecies(ctx context.Context, request api.DeleteAnimalSpeciesRequestObject) (api.DeleteAnimalSpeciesResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimalSpecies204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalSpecie",
			Object:    shared.AuthzAnimalSpecieID(request.Body.Ids[0]), // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalSpecies403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animalSpecies_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().DeleteAnimalSpecies(ctx, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete an animalSpecies by ids: %w", err)
	}

	resp := api.DeleteAnimalSpecies204Response{}

	// @TODO: use outbox pattern
	for _, id := range request.Body.Ids {
		if err := deleteAnimalSpecieRelationTuples(ctx, a.Keto, id); err != nil {
			a.Log.Error("failed to delete animalSpecies relation-tuple", slog.Int64("animal_specie_id", id), slog.Any("error", err))
			return api.DeleteAnimalSpeciesdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animalSpecies_permissions", "failed to delete permissions")}, nil
		}
	}

	return resp, nil
}

func createAnimalSpecieRelationTuples(ctx context.Context, c *keto.Client, animalSpecieID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "AnimalSpecie",
					Object:    shared.AuthzAnimalSpecieID(animalSpecieID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("AnimalSpecies", "animal_species", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert animalSpecie relation tuples: %w", err)
	}

	return nil
}

func deleteAnimalSpecieRelationTuples(ctx context.Context, c *keto.Client, animalSpecieID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: "AnimalSpecie",
					Object:    shared.AuthzAnimalSpecieID(animalSpecieID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("AnimalSpecies", "animal_species", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to delete animalSpecie relation tuples: %w", err)
	}

	return nil
}
