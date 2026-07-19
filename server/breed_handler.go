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

func (a *ApiHandler) ListBreeds(ctx context.Context, request api.ListBreedsRequestObject) (api.ListBreedsResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Breeds",
			Object:    "breeds",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListBreeds403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("breed_permission", "permission denied")}, nil
	}

	filters := dbtype.ListBreedsFilters{ListBreedsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	breeds, err := a.persistor.Animal().ListBreeds(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list breeds: %w", err)
	}

	breedsData := make([]api.Breed, len(breeds.Data))
	for i, breed := range breeds.Data {
		breedsData[i] = dto.BreedToResponse(breed)
	}

	resp := api.ListBreeds200JSONResponse{
		Data: breedsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, breeds.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetBreed(ctx context.Context, request api.GetBreedRequestObject) (api.GetBreedResponseObject, error) {
	breed, err := a.persistor.Animal().GetBreedByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrBreedNotFound) {
			return api.GetBreed404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("breed_not_found", "breed not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an breed by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Breed",
			Object:    shared.AuthzBreedID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetBreed403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("breed_permission", "permission denied")}, nil
	}

	resp := api.GetBreed200JSONResponse(dto.BreedToResponse(breed))

	return resp, nil
}

func (a *ApiHandler) CreateBreed(ctx context.Context, request api.CreateBreedRequestObject) (api.CreateBreedResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Breeds",
			Object:    "breeds",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateBreed403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("breed_permission", "permission denied")}, nil
	}

	breedSetter := models.BreedSetter{
		AnimalSpecieID: omit.From(request.Body.AnimalSpecieID),
		Name:           omit.From(request.Body.Name),
	}

	breed, err := a.persistor.Animal().CreateBreed(ctx, breedSetter)
	if err != nil {
		msg := "could not create an animal type"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateBreed400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "breed_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "breed integrity error"
			return api.CreateBreed400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "breed_save", msg, reason)}, nil
		}

		return api.CreateBreeddefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "breed_save", msg, reason)}, nil
	}

	resp := api.CreateBreed201JSONResponse(dto.BreedToResponse(breed))

	// @TODO: use outbox pattern
	if err := createBreedRelationTuples(ctx, a.Keto, resp.ID); err != nil {
		a.Log.Error("failed to insert breed relation-tuple", slog.Int64("breed_id", resp.ID), slog.Any("error", err))
		return api.CreateBreeddefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "breed_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) UpdateBreed(ctx context.Context, request api.UpdateBreedRequestObject) (api.UpdateBreedResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Breed",
			Object:    shared.AuthzBreedID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateBreed403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("breed_permission", "permission denied")}, nil
	}

	breedSetter := models.BreedSetter{
		Name: omit.FromPtr(request.Body.Name),
	}

	breed, err := a.persistor.Animal().UpdateBreed(ctx, request.ID, breedSetter)
	if err != nil {
		msg := "could not update an animal type"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateBreed400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "breed_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "breed integrity error"
			return api.UpdateBreed400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "breed_update", msg, reason)}, nil
		}

		return api.UpdateBreeddefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "breed_update", msg, reason)}, nil
	}

	resp := api.UpdateBreed201JSONResponse(dto.BreedToResponse(breed))

	return resp, nil
}

func (a *ApiHandler) DeleteBreed(ctx context.Context, request api.DeleteBreedRequestObject) (api.DeleteBreedResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Breed",
			Object:    shared.AuthzBreedID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteBreed403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("breed_permission", "permission denied")}, nil
	}

	_, err := a.persistor.Animal().DeleteBreedByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an breed by id: %w", err)
	}

	resp := api.DeleteBreed204Response{}

	// @TODO: use outbox pattern
	if err := deleteBreedRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete breed relation-tuple", slog.Int64("breed_id", request.ID), slog.Any("error", err))
		return api.DeleteBreeddefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "breed_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) DeleteBreeds(ctx context.Context, request api.DeleteBreedsRequestObject) (api.DeleteBreedsResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteBreeds204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Breed",
			Object:    shared.AuthzBreedID(request.Body.Ids[0]), // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteBreeds403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("breeds_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().DeleteBreeds(ctx, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete an breeds by ids: %w", err)
	}

	resp := api.DeleteBreeds204Response{}

	// @TODO: use outbox pattern
	for _, id := range request.Body.Ids {
		if err := deleteBreedRelationTuples(ctx, a.Keto, id); err != nil {
			a.Log.Error("failed to delete breeds relation-tuple", slog.Int64("breed_id", id), slog.Any("error", err))
			return api.DeleteBreedsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "breeds_permissions", "failed to delete permissions")}, nil
		}
	}

	return resp, nil
}

func createBreedRelationTuples(ctx context.Context, c *keto.Client, breedID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Breed",
					Object:    shared.AuthzBreedID(breedID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Breeds", "breeds", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert breed relation tuples: %w", err)
	}

	return nil
}

func deleteBreedRelationTuples(ctx context.Context, c *keto.Client, breedID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Breed",
					Object:    shared.AuthzBreedID(breedID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Breeds", "breeds", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to delete breed relation tuples: %w", err)
	}

	return nil
}
