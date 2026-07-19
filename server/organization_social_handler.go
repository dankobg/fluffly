package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListOrganizationSocials(ctx context.Context, request api.ListOrganizationSocialsRequestObject) (api.ListOrganizationSocialsResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListOrganizationSocials403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_socials_permission", "permission denied")}, nil
	}

	filters := dbtype.ListOrganizationSocialsFilters{ListOrganizationSocialsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	orgSocials, err := a.persistor.Organization().ListOrganizationSocials(ctx, request.ID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization socials: %w", err)
	}

	orgSocialsData := make([]api.OrganizationSocial, len(orgSocials.Data))
	for i, social := range orgSocials.Data {
		orgSocialsData[i] = dto.OrganizationSocialToResponse(social)
	}

	resp := api.ListOrganizationSocials200JSONResponse{
		Data: orgSocialsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, orgSocials.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetOrganizationSocial(ctx context.Context, request api.GetOrganizationSocialRequestObject) (api.GetOrganizationSocialResponseObject, error) {
	orgSocial, err := a.persistor.Organization().GetOrganizationSocial(ctx, request.ID, request.SocialID)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationSocialNotFound) {
			return api.GetOrganizationSocial404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("organization_social_not_found", "organization social not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an organization social by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetOrganizationSocialdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_social_permission", "permission denied")}, nil
	}

	resp := api.GetOrganizationSocial200JSONResponse(dto.OrganizationSocialToResponse(orgSocial))

	return resp, nil
}

func (a *ApiHandler) CreateOrganizationSocials(ctx context.Context, request api.CreateOrganizationSocialsRequestObject) (api.CreateOrganizationSocialsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateOrganizationSocials403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_socials_permission", "permission denied")}, nil
	}

	organizationCreateSocialSetters := make([]models.OrganizationSocialSetter, len(request.Body.Socials))
	for i, social := range request.Body.Socials {
		organizationCreateSocialSetters[i] = models.OrganizationSocialSetter{
			OrganizationID: omit.From(request.ID),
			Platform:       omit.From(social.Platform),
			URL:            omit.From(social.URL),
		}
	}

	orgSocials, err := a.persistor.Organization().CreateOrganizationSocials(ctx, request.ID, organizationCreateSocialSetters)
	if err != nil {
		msg := "could not create organization socials"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateOrganizationSocials400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_socials_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.CreateOrganizationSocials400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_socials_save", msg, reason)}, nil
		}

		return api.CreateOrganizationSocialsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_socials_save", msg, reason)}, nil
	}

	var data []api.OrganizationSocial
	for _, s := range orgSocials {
		data = append(data, dto.OrganizationSocialToResponse(s))
	}

	resp := api.CreateOrganizationSocials200JSONResponse(api.CreateOrganizationSocials200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateOrganizationSocialsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_socials_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizationSocials(ctx context.Context, request api.DeleteOrganizationSocialsRequestObject) (api.DeleteOrganizationSocialsResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteOrganizationSocials204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizationSocials403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_socials_permission", "permission denied")}, nil
	}

	if err := a.persistor.Organization().DeleteOrganizationSocials(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a organization socials by ids: %w", err)
	}

	resp := api.DeleteOrganizationSocials204Response{}

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete organization socials relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
	// 		return api.DeleteOrganizationSocialsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_socials_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateOrganizationSocial(ctx context.Context, request api.UpdateOrganizationSocialRequestObject) (api.UpdateOrganizationSocialResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateOrganizationSocial403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_social_permission", "permission denied")}, nil
	}

	orgSocialSetter := models.OrganizationSocialSetter{
		Platform: omit.From(request.Body.Platform),
		URL:      omit.From(request.Body.URL),
	}

	orgSocial, err := a.persistor.Organization().UpdateOrganizationSocial(ctx, request.ID, request.SocialID, orgSocialSetter)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationSocialNotFound) {
			return api.UpdateOrganizationSocial404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("organization_social_not_found", "organization social not found")}, nil
		}

		msg := "could not update a social"

		var (
			reason string
			e1     postgres.ErrCountryUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateOrganizationSocial400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_social_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "social integrity error"
			return api.UpdateOrganizationSocial400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_social_update", msg, reason)}, nil
		}

		return api.UpdateOrganizationSocialdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_social_update", msg, reason)}, nil
	}

	resp := api.UpdateOrganizationSocial200JSONResponse(dto.OrganizationSocialToResponse(orgSocial))

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizationSocial(ctx context.Context, request api.DeleteOrganizationSocialRequestObject) (api.DeleteOrganizationSocialResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizationSocial403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_social_permission", "permission denied")}, nil
	}

	if _, err := a.persistor.Organization().DeleteOrganizationSocial(ctx, request.ID, request.SocialID); err != nil {
		return nil, fmt.Errorf("failed to delete a organization social: %w", err)
	}

	resp := api.DeleteOrganizationSocial204Response{}

	// @TODO: use outbox pattern
	// if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete organization socials relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
	// 	return api.DeleteOrganizationSocialsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_socials_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
