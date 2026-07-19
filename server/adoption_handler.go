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
	"github.com/dankobg/fluffly/shared"
	"github.com/google/uuid"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) GetAdoption(ctx context.Context, request api.GetAdoptionRequestObject) (api.GetAdoptionResponseObject, error) {
	sess := GetSession(ctx)

	filters := dbtype.GetAdoptionByIDFilters{GetAdoptionParams: request.Params}

	adoptionWithJoinData, err := a.persistor.Adoption().GetAdoptionByID(ctx, request.ID, filters)
	if err != nil {
		if errors.Is(err, postgres.ErrAdoptionNotFound) {
			return api.GetAdoption404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("adoption_not_found", "adoption not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an adoption by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Adoption",
			Object:    shared.AuthzAdoptionID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAdoptiondefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "adoption_permission", "permission denied")}, nil
	}

	resp := api.GetAdoption200JSONResponse(dto.AdoptionWithJoinDataToResponse(adoptionWithJoinData, a.uploader))

	return resp, nil
}

func (a *ApiHandler) ListAdoptions(ctx context.Context, request api.ListAdoptionsRequestObject) (api.ListAdoptionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Adoptions",
			Object:    "adoptions",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAdoptionsdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "adoption_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAdoptionsFilters{ListAdoptionsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	adoptions, err := a.persistor.Adoption().ListAdoptions(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list adoptions: %w", err)
	}

	adoptionsData := make([]api.Adoption, len(adoptions.Data))
	for i, adoptionWithJoinData := range adoptions.Data {
		adoptionsData[i] = dto.AdoptionWithJoinDataToResponse(adoptionWithJoinData, a.uploader)
	}

	resp := api.ListAdoptions200JSONResponse{
		Data: adoptionsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, adoptions.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ApplyForAdoption(ctx context.Context, request api.ApplyForAdoptionRequestObject) (api.ApplyForAdoptionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "adopt",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ApplyForAdoption403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("adoption_apply_permission", "permission denied")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)
	adoption, err := a.persistor.Animal().ApplyForAdoption(ctx, request.ID, userID, request.Body.OrganizationID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalAlreadyAdopted) {
			return api.ApplyForAdoption400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "start_adoption", "failed to adopt animal", "animal already adopted")}, nil
		}

		return api.ApplyForAdoptiondefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "start_adoption", "failed to adopt animal", err.Error())}, nil
	}

	// @TODO: use outbox pattern
	if err := createAdoptionRelationTuples(ctx, a.Keto, sess.Identity.Id, adoption.ID); err != nil {
		a.Log.Error("failed to insert adoption relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("adoption_id", request.ID), slog.Any("error", err))
		return api.ApplyForAdoptiondefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "adoption_permissions", "failed to create permissions")}, nil
	}

	return api.EmptyResponseResponse{}, nil
}

func (a *ApiHandler) ApproveAdoption(ctx context.Context, request api.ApproveAdoptionRequestObject) (api.ApproveAdoptionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ApproveAdoption403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("approve_adoption_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().ApproveAdoption(ctx, request.AdoptionID); err != nil {
		return api.ApproveAdoptiondefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "adopt_animal", "failed to approve animal adoption", err.Error())}, nil
	}

	return api.EmptyResponseResponse{}, nil
}

func (a *ApiHandler) RejectAdoption(ctx context.Context, request api.RejectAdoptionRequestObject) (api.RejectAdoptionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.RejectAdoption403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("reject_adoption_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().RejectAdoption(ctx, request.AdoptionID); err != nil {
		return api.RejectAdoptiondefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "reject_adoption", "failed to reject animal adoption", err.Error())}, nil
	}

	return api.EmptyResponseResponse{}, nil
}

func createAdoptionRelationTuples(ctx context.Context, c *keto.Client, identityID string, adoptionID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Adoption",
					Object:    shared.AuthzAdoptionID(adoptionID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Adoption",
					Object:    shared.AuthzAdoptionID(adoptionID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Adoptions", "adoptions", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert adoption relation tuples: %w", err)
	}

	return nil
}
