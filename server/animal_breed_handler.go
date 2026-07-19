package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/convert"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListAnimalBreeds(ctx context.Context, request api.ListAnimalBreedsRequestObject) (api.ListAnimalBreedsResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		slog.Error("list animal breeds permission denied", slog.Any("error", err))
		return api.ListAnimalBreeds403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_breeds_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalBreedsFilters{ListAnimalBreedsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	animalBreeds, err := a.persistor.Animal().ListAnimalBreeds(ctx, request.ID, filters)
	if err != nil {
		slog.Error("failed to list animal breeds", slog.Any("error", err))
		return nil, fmt.Errorf("failed to list animal breeds: %w", err)
	}

	animalBreedsData := make([]api.AnimalBreed, len(animalBreeds.Data))
	for i, breed := range animalBreeds.Data {
		animalBreedsData[i] = dto.AnimalBreedWithJoinDataToResponse(breed)
	}

	resp := api.ListAnimalBreeds200JSONResponse{
		Data: animalBreedsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animalBreeds.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetAnimalBreed(ctx context.Context, request api.GetAnimalBreedRequestObject) (api.GetAnimalBreedResponseObject, error) {
	animalBreed, err := a.persistor.Animal().GetAnimalBreed(ctx, request.ID, request.BreedID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalBreedNotFound) {
			return api.GetAnimalBreed404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_breed_not_found", "animal breed not found")}, nil
		}

		slog.Error("failed to get animal breed", slog.Int64("id", request.ID), slog.Int64("breed_id", request.BreedID), slog.Any("error", err))

		return nil, fmt.Errorf("failed to get an animal breed by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		slog.Error("get animal breed permission", slog.Int64("id", request.ID), slog.Int64("breed_id", request.BreedID), slog.Any("error", err))
		return api.GetAnimalBreeddefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "animal_breed_permission", "permission denied")}, nil
	}

	resp := api.GetAnimalBreed200JSONResponse(dto.AnimalBreedWithJoinDataToResponse(animalBreed))

	return resp, nil
}

func (a *ApiHandler) CreateAnimalBreeds(ctx context.Context, request api.CreateAnimalBreedsRequestObject) (api.CreateAnimalBreedsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		slog.Error("create animal breed permission", slog.Int64("id", request.ID), slog.Any("error", err))
		return api.CreateAnimalBreeds403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_breeds_permission", "permission denied")}, nil
	}

	animalCreateBreedSetters := make([]models.AnimalBreedSetter, len(request.Body.AnimalBreeds))
	for i, breed := range request.Body.AnimalBreeds {
		animalCreateBreedSetters[i] = models.AnimalBreedSetter{
			AnimalID: omit.From(request.ID),
			BreedID:  omit.From(breed.BreedID),
			Primary:  omit.FromPtr(breed.Primary),
		}
	}

	animalBreeds, err := a.persistor.Animal().CreateAnimalBreeds(ctx, request.ID, animalCreateBreedSetters)
	if err != nil {
		msg := "could not create animal breeds"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			slog.Error("create animal breeds", slog.String("reason", reason), slog.Int64("id", request.ID), slog.Any("error", err))

			return api.CreateAnimalBreeds400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_breeds_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			slog.Error("create animal breeds", slog.String("reason", reason), slog.Int64("id", request.ID), slog.Any("error", err))

			return api.CreateAnimalBreeds400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_breeds_save", msg, reason)}, nil
		}

		return api.CreateAnimalBreedsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_breeds_save", msg, reason)}, nil
	}

	var data []api.AnimalBreed
	for _, s := range animalBreeds {
		data = append(data, dto.AnimalBreedToResponse(s))
	}

	resp := api.CreateAnimalBreeds200JSONResponse(api.CreateAnimalBreeds200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateAnimalBreedsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_breeds_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalBreeds(ctx context.Context, request api.DeleteAnimalBreedsRequestObject) (api.DeleteAnimalBreedsResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimalBreeds204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalBreeds403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_breeds_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().DeleteAnimalBreeds(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a animal breeds by ids: %w", err)
	}

	resp := api.DeleteAnimalBreeds204Response{}

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete animal breeds relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 		return api.DeleteAnimalBreedsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_breeds_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateAnimalBreed(ctx context.Context, request api.UpdateAnimalBreedRequestObject) (api.UpdateAnimalBreedResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimalBreed403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_breed_permission", "permission denied")}, nil
	}

	animalBreedSetter := models.AnimalBreedSetter{
		BreedID: omit.FromPtr(request.Body.BreedID),
		Primary: convert.NullableToOmit(request.Body.Primary),
	}

	animalBreed, err := a.persistor.Animal().UpdateAnimalBreed(ctx, request.ID, request.BreedID, animalBreedSetter)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalBreedNotFound) {
			return api.UpdateAnimalBreed404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_breed_not_found", "animal breed not found")}, nil
		}

		msg := "could not update a breed"

		var (
			reason string
			e1     postgres.ErrCountryUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimalBreed400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_breed_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "breed integrity error"
			return api.UpdateAnimalBreed400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_breed_update", msg, reason)}, nil
		}

		return api.UpdateAnimalBreeddefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_breed_update", msg, reason)}, nil
	}

	resp := api.UpdateAnimalBreed200JSONResponse(dto.AnimalBreedToResponse(animalBreed))

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalBreed(ctx context.Context, request api.DeleteAnimalBreedRequestObject) (api.DeleteAnimalBreedResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalBreed403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_breed_permission", "permission denied")}, nil
	}

	if _, err := a.persistor.Animal().DeleteAnimalBreed(ctx, request.ID, request.BreedID); err != nil {
		return nil, fmt.Errorf("failed to delete a animal breed: %w", err)
	}

	resp := api.DeleteAnimalBreed204Response{}

	// @TODO: use outbox pattern
	// if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete animal breeds relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 	return api.DeleteAnimalBreedsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_breeds_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
