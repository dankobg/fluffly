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
	"github.com/google/uuid"
	orykratos "github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListIdentities(ctx context.Context, request api.ListIdentitiesRequestObject) (api.ListIdentitiesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListIdentities403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.ListIdentities(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	if request.Params.Ids != nil {
		req = req.Ids(*request.Params.Ids)
	}

	if request.Params.Consistency != nil {
		req = req.Consistency(string(*request.Params.Consistency))
	}

	if request.Params.CredentialsIdentifier != nil {
		req = req.CredentialsIdentifier(*request.Params.CredentialsIdentifier)
	}

	if request.Params.IncludeCredential != nil {
		req = req.IncludeCredential(*request.Params.IncludeCredential)
	}

	if request.Params.PreviewCredentialsIdentifierSimilar != nil {
		req = req.PreviewCredentialsIdentifierSimilar(*request.Params.PreviewCredentialsIdentifierSimilar)
	}

	identities, identitiesResp, err := req.Execute()
	if err != nil {
		return api.ListIdentitiesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identities_list", "failed to list identities")}, nil
	}

	defer func() { _ = identitiesResp.Body.Close() }()

	resp := make(api.ListIdentities200JSONResponse, 0, len(identities))
	for _, identity := range identities {
		res, err := dto.IdentityToResponse(identity)
		if err != nil {
			return nil, err
		}

		resp = append(resp, res)
	}

	return resp, nil
}

func (a *ApiHandler) GetIdentity(ctx context.Context, request api.GetIdentityRequestObject) (api.GetIdentityResponseObject, error) {
	sess := GetSession(ctx)

	req := a.Kratos.Admin.IdentityAPI.GetIdentity(ctx, request.ID)
	if request.Params.IncludeCredential != nil {
		includeParams := make([]string, 0, len(*request.Params.IncludeCredential))
		for _, iparam := range *request.Params.IncludeCredential {
			includeParams = append(includeParams, string(iparam))
		}

		req = req.IncludeCredential(includeParams)
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return api.GetIdentity404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("identity_not_found", "identity not found")}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	return api.GetIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) CreateIdentity(ctx context.Context, request api.CreateIdentityRequestObject) (api.CreateIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.CreateIdentity(ctx)

	if request.Body != nil {
		var credentials *orykratos.IdentityWithCredentials
		if request.Body.Credentials != nil {
			credentials = orykratos.NewIdentityWithCredentials()
			if request.Body.Credentials.Password != nil {
				credentials.Password = orykratos.NewIdentityWithCredentialsPassword()
				if request.Body.Credentials.Password.Config != nil {
					credentials.Password.Config = orykratos.NewIdentityWithCredentialsPasswordConfig()
				}

				if request.Body.Credentials.Password.Config.Password != nil {
					credentials.Password.Config.Password = request.Body.Credentials.Password.Config.Password
				}

				if request.Body.Credentials.Password.Config.HashedPassword != nil {
					credentials.Password.Config.HashedPassword = request.Body.Credentials.Password.Config.HashedPassword
				}
			}

			if request.Body.Credentials.Oidc != nil {
				credentials.Oidc = orykratos.NewIdentityWithCredentialsOidc()
				if request.Body.Credentials.Oidc.Config != nil {
					credentials.Oidc.Config = orykratos.NewIdentityWithCredentialsOidcConfig()
					if request.Body.Credentials.Oidc.Config.Providers != nil {
						providers := make([]orykratos.IdentityWithCredentialsOidcConfigProvider, 0, len(*request.Body.Credentials.Oidc.Config.Providers))
						for _, p := range *request.Body.Credentials.Oidc.Config.Providers {
							providers = append(providers, orykratos.IdentityWithCredentialsOidcConfigProvider{
								Provider: p.Provider,
								Subject:  p.Subject,
							})
						}

						credentials.Oidc.Config.Providers = providers
					}
				}
			}
		}

		recoveryAddresses := make([]orykratos.RecoveryIdentityAddress, 0)

		if request.Body.RecoveryAddresses != nil {
			for _, recAddr := range *request.Body.RecoveryAddresses {
				recoveryAddresses = append(recoveryAddresses, orykratos.RecoveryIdentityAddress{
					Id:        new(recAddr.ID.String()),
					Value:     recAddr.Value,
					Via:       recAddr.Via,
					CreatedAt: recAddr.CreatedAt,
					UpdatedAt: recAddr.UpdatedAt,
				})
			}
		}

		verifiableAddresses := make([]orykratos.VerifiableIdentityAddress, 0)

		if request.Body.VerifiableAddresses != nil {
			for _, verAddr := range *request.Body.VerifiableAddresses {
				var id *string
				if verAddr.ID != nil {
					id = new(verAddr.ID.String())
				}

				verifiableAddress := orykratos.VerifiableIdentityAddress{
					Id:        id,
					Status:    verAddr.Status,
					Value:     verAddr.Value,
					Verified:  verAddr.Verified,
					Via:       string(verAddr.Via),
					CreatedAt: verAddr.CreatedAt,
					UpdatedAt: verAddr.UpdatedAt,
				}
				if verAddr.VerifiedAt.IsSpecified() && !verAddr.VerifiedAt.IsNull() {
					verifiableAddress.VerifiedAt = new(verAddr.VerifiedAt.MustGet().UTC())
				}

				verifiableAddresses = append(verifiableAddresses, verifiableAddress)
			}
		}

		req = req.CreateIdentityBody(orykratos.CreateIdentityBody{
			Credentials:         credentials,
			MetadataAdmin:       request.Body.MetadataAdmin,
			MetadataPublic:      request.Body.MetadataPublic,
			RecoveryAddresses:   recoveryAddresses,
			SchemaId:            request.Body.SchemaID,
			State:               (*string)(request.Body.State),
			Traits:              request.Body.Traits,
			VerifiableAddresses: verifiableAddresses,
		})
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return nil, err
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	identityID, err := uuid.Parse(identity.Id)
	if err != nil {
		return nil, err
	}

	if _, err := a.persistor.User().CreateUser(ctx, models.UserSetter{ID: omit.From(identityID)}); err != nil {
		return nil, err
	}

	// @TODO: use outbox pattern
	if err := createUserRelationTuples(ctx, a.Keto, identity.Id); err != nil {
		a.Log.Error("failed to insert identity relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Any("error", err))
		return api.CreateIdentitydefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "identity_permissions", "failed to create permissions")}, nil
	}

	return api.CreateIdentity201JSONResponse(resp), nil
}

func (a *ApiHandler) UpdateIdentity(ctx context.Context, request api.UpdateIdentityRequestObject) (api.UpdateIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.UpdateIdentity(ctx, request.ID)
	if request.Body != nil {
		var credentials *orykratos.IdentityWithCredentials
		if request.Body.Credentials != nil {
			credentials = orykratos.NewIdentityWithCredentials()
			if request.Body.Credentials.Password != nil {
				credentials.Password = orykratos.NewIdentityWithCredentialsPassword()
				if request.Body.Credentials.Password.Config != nil {
					credentials.Password.Config = orykratos.NewIdentityWithCredentialsPasswordConfig()
				}

				if request.Body.Credentials.Password.Config.Password != nil {
					credentials.Password.Config.Password = request.Body.Credentials.Password.Config.Password
				}

				if request.Body.Credentials.Password.Config.HashedPassword != nil {
					credentials.Password.Config.HashedPassword = request.Body.Credentials.Password.Config.HashedPassword
				}
			}

			if request.Body.Credentials.Oidc != nil {
				credentials.Oidc = orykratos.NewIdentityWithCredentialsOidc()
				if request.Body.Credentials.Oidc.Config != nil {
					credentials.Oidc.Config = orykratos.NewIdentityWithCredentialsOidcConfig()
					if request.Body.Credentials.Oidc.Config.Providers != nil {
						providers := make([]orykratos.IdentityWithCredentialsOidcConfigProvider, 0, len(*request.Body.Credentials.Oidc.Config.Providers))
						for _, p := range *request.Body.Credentials.Oidc.Config.Providers {
							providers = append(providers, orykratos.IdentityWithCredentialsOidcConfigProvider{
								Provider: p.Provider,
								Subject:  p.Subject,
							})
						}

						credentials.Oidc.Config.Providers = providers
					}
				}
			}
		}

		req = req.UpdateIdentityBody(orykratos.UpdateIdentityBody{
			Credentials:    credentials,
			MetadataAdmin:  request.Body.MetadataAdmin,
			MetadataPublic: request.Body.MetadataPublic,
			SchemaId:       request.Body.SchemaID,
			State:          string(request.Body.State),
			Traits:         request.Body.Traits,
		})
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return api.UpdateIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_edit", "failed to edit identity")}, nil
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	return api.UpdateIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) DeleteIdentity(ctx context.Context, request api.DeleteIdentityRequestObject) (api.DeleteIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DeleteIdentity(ctx, request.ID)

	deleteIdentityResp, err := req.Execute()
	if err != nil {
		return api.DeleteIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete", "failed to delete identity")}, nil
	}

	defer func() { _ = deleteIdentityResp.Body.Close() }()

	// @TODO: use outbox pattern
	if err := deleteUserRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete identity relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Any("error", err))
		return api.DeleteIdentitydefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "identity_permissions", "failed to delete permissions")}, nil
	}

	return api.DeleteIdentity204Response{}, nil
}

func (a *ApiHandler) PatchIdentity(ctx context.Context, request api.PatchIdentityRequestObject) (api.PatchIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.PatchIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.PatchIdentity(ctx, request.ID)
	if request.Body != nil {
		patches := make([]orykratos.JsonPatch, 0, len(*request.Body))
		for _, x := range *request.Body {
			patches = append(patches, orykratos.JsonPatch{
				From:  x.From,
				Op:    string(x.Op),
				Path:  x.Path,
				Value: x.Value,
			})
		}

		req = req.JsonPatch(patches)
	}

	identity, identityResp, err := req.Execute()
	if err != nil {
		return api.PatchIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_patch", "failed to patch identity")}, nil
	}

	defer func() { _ = identityResp.Body.Close() }()

	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}

	// @TODO: delete tuples if patch OP == "remove"

	return api.PatchIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) BatchPatchIdentities(ctx context.Context, request api.BatchPatchIdentitiesRequestObject) (api.BatchPatchIdentitiesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.BatchPatchIdentities403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.BatchPatchIdentities(ctx)

	if request.Body != nil {
		patch := orykratos.PatchIdentitiesBody{}
		req = req.PatchIdentitiesBody(patch)
	}

	batchPatchIdentities, batchPatchIdentitiesResp, err := req.Execute()
	if err != nil {
		return api.BatchPatchIdentitiesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_batch_patch", "failed to batch patch identity")}, nil
	}

	defer func() { _ = batchPatchIdentitiesResp.Body.Close() }()

	identitiesPatches := make([]api.IdentityPatchResponse, 0)

	for _, x := range batchPatchIdentities.Identities {
		var (
			defaultErr   any
			identityUUID *api.UUID
		)

		if x.Identity != nil {
			identityUUIDParsed, err := uuid.Parse(*x.Identity)
			if err != nil {
				defaultErr = err
			} else {
				identityUUID = new(identityUUIDParsed)
			}
		}

		var patchUUID *api.UUID

		if x.PatchId != nil {
			patchUUIDParsed, err := uuid.Parse(*x.PatchId)
			if err != nil {
				defaultErr = err
			} else {
				patchUUID = new(patchUUIDParsed)
			}
		}

		identitiesPatches = append(identitiesPatches, api.IdentityPatchResponse{
			Action:   (*api.IdentityPatchResponseAction)(x.Action),
			Identity: identityUUID,
			PatchID:  patchUUID,
			Error:    &defaultErr,
		})
	}

	resp := api.BatchPatchIdentitiesResponse{
		Identities: &identitiesPatches,
	}

	// @TODO: delete tuples if patch OP == "remove"

	return api.BatchPatchIdentities200JSONResponse(resp), nil
}

func (a *ApiHandler) DeleteIdentityCredentials(ctx context.Context, request api.DeleteIdentityCredentialsRequestObject) (api.DeleteIdentityCredentialsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteIdentityCredentials403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DeleteIdentityCredentials(ctx, request.ID, string(request.Type))
	if request.Params.Identifier != nil {
		req = req.Identifier(*request.Params.Identifier)
	}

	deleteIdentityCredentialsResp, err := req.Execute()
	if err != nil {
		return api.DeleteIdentityCredentialsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete_credentials", "failed to delete identity")}, nil
	}

	defer func() { _ = deleteIdentityCredentialsResp.Body.Close() }()

	return api.DeleteIdentityCredentials204Response{}, nil
}

func (a *ApiHandler) DeleteIdentitySessions(ctx context.Context, request api.DeleteIdentitySessionsRequestObject) (api.DeleteIdentitySessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteIdentitySessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DeleteIdentitySessions(ctx, request.ID)

	deleteIdentitySessionsResp, err := req.Execute()
	if err != nil {
		return api.DeleteIdentitySessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete_sessions", "failed to delete identity sessions")}, nil
	}

	defer func() { _ = deleteIdentitySessionsResp.Body.Close() }()

	return api.DeleteIdentitySessions204Response{}, nil
}

func (a *ApiHandler) ListIdentitySessions(ctx context.Context, request api.ListIdentitySessionsRequestObject) (api.ListIdentitySessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListIdentitySessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.ListIdentitySessions(ctx, request.ID)
	if request.Params.Active != nil {
		req = req.Active(*request.Params.Active)
	}

	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}

	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}

	identitySessions, identitySessionsResp, err := req.Execute()
	if err != nil {
		return api.ListIdentitySessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identities_list", "failed to list identities")}, nil
	}

	defer func() { _ = identitySessionsResp.Body.Close() }()

	resp := make(api.ListIdentitySessions200JSONResponse, 0, len(identitySessions))
	for _, sess := range identitySessions {
		res, err := dto.SessionToResponse(sess)
		if err != nil {
			return nil, err
		}

		resp = append(resp, res)
	}

	return resp, nil
}

func (a *ApiHandler) CreateRecoveryCodeForIdentity(ctx context.Context, request api.CreateRecoveryCodeForIdentityRequestObject) (api.CreateRecoveryCodeForIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.Body.IdentityID.String()),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateRecoveryCodeForIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.CreateRecoveryCodeForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryCodeForIdentityBody(orykratos.CreateRecoveryCodeForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
			FlowType:   request.Body.FlowType,
		})
	}

	recoveryCodeForIdentity, recoveryCodeForIdentityResp, err := req.Execute()
	if err != nil {
		return api.CreateRecoveryCodeForIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_create_recovery_code", "failed to create recovery code for identity")}, nil
	}

	defer func() { _ = recoveryCodeForIdentityResp.Body.Close() }()

	resp, err := dto.RecoveryCodeForIdentityToResponse(*recoveryCodeForIdentity)
	if err != nil {
		return nil, err
	}

	return api.CreateRecoveryCodeForIdentity201JSONResponse(resp), nil
}

func (a *ApiHandler) CreateRecoveryLinkForIdentity(ctx context.Context, request api.CreateRecoveryLinkForIdentityRequestObject) (api.CreateRecoveryLinkForIdentityResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    shared.AuthzIdentityID(request.Body.IdentityID.String()),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateRecoveryLinkForIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "permission denied")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.CreateRecoveryLinkForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryLinkForIdentityBody(orykratos.CreateRecoveryLinkForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
		})
	}

	recoveryLinkForIdentity, recoveryLinkForIdentityResp, err := req.Execute()
	if err != nil {
		return api.CreateRecoveryLinkForIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_create_recovery_link", "failed to create recovery link for identity")}, nil
	}

	defer func() { _ = recoveryLinkForIdentityResp.Body.Close() }()

	resp, err := dto.RecoveryLinkForIdentityToResponse(*recoveryLinkForIdentity)
	if err != nil {
		return nil, err
	}

	return api.CreateRecoveryLinkForIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) ListMyAnimals(ctx context.Context, request api.ListMyAnimalsRequestObject) (api.ListMyAnimalsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListMyAnimalsdefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("animal_permission", "permission denied")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)

	filters := dbtype.ListMyAnimalsFilters{ListMyAnimalsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	if request.Params.Properties != nil && request.Params.AnimalSpecieID != nil && len(*request.Params.AnimalSpecieID) > 0 {
		speciesSchema, err := a.LoadSpeciesPropertiesJsonSchema((*request.Params.AnimalSpecieID)[0])
		if err != nil {
			return api.ListMyAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("schema_error", "invalid json schema")}, nil
		}

		propertiesFilters, err := dbtype.CreateDynamicPropertiesFilter(speciesSchema, *request.Params.Properties)
		if err != nil {
			return api.ListMyAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("filters_err", "failed to create dynamic filters")}, nil
		}

		filters.PropertiesFilters = propertiesFilters
	}

	animals, err := a.persistor.User().ListMyAnimals(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animals: %w", err)
	}

	animalsData := make([]api.Animal, len(animals.Data))
	for i, animalWithJoinData := range animals.Data {
		animalsData[i] = dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader)
	}

	resp := api.ListMyAnimals200JSONResponse{
		Data: animalsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animals.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListMyFavoriteAnimals(ctx context.Context, request api.ListMyFavoriteAnimalsRequestObject) (api.ListMyFavoriteAnimalsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListMyFavoriteAnimalsdefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("animal_permission", "permission denied")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)

	filters := dbtype.ListMyFavoriteAnimalsFilters{ListMyFavoriteAnimalsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	if request.Params.Properties != nil && request.Params.AnimalSpecieID != nil && len(*request.Params.AnimalSpecieID) > 0 {
		speciesSchema, err := a.LoadSpeciesPropertiesJsonSchema((*request.Params.AnimalSpecieID)[0])
		if err != nil {
			return api.ListMyFavoriteAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("schema_error", "invalid json schema")}, nil
		}

		propertiesFilters, err := dbtype.CreateDynamicPropertiesFilter(speciesSchema, *request.Params.Properties)
		if err != nil {
			return api.ListMyFavoriteAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("filters_err", "failed to create dynamic filters")}, nil
		}

		filters.PropertiesFilters = propertiesFilters
	}

	animals, err := a.persistor.User().ListMyFavoriteAnimals(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animals: %w", err)
	}

	animalsData := make([]api.Animal, len(animals.Data))
	for i, animalWithJoinData := range animals.Data {
		animalsData[i] = dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader)
	}

	resp := api.ListMyFavoriteAnimals200JSONResponse{
		Data: animalsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animals.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) ListMyAdoptions(ctx context.Context, request api.ListMyAdoptionsRequestObject) (api.ListMyAdoptionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListMyAdoptionsdefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("animal_permission", "permission denied")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)

	filters := dbtype.ListMyAdoptionsFilters{ListMyAdoptionsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	if request.Params.Properties != nil && request.Params.AnimalSpecieID != nil && len(*request.Params.AnimalSpecieID) > 0 {
		speciesSchema, err := a.LoadSpeciesPropertiesJsonSchema((*request.Params.AnimalSpecieID)[0])
		if err != nil {
			return api.ListMyAdoptionsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("schema_error", "invalid json schema")}, nil
		}

		propertiesFilters, err := dbtype.CreateDynamicPropertiesFilter(speciesSchema, *request.Params.Properties)
		if err != nil {
			return api.ListMyAdoptionsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("filters_err", "failed to create dynamic filters")}, nil
		}

		filters.PropertiesFilters = propertiesFilters
	}

	animals, err := a.persistor.User().ListMyAdoptions(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animals: %w", err)
	}

	animalsData := make([]api.Animal, len(animals.Data))
	for i, animalWithJoinData := range animals.Data {
		animalsData[i] = dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader)
	}

	resp := api.ListMyAdoptions200JSONResponse{
		Data: animalsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animals.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) LikeAnimal(ctx context.Context, request api.LikeAnimalRequestObject) (api.LikeAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "like",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.LikeAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)
	if err := a.persistor.Animal().LikeAnimal(ctx, userID, request.ID); err != nil {
		msg := "failed to like an animal"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.LikeAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_like", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.LikeAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_like", msg, reason)}, nil
		}

		return api.LikeAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_like", msg, reason)}, nil
	}

	return api.EmptyResponseResponse{}, nil
}

func createUserRelationTuples(ctx context.Context, c *keto.Client, identityID string) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Group",
					Object:    "customer",
					Relation:  "members",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    shared.AuthzIdentityID(identityID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    shared.AuthzIdentityID(identityID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities", ""),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Analytic",
					Object:    shared.AuthzAnalyticsID(identityID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert identity relation tuples: %w", err)
	}

	return nil
}

func deleteUserRelationTuples(ctx context.Context, c *keto.Client, identityID string) error {
	ownersResp, err := c.Read.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: new("Identity"),
			Object:    new(shared.AuthzIdentityID(identityID)),
			Relation:  new("owners"),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list identity relation tuples: %w", err)
	}

	tuplesToDelete := []*rts.RelationTupleDelta{
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "customer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
			},
		},
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "developer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
			},
		},
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(identityID),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Identities", "identities", ""),
			},
		},
	}

	for _, tuple := range ownersResp.RelationTuples {
		subject := tuple.GetSubject()
		if subject == nil {
			continue
		}

		tuplesToDelete = append(tuplesToDelete, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    shared.AuthzIdentityID(identityID),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(subject.GetId()),
			},
		})
	}

	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuplesToDelete,
	}); err != nil {
		return fmt.Errorf("failed to delete identity relation tuples: %w", err)
	}

	return nil
}

func (a *ApiHandler) ListMyOrganizations(ctx context.Context, request api.ListMyOrganizationsRequestObject) (api.ListMyOrganizationsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
			// Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListMyOrganizationsdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_permission", "permission denied")}, nil
	}

	filters := dbtype.ListMyOrganizationsFilters{ListMyOrganizationsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	userID := uuid.MustParse(sess.Identity.Id)

	organizations, err := a.persistor.User().ListMyOrganizations(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list my organizations: %w", err)
	}

	organizationsData := make([]api.Organization, len(organizations.Data))
	for i, organizationWithJoinData := range organizations.Data {
		organizationsData[i] = dto.OrganizationWithJoinDataToResponse(organizationWithJoinData, a.uploader)
	}

	resp := api.ListMyOrganizations200JSONResponse{
		Data: organizationsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, organizations.TotalCount),
	}

	return resp, nil
}
