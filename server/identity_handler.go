package server

import (
	"context"
	"net/http"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/dbmodel"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/ptr"
	"github.com/google/uuid"
	"github.com/ory/client-go"
)

func (a *ApiHandler) ListIdentities(ctx context.Context, request api.ListIdentitiesRequestObject) (api.ListIdentitiesResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.ListIdentities(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}
	if request.Params.PageToken != nil {
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
	identities, _, err := req.Execute()
	if err != nil {
		return make(api.ListIdentities200JSONResponse, 0), nil
	}
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
	req := a.Kratos.Admin.IdentityAPI.GetIdentity(ctx, request.ID)
	if request.Params.IncludeCredential != nil {
		includeParams := make([]string, 0, len(*request.Params.IncludeCredential))
		for _, iparam := range *request.Params.IncludeCredential {
			includeParams = append(includeParams, string(iparam))
		}
		req = req.IncludeCredential(includeParams)
	}
	identity, _, err := req.Execute()
	if err != nil {
		return api.GetIdentity404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "identity not found"}}, nil
	}
	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}
	return api.GetIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) CreateIdentity(ctx context.Context, request api.CreateIdentityRequestObject) (api.CreateIdentityResponseObject, error) {
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
	identity, _, err := req.Execute()
	if err != nil {
		return nil, err
	}
	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}
	identityID, err := uuid.Parse(identity.Id)
	if err != nil {
		return nil, err
	}
	if _, err := a.persistor.User().Create(ctx, dbmodel.UserSetter{ID: omit.From(identityID)}); err != nil {
		return nil, err
	}
	return api.CreateIdentity201JSONResponse(resp), nil
}

func (a *ApiHandler) UpdateIdentity(ctx context.Context, request api.UpdateIdentityRequestObject) (api.UpdateIdentityResponseObject, error) {
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
	identity, _, err := req.Execute()
	if err != nil {
		return api.UpdateIdentity404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "identity not found"}}, nil
	}
	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}
	return api.UpdateIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) DeleteIdentity(ctx context.Context, request api.DeleteIdentityRequestObject) (api.DeleteIdentityResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.DeleteIdentity(ctx, request.ID)
	_, err := req.Execute()
	if err != nil {
		return api.DeleteIdentity404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "identity not found"}}, nil
	}
	return api.DeleteIdentity204Response{}, nil
}

func (a *ApiHandler) PatchIdentity(ctx context.Context, request api.PatchIdentityRequestObject) (api.PatchIdentityResponseObject, error) {
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
	identity, _, err := req.Execute()
	if err != nil {
		return api.PatchIdentity404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "identity not found"}}, nil
	}
	resp, err := dto.IdentityToResponse(*identity)
	if err != nil {
		return nil, err
	}
	return api.PatchIdentity200JSONResponse(resp), nil
}

func (a *ApiHandler) BatchPatchIdentities(ctx context.Context, request api.BatchPatchIdentitiesRequestObject) (api.BatchPatchIdentitiesResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.BatchPatchIdentities(ctx)
	if request.Body != nil {
		patch := client.PatchIdentitiesBody{}
		req = req.PatchIdentitiesBody(patch)
	}
	batchPatchIdentities, _, err := req.Execute()
	if err != nil {
		return api.BatchPatchIdentities400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusNotFound, Message: "batch patch failed"}}, nil
	}
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
	return api.BatchPatchIdentities200JSONResponse(resp), nil
}

func (a *ApiHandler) DeleteIdentityCredentials(ctx context.Context, request api.DeleteIdentityCredentialsRequestObject) (api.DeleteIdentityCredentialsResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.DeleteIdentityCredentials(ctx, request.ID, string(request.Type))
	if request.Params.Identifier != nil {
		req = req.Identifier(*request.Params.Identifier)
	}
	_, err := req.Execute()
	if err != nil {
		return api.DeleteIdentityCredentials404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "identity not found"}}, nil
	}
	return api.DeleteIdentityCredentials204Response{}, nil
}

func (a *ApiHandler) DeleteIdentitySessions(ctx context.Context, request api.DeleteIdentitySessionsRequestObject) (api.DeleteIdentitySessionsResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.DeleteIdentitySessions(ctx, request.ID)
	_, err := req.Execute()
	if err != nil {
		return api.DeleteIdentitySessions404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "identity not found"}}, nil
	}
	return api.DeleteIdentitySessions204Response{}, nil
}

func (a *ApiHandler) ListIdentitySessions(ctx context.Context, request api.ListIdentitySessionsRequestObject) (api.ListIdentitySessionsResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.ListIdentitySessions(ctx, request.ID)
	if request.Params.Active != nil {
		req = req.Active(*request.Params.Active)
	}
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}
	if request.Params.PageToken != nil {
		req = req.PageToken(*request.Params.PageToken)
	}
	identitySessions, _, err := req.Execute()
	if err != nil {
		return make(api.ListIdentitySessions200JSONResponse, 0), nil
	}
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

func (a *ApiHandler) ListIdentitySchemas(ctx context.Context, request api.ListIdentitySchemasRequestObject) (api.ListIdentitySchemasResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.ListIdentitySchemas(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}
	if request.Params.PageToken != nil {
		req = req.PageToken(*request.Params.PageToken)
	}
	schemaContainers, _, err := req.Execute()
	if err != nil {
		return make(api.ListIdentitySchemas200JSONResponse, 0), nil
	}
	resp := make(api.ListIdentitySchemas200JSONResponse, 0, len(schemaContainers))
	for _, sc := range schemaContainers {
		res, err := dto.SchemaContainerToResponse(sc)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func (a *ApiHandler) GetIdentitySchema(ctx context.Context, request api.GetIdentitySchemaRequestObject) (api.GetIdentitySchemaResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.GetIdentitySchema(ctx, request.ID)
	identitySchema, _, err := req.Execute()
	if err != nil {
		return api.GetIdentitySchema404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "schema not found"}}, nil
	}
	return api.GetIdentitySchema200JSONResponse(identitySchema), nil
}

func (a *ApiHandler) CreateRecoveryCodeForIdentity(ctx context.Context, request api.CreateRecoveryCodeForIdentityRequestObject) (api.CreateRecoveryCodeForIdentityResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.CreateRecoveryCodeForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryCodeForIdentityBody(client.CreateRecoveryCodeForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
			FlowType:   request.Body.FlowType,
		})
	}
	recoveryCodeForIdentity, _, err := req.Execute()
	if err != nil {
		return api.CreateRecoveryCodeForIdentity400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: "failed to create code"}}, nil
	}
	resp, err := dto.RecoveryCodeForIdentityToResponse(*recoveryCodeForIdentity)
	if err != nil {
		return nil, err
	}
	return api.CreateRecoveryCodeForIdentity201JSONResponse(resp), nil
}

func (a *ApiHandler) CreateRecoveryLinkForIdentity(ctx context.Context, request api.CreateRecoveryLinkForIdentityRequestObject) (api.CreateRecoveryLinkForIdentityResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.CreateRecoveryLinkForIdentity(ctx)
	if request.Body != nil {
		req = req.CreateRecoveryLinkForIdentityBody(client.CreateRecoveryLinkForIdentityBody{
			IdentityId: request.Body.IdentityID.String(),
			ExpiresIn:  request.Body.ExpiresIn,
		})
	}
	recoveryLinkForIdentity, _, err := req.Execute()
	if err != nil {
		return api.CreateRecoveryLinkForIdentity400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: "failed to create link"}}, nil
	}
	resp, err := dto.RecoveryLinkForIdentityToResponse(*recoveryLinkForIdentity)
	if err != nil {
		return nil, err
	}
	return api.CreateRecoveryLinkForIdentity200JSONResponse(resp), nil
}
