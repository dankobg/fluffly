package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/ptr"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/ory/client-go"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListIdentities(ctx context.Context, request api.ListIdentitiesRequestObject) (api.ListIdentitiesResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identities",
			Object:    "identities",
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ListIdentities403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
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
	defer identitiesResp.Body.Close()
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
			Object:    AuthzIdentityID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}
	defer identityResp.Body.Close()
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
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.CreateIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.CreateIdentity(ctx)
	if request.Body != nil {
		var credentials *client.IdentityWithCredentials
		if request.Body.Credentials != nil {
			credentials = client.NewIdentityWithCredentials()
			if request.Body.Credentials.Password != nil {
				credentials.Password = client.NewIdentityWithCredentialsPassword()
				if request.Body.Credentials.Password.Config != nil {
					credentials.Password.Config = client.NewIdentityWithCredentialsPasswordConfig()
				}
				if request.Body.Credentials.Password.Config.Password != nil {
					credentials.Password.Config.Password = request.Body.Credentials.Password.Config.Password
				}
				if request.Body.Credentials.Password.Config.HashedPassword != nil {
					credentials.Password.Config.HashedPassword = request.Body.Credentials.Password.Config.HashedPassword
				}
			}
			if request.Body.Credentials.Oidc != nil {
				credentials.Oidc = client.NewIdentityWithCredentialsOidc()
				if request.Body.Credentials.Oidc.Config != nil {
					credentials.Oidc.Config = client.NewIdentityWithCredentialsOidcConfig()
					if request.Body.Credentials.Oidc.Config.Providers != nil {
						providers := make([]client.IdentityWithCredentialsOidcConfigProvider, 0, len(*request.Body.Credentials.Oidc.Config.Providers))
						for _, p := range *request.Body.Credentials.Oidc.Config.Providers {
							providers = append(providers, client.IdentityWithCredentialsOidcConfigProvider{
								Provider: p.Provider,
								Subject:  p.Subject,
							})
						}
						credentials.Oidc.Config.Providers = providers
					}
				}
			}
		}
		recoveryAddresses := make([]client.RecoveryIdentityAddress, 0)
		if request.Body.RecoveryAddresses != nil {
			for _, recAddr := range *request.Body.RecoveryAddresses {
				recoveryAddresses = append(recoveryAddresses, client.RecoveryIdentityAddress{
					Id:        recAddr.ID.String(),
					Value:     recAddr.Value,
					Via:       recAddr.Via,
					CreatedAt: recAddr.CreatedAt,
					UpdatedAt: recAddr.UpdatedAt,
				})
			}
		}
		verifiableAddresses := make([]client.VerifiableIdentityAddress, 0)
		if request.Body.VerifiableAddresses != nil {
			for _, verAddr := range *request.Body.VerifiableAddresses {
				var id *string
				if verAddr.ID != nil {
					id = ptr.Of(verAddr.ID.String())
				}
				verifiableAddress := client.VerifiableIdentityAddress{
					Id:        id,
					Status:    verAddr.Status,
					Value:     verAddr.Value,
					Verified:  verAddr.Verified,
					Via:       string(verAddr.Via),
					CreatedAt: verAddr.CreatedAt,
					UpdatedAt: verAddr.UpdatedAt,
				}
				if verAddr.VerifiedAt.IsSpecified() && !verAddr.VerifiedAt.IsNull() {
					verifiableAddress.VerifiedAt = ptr.Of(verAddr.VerifiedAt.MustGet().UTC())
				}
				verifiableAddresses = append(verifiableAddresses, verifiableAddress)
			}
		}
		req = req.CreateIdentityBody(client.CreateIdentityBody{
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
	defer identityResp.Body.Close()
	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}
	identityID, err := uuid.Parse(identity.Id)
	if err != nil {
		return nil, err
	}
	if _, err := a.persistor.User().CreateUser(ctx, dbtype.UserSetter{ID: nullable.NewNullableWithValue(identityID)}); err != nil {
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
			Object:    AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.UpdateIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.UpdateIdentity(ctx, request.ID)
	if request.Body != nil {
		var credentials *client.IdentityWithCredentials
		if request.Body.Credentials != nil {
			credentials = client.NewIdentityWithCredentials()
			if request.Body.Credentials.Password != nil {
				credentials.Password = client.NewIdentityWithCredentialsPassword()
				if request.Body.Credentials.Password.Config != nil {
					credentials.Password.Config = client.NewIdentityWithCredentialsPasswordConfig()
				}
				if request.Body.Credentials.Password.Config.Password != nil {
					credentials.Password.Config.Password = request.Body.Credentials.Password.Config.Password
				}
				if request.Body.Credentials.Password.Config.HashedPassword != nil {
					credentials.Password.Config.HashedPassword = request.Body.Credentials.Password.Config.HashedPassword
				}
			}
			if request.Body.Credentials.Oidc != nil {
				credentials.Oidc = client.NewIdentityWithCredentialsOidc()
				if request.Body.Credentials.Oidc.Config != nil {
					credentials.Oidc.Config = client.NewIdentityWithCredentialsOidcConfig()
					if request.Body.Credentials.Oidc.Config.Providers != nil {
						providers := make([]client.IdentityWithCredentialsOidcConfigProvider, 0, len(*request.Body.Credentials.Oidc.Config.Providers))
						for _, p := range *request.Body.Credentials.Oidc.Config.Providers {
							providers = append(providers, client.IdentityWithCredentialsOidcConfigProvider{
								Provider: p.Provider,
								Subject:  p.Subject,
							})
						}
						credentials.Oidc.Config.Providers = providers
					}
				}
			}
		}
		req = req.UpdateIdentityBody(client.UpdateIdentityBody{
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
	defer identityResp.Body.Close()
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
			Object:    AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DeleteIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DeleteIdentity(ctx, request.ID)
	deleteIdentityResp, err := req.Execute()
	if err != nil {
		return api.DeleteIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete", "failed to delete identity")}, nil

	}
	defer deleteIdentityResp.Body.Close()

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
			Object:    AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.PatchIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.PatchIdentity(ctx, request.ID)
	if request.Body != nil {
		patches := make([]client.JsonPatch, 0, len(*request.Body))
		for _, x := range *request.Body {
			patches = append(patches, client.JsonPatch{
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
	defer identityResp.Body.Close()
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
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.BatchPatchIdentities403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.BatchPatchIdentities(ctx)
	if request.Body != nil {
		patch := client.PatchIdentitiesBody{}
		req = req.PatchIdentitiesBody(patch)
	}
	batchPatchIdentities, batchPatchIdentitiesResp, err := req.Execute()
	if err != nil {
		return api.BatchPatchIdentitiesdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_batch_patch", "failed to batch patch identity")}, nil
	}
	defer batchPatchIdentitiesResp.Body.Close()
	identitiesPatches := make([]api.IdentityPatchResponse, 0)
	for _, x := range batchPatchIdentities.Identities {
		var defaultErr any
		var identityUUID *api.UUID
		if x.Identity != nil {
			identityUUIDParsed, err := uuid.Parse(*x.Identity)
			if err != nil {
				defaultErr = err
			} else {
				identityUUID = ptr.Of(identityUUIDParsed)
			}
		}
		var patchUUID *api.UUID
		if x.PatchId != nil {
			patchUUIDParsed, err := uuid.Parse(*x.PatchId)
			if err != nil {
				defaultErr = err
			} else {
				patchUUID = ptr.Of(patchUUIDParsed)
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
			Object:    AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DeleteIdentityCredentials403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DeleteIdentityCredentials(ctx, request.ID, string(request.Type))
	if request.Params.Identifier != nil {
		req = req.Identifier(*request.Params.Identifier)
	}
	deleteIdentityCredentialsResp, err := req.Execute()
	if err != nil {
		return api.DeleteIdentityCredentialsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete_credentials", "failed to delete identity")}, nil
	}
	defer deleteIdentityCredentialsResp.Body.Close()
	return api.DeleteIdentityCredentials204Response{}, nil
}

func (a *ApiHandler) DeleteIdentitySessions(ctx context.Context, request api.DeleteIdentitySessionsRequestObject) (api.DeleteIdentitySessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    AuthzIdentityID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DeleteIdentitySessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DeleteIdentitySessions(ctx, request.ID)
	deleteIdentitySessionsResp, err := req.Execute()
	if err != nil {
		return api.DeleteIdentitySessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_delete_sessions", "failed to delete identity sessions")}, nil
	}
	defer deleteIdentitySessionsResp.Body.Close()
	return api.DeleteIdentitySessions204Response{}, nil
}

func (a *ApiHandler) ListIdentitySessions(ctx context.Context, request api.ListIdentitySessionsRequestObject) (api.ListIdentitySessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Identity",
			Object:    AuthzIdentityID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ListIdentitySessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
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
	defer identitySessionsResp.Body.Close()
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
			Object:    AuthzIdentityID(request.Body.IdentityID.String()),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.CreateRecoveryCodeForIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.CreateRecoveryCodeForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryCodeForIdentityBody(client.CreateRecoveryCodeForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
			FlowType:   request.Body.FlowType,
		})
	}
	recoveryCodeForIdentity, recoveryCodeForIdentityResp, err := req.Execute()
	if err != nil {
		return api.CreateRecoveryCodeForIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_create_recovery_code", "failed to create recovery code for identity")}, nil
	}
	defer recoveryCodeForIdentityResp.Body.Close()
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
			Object:    AuthzIdentityID(request.Body.IdentityID.String()),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.CreateRecoveryLinkForIdentity403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("identity_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.CreateRecoveryLinkForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryLinkForIdentityBody(client.CreateRecoveryLinkForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
		})
	}
	recoveryLinkForIdentity, recoveryLinkForIdentityResp, err := req.Execute()
	if err != nil {
		return api.CreateRecoveryLinkForIdentitydefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "identity_create_recovery_link", "failed to create recovery link for identity")}, nil
	}
	defer recoveryLinkForIdentityResp.Body.Close()
	resp, err := dto.RecoveryLinkForIdentityToResponse(*recoveryLinkForIdentity)
	if err != nil {
		return nil, err
	}
	return api.CreateRecoveryLinkForIdentity200JSONResponse(resp), nil
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
					Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    AuthzIdentityID(identityID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Identity",
					Object:    AuthzIdentityID(identityID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Identities", "identities", ""),
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
			Namespace: ptr.Of("Identity"),
			Object:    ptr.Of(AuthzIdentityID(identityID)),
			Relation:  ptr.Of("owners"),
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
				Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
			},
		},
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Group",
				Object:    "developer",
				Relation:  "members",
				Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
			},
		},
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Identity",
				Object:    AuthzIdentityID(identityID),
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
				Object:    AuthzIdentityID(identityID),
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
