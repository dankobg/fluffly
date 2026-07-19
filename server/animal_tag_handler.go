package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListAnimalTags(ctx context.Context, request api.ListAnimalTagsRequestObject) (api.ListAnimalTagsResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAnimalTags403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_tags_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalTagsFilters{ListAnimalTagsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	orgTags, err := a.persistor.Animal().ListAnimalTags(ctx, request.ID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animal tags: %w", err)
	}

	orgTagsData := make([]api.AnimalTag, len(orgTags.Data))
	for i, tag := range orgTags.Data {
		orgTagsData[i] = dto.AnimalTagToResponse(tag)
	}

	resp := api.ListAnimalTags200JSONResponse{
		Data: orgTagsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, orgTags.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetAnimalTag(ctx context.Context, request api.GetAnimalTagRequestObject) (api.GetAnimalTagResponseObject, error) {
	orgTag, err := a.persistor.Animal().GetAnimalTag(ctx, request.ID, request.TagID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalTagNotFound) {
			return api.GetAnimalTag404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_tag_not_found", "animal tag not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an animal tag by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnimalTagdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "animal_tag_permission", "permission denied")}, nil
	}

	resp := api.GetAnimalTag200JSONResponse(dto.AnimalTagToResponse(orgTag))

	return resp, nil
}

func (a *ApiHandler) CreateAnimalTags(ctx context.Context, request api.CreateAnimalTagsRequestObject) (api.CreateAnimalTagsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateAnimalTags403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_tags_permission", "permission denied")}, nil
	}

	animalCreateTagSetters := make([]models.AnimalTagSetter, len(request.Body.Tags))
	for i, tag := range request.Body.Tags {
		animalCreateTagSetters[i] = models.AnimalTagSetter{
			Name:     omit.FromPtr(tag.Name),
			AnimalID: omitnull.From(request.ID),
		}
	}

	orgTags, err := a.persistor.Animal().CreateAnimalTags(ctx, request.ID, animalCreateTagSetters)
	if err != nil {
		msg := "could not create animal tags"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateAnimalTags400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_tags_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.CreateAnimalTags400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_tags_save", msg, reason)}, nil
		}

		return api.CreateAnimalTagsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_tags_save", msg, reason)}, nil
	}

	var data []api.AnimalTag
	for _, s := range orgTags {
		data = append(data, dto.AnimalTagToResponse(s))
	}

	resp := api.CreateAnimalTags200JSONResponse(api.CreateAnimalTags200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateAnimalTagsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_tags_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalTags(ctx context.Context, request api.DeleteAnimalTagsRequestObject) (api.DeleteAnimalTagsResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimalTags204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalTags403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_tags_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().DeleteAnimalTags(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a animal tags by ids: %w", err)
	}

	resp := api.DeleteAnimalTags204Response{}

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete animal tags relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 		return api.DeleteAnimalTagsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_tags_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateAnimalTag(ctx context.Context, request api.UpdateAnimalTagRequestObject) (api.UpdateAnimalTagResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimalTag403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_tag_permission", "permission denied")}, nil
	}

	orgTag, err := a.persistor.Animal().UpdateAnimalTag(ctx, request.ID, request.TagID, models.AnimalTagSetter{
		Name: omit.FromPtr(request.Body.Name),
	})
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalTagNotFound) {
			return api.UpdateAnimalTag404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_tag_not_found", "animal tag not found")}, nil
		}

		msg := "could not update a tag"

		var (
			reason string
			e1     postgres.ErrCountryUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimalTag400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_tag_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "tag integrity error"
			return api.UpdateAnimalTag400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_tag_update", msg, reason)}, nil
		}

		return api.UpdateAnimalTagdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_tag_update", msg, reason)}, nil
	}

	resp := api.UpdateAnimalTag200JSONResponse(dto.AnimalTagToResponse(orgTag))

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalTag(ctx context.Context, request api.DeleteAnimalTagRequestObject) (api.DeleteAnimalTagResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalTag403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_tag_permission", "permission denied")}, nil
	}

	if _, err := a.persistor.Animal().DeleteAnimalTag(ctx, request.ID, request.TagID); err != nil {
		return nil, fmt.Errorf("failed to delete a animal tag: %w", err)
	}

	resp := api.DeleteAnimalTag204Response{}

	// @TODO: use outbox pattern
	// if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete animal tags relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 	return api.DeleteAnimalTagsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_tags_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
