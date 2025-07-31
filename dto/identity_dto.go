package dto

import (
	"fmt"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	kratos "github.com/ory/client-go"
)

func IdentityToResponse(identity kratos.Identity) (api.Identity, error) {
	credentials := map[string]api.IdentityCredentials{}
	if identity.Credentials != nil {
		for k, v := range *identity.Credentials {
			credentials[k] = api.IdentityCredentials{
				Config:      &v.Config,
				Identifiers: &v.Identifiers,
				Type:        (*api.IdentityCredentialsType)(v.Type),
				Version:     new(int64),
				CreatedAt:   v.CreatedAt,
				UpdatedAt:   v.UpdatedAt,
			}
		}
	}
	recoveryAddresses := make([]api.RecoveryIdentityAddress, 0, len(identity.RecoveryAddresses))
	for _, recAddr := range identity.RecoveryAddresses {
		id, err := uuid.Parse(recAddr.Id)
		if err != nil {
			return api.Identity{}, fmt.Errorf("failed to parse recovery address uuid: %w", err)
		}
		recoveryAddresses = append(recoveryAddresses, api.RecoveryIdentityAddress{
			ID:        id,
			Value:     recAddr.Value,
			Via:       recAddr.Via,
			CreatedAt: recAddr.CreatedAt,
			UpdatedAt: recAddr.UpdatedAt,
		})
	}
	verifiableAddresses := make([]api.VerifiableIdentityAddress, 0, len(identity.VerifiableAddresses))
	for _, verAddr := range identity.VerifiableAddresses {
		var id *api.UUID
		if verAddr.Id != nil {
			parsed, err := uuid.Parse(*verAddr.Id)
			if err != nil {
				return api.Identity{}, fmt.Errorf("failed to parse verifiable address uuid: %w", err)
			}
			id = &parsed
		}
		verifiableAddresses = append(verifiableAddresses, api.VerifiableIdentityAddress{

			ID:         id,
			Status:     verAddr.Status,
			Value:      verAddr.Value,
			Verified:   verAddr.Verified,
			VerifiedAt: nullable.NewNullableWithValue(*verAddr.VerifiedAt),
			Via:        api.VerifiableIdentityAddressVia(verAddr.Via),
			CreatedAt:  verAddr.CreatedAt,
			UpdatedAt:  verAddr.UpdatedAt,
		})
	}
	id, err := uuid.Parse(identity.Id)
	if err != nil {
		return api.Identity{}, fmt.Errorf("failed to parse identity uuid: %w", err)
	}
	resp := api.Identity{
		ID:                  id,
		Credentials:         &credentials,
		MetadataAdmin:       nullable.NewNullableWithValue(identity.MetadataAdmin),
		MetadataPublic:      nullable.NewNullableWithValue(identity.MetadataPublic),
		RecoveryAddresses:   &recoveryAddresses,
		SchemaID:            identity.SchemaId,
		SchemaURL:           identity.SchemaUrl,
		State:               (*api.IdentityState)(identity.State),
		StateChangedAt:      nullable.NewNullableWithValue(*identity.StateChangedAt),
		Traits:              identity.Traits,
		VerifiableAddresses: &verifiableAddresses,
		CreatedAt:           identity.CreatedAt,
		UpdatedAt:           identity.UpdatedAt,
	}
	return resp, nil
}

func SessionToResponse(sess kratos.Session) (api.Session, error) {
	authMethods := make([]api.SessionAuthenticationMethod, 0, len(sess.AuthenticationMethods))
	for _, am := range sess.AuthenticationMethods {
		authMethods = append(authMethods, api.SessionAuthenticationMethod{
			Aal:          (*api.AuthenticatorAssuranceLevel)(am.Aal),
			CompletedAt:  am.CompletedAt,
			Method:       (*api.SessionAuthenticationMethodMethod)(am.Method),
			Organization: am.Organization,
			Provider:     am.Provider,
		})
	}
	sessDevices := make([]api.SessionDevice, 0, len(sess.Devices))
	for _, x := range sess.Devices {
		id, err := uuid.Parse(x.Id)
		if err != nil {
			return api.Session{}, fmt.Errorf("failed to parse devices uuid: %w", err)
		}
		sessDevices = append(sessDevices, api.SessionDevice{
			ID:        id,
			IPAddress: x.IpAddress,
			Location:  x.Location,
			UserAgent: x.UserAgent,
		})
	}
	var identity *api.Identity
	if sess.Identity != nil {
		ident, err := IdentityToResponse(*sess.Identity)
		if err != nil {
			return api.Session{}, err
		}
		identity = &ident
	}
	id, err := uuid.Parse(sess.Id)
	if err != nil {
		return api.Session{}, fmt.Errorf("failed to parse session uuid: %w", err)
	}
	resp := api.Session{
		Active:                      sess.Active,
		AuthenticatedAt:             sess.AuthenticatedAt,
		AuthenticationMethods:       &authMethods,
		AuthenticatorAssuranceLevel: (*api.AuthenticatorAssuranceLevel)(sess.AuthenticatorAssuranceLevel),
		Devices:                     &sessDevices,
		ExpiresAt:                   sess.ExpiresAt,
		ID:                          id,
		Identity:                    identity,
		IssuedAt:                    sess.IssuedAt,
		Tokenized:                   sess.Tokenized,
	}
	return resp, nil
}

func SchemaContainerToResponse(sc kratos.IdentitySchemaContainer) (api.IdentitySchemaContainer, error) {
	resp := api.IdentitySchemaContainer{
		ID:     sc.Id,
		Schema: &sc.Schema,
	}
	return resp, nil
}

func RecoveryCodeForIdentityToResponse(code kratos.RecoveryCodeForIdentity) (api.RecoveryCodeForIdentity, error) {
	resp := api.RecoveryCodeForIdentity{
		ExpiresAt:    code.ExpiresAt,
		RecoveryCode: code.RecoveryCode,
		RecoveryLink: code.RecoveryLink,
	}
	return resp, nil
}

func RecoveryLinkForIdentityToResponse(link kratos.RecoveryLinkForIdentity) (api.RecoveryLinkForIdentity, error) {
	resp := api.RecoveryLinkForIdentity{
		ExpiresAt:    link.ExpiresAt,
		RecoveryLink: link.RecoveryLink,
	}
	return resp, nil
}
