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

func (a *ApiHandler) ListAnimalTypes(ctx context.Context, request api.ListAnimalTypesRequestObject) (api.ListAnimalTypesResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalTypes",
			Object:    "animal_types",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAnimalTypes403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_type_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalTypesFilters{ListAnimalTypesParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	animalTypes, err := a.persistor.Animal().ListAnimalTypes(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animalTypes: %w", err)
	}

	animalTypesData := make([]api.AnimalType, len(animalTypes.Data))
	for i, animalType := range animalTypes.Data {
		animalTypesData[i] = dto.AnimalTypeToResponse(animalType)
	}

	resp := api.ListAnimalTypes200JSONResponse{
		Data: animalTypesData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animalTypes.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetAnimalType(ctx context.Context, request api.GetAnimalTypeRequestObject) (api.GetAnimalTypeResponseObject, error) {
	animalType, err := a.persistor.Animal().GetAnimalTypeByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalTypeNotFound) {
			return api.GetAnimalType404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_type_not_found", "animalType not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an animalType by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalType",
			Object:    shared.AuthzAnimalTypeID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnimalType403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_type_permission", "permission denied")}, nil
	}

	resp := api.GetAnimalType200JSONResponse(dto.AnimalTypeToResponse(animalType))

	return resp, nil
}

func (a *ApiHandler) CreateAnimalType(ctx context.Context, request api.CreateAnimalTypeRequestObject) (api.CreateAnimalTypeResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalTypes",
			Object:    "animal_types",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateAnimalType403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_type_permission", "permission denied")}, nil
	}

	animalTypeSetter := models.AnimalTypeSetter{
		Name: omit.From(request.Body.Name),
	}

	animalType, err := a.persistor.Animal().CreateAnimalType(ctx, animalTypeSetter)
	if err != nil {
		msg := "could not create an animal type"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateAnimalType400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_type_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animalType integrity error"
			return api.CreateAnimalType400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_type_save", msg, reason)}, nil
		}

		return api.CreateAnimalTypedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_type_save", msg, reason)}, nil
	}

	resp := api.CreateAnimalType201JSONResponse(dto.AnimalTypeToResponse(animalType))

	// @TODO: use outbox pattern
	if err := createAnimalTypeRelationTuples(ctx, a.Keto, resp.ID); err != nil {
		a.Log.Error("failed to insert animalType relation-tuple", slog.Int64("animal_type_id", resp.ID), slog.Any("error", err))
		return api.CreateAnimalTypedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_type_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) UpdateAnimalType(ctx context.Context, request api.UpdateAnimalTypeRequestObject) (api.UpdateAnimalTypeResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalType",
			Object:    shared.AuthzAnimalTypeID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimalType403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_type_permission", "permission denied")}, nil
	}

	animalType, err := a.persistor.Animal().UpdateAnimalType(ctx, request.ID, models.AnimalTypeSetter{
		Name: omit.FromPtr(request.Body.Name),
	})
	if err != nil {
		msg := "could not update an animal type"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimalType400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_type_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animalType integrity error"
			return api.UpdateAnimalType400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_type_update", msg, reason)}, nil
		}

		return api.UpdateAnimalTypedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_type_update", msg, reason)}, nil
	}

	resp := api.UpdateAnimalType201JSONResponse(dto.AnimalTypeToResponse(animalType))

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalType(ctx context.Context, request api.DeleteAnimalTypeRequestObject) (api.DeleteAnimalTypeResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalType",
			Object:    shared.AuthzAnimalTypeID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalType403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_type_permission", "permission denied")}, nil
	}

	_, err := a.persistor.Animal().DeleteAnimalTypeByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an animalType by id: %w", err)
	}

	resp := api.DeleteAnimalType204Response{}

	// @TODO: use outbox pattern
	if err := deleteAnimalTypeRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete animalType relation-tuple", slog.Int64("animal_type_id", request.ID), slog.Any("error", err))
		return api.DeleteAnimalTypedefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_type_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalTypes(ctx context.Context, request api.DeleteAnimalTypesRequestObject) (api.DeleteAnimalTypesResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimalTypes204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "AnimalType",
			Object:    shared.AuthzAnimalTypeID(request.Body.Ids[0]), // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalTypes403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animalTypes_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().DeleteAnimalTypes(ctx, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete an animalTypes by ids: %w", err)
	}

	resp := api.DeleteAnimalTypes204Response{}

	// @TODO: use outbox pattern
	for _, id := range request.Body.Ids {
		if err := deleteAnimalTypeRelationTuples(ctx, a.Keto, id); err != nil {
			a.Log.Error("failed to delete animalTypes relation-tuple", slog.Int64("animal_type_id", id), slog.Any("error", err))
			return api.DeleteAnimalTypesdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animalTypes_permissions", "failed to delete permissions")}, nil
		}
	}

	return resp, nil
}

func createAnimalTypeRelationTuples(ctx context.Context, c *keto.Client, animalTypeID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "AnimalType",
					Object:    shared.AuthzAnimalTypeID(animalTypeID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("AnimalTypes", "animal_types", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert animalType relation tuples: %w", err)
	}

	return nil
}

func deleteAnimalTypeRelationTuples(ctx context.Context, c *keto.Client, animalTypeID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: "AnimalType",
					Object:    shared.AuthzAnimalTypeID(animalTypeID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("AnimalTypes", "animal_types", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to delete animalType relation tuples: %w", err)
	}

	return nil
}
