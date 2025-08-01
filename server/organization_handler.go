package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/model"
	"github.com/dankobg/fluffly/dto"
)

func (a *ApiHandler) CreateOrganization(ctx context.Context, request api.CreateOrganizationRequestObject) (api.CreateOrganizationResponseObject, error) {
	organizationSetter := model.OrganizationSetter{
		// ContactID: omit.From(request.Body.ContactID),
		Name: omit.From(request.Body.Name),
	}
	if request.Body.Website.IsSpecified() {
		organizationSetter.Website = omitnull.From(request.Body.Website.MustGet())
	}
	if request.Body.MissionStatement.IsSpecified() {
		organizationSetter.MissionStatement = omitnull.From(request.Body.MissionStatement.MustGet())
	}
	if request.Body.AdoptionPolicy.IsSpecified() {
		organizationSetter.AdoptionPolicy = omitnull.From(request.Body.AdoptionPolicy.MustGet())
	}
	if request.Body.AdoptionURL.IsSpecified() {
		organizationSetter.AdoptionURL = omitnull.From(request.Body.AdoptionURL.MustGet())
	}
	if request.Body.Distance.IsSpecified() {
		organizationSetter.Distance = omitnull.From(request.Body.Distance.MustGet())
	}
	if request.Body.Facebook.IsSpecified() {
		organizationSetter.Facebook = omitnull.From(request.Body.Facebook.MustGet())
	}
	if request.Body.Twitter.IsSpecified() {
		organizationSetter.Twitter = omitnull.From(request.Body.Twitter.MustGet())
	}
	if request.Body.Youtube.IsSpecified() {
		organizationSetter.Youtube = omitnull.From(request.Body.Youtube.MustGet())
	}
	if request.Body.Instagram.IsSpecified() {
		organizationSetter.Instagram = omitnull.From(request.Body.Instagram.MustGet())
	}
	if request.Body.Pinterest.IsSpecified() {
		organizationSetter.Pinterest = omitnull.From(request.Body.Pinterest.MustGet())
	}
	organization, err := a.persistor.Organization().Create(ctx, organizationSetter)
	if err != nil {
		return nil, fmt.Errorf("failed to create an organization: %w", err)
	}
	resp := api.CreateOrganization201JSONResponse(dto.OrganizationToResponse(organization))
	return resp, nil
}

func (a *ApiHandler) UpdateOrganization(ctx context.Context, request api.UpdateOrganizationRequestObject) (api.UpdateOrganizationResponseObject, error) {
	organizationSetter := model.OrganizationSetter{}
	if request.Body.Name.IsSpecified() && !request.Body.Name.IsNull() {
		organizationSetter.Name = omit.From(request.Body.Name.MustGet())
	}
	if request.Body.Website.IsSpecified() {
		organizationSetter.Website = omitnull.From(request.Body.Website.MustGet())
	}
	if request.Body.MissionStatement.IsSpecified() {
		organizationSetter.MissionStatement = omitnull.From(request.Body.MissionStatement.MustGet())
	}
	if request.Body.AdoptionPolicy.IsSpecified() {
		organizationSetter.AdoptionPolicy = omitnull.From(request.Body.AdoptionPolicy.MustGet())
	}
	if request.Body.AdoptionURL.IsSpecified() {
		organizationSetter.AdoptionURL = omitnull.From(request.Body.AdoptionURL.MustGet())
	}
	if request.Body.Distance.IsSpecified() {
		organizationSetter.Distance = omitnull.From(request.Body.Distance.MustGet())
	}
	if request.Body.Facebook.IsSpecified() {
		organizationSetter.Facebook = omitnull.From(request.Body.Facebook.MustGet())
	}
	if request.Body.Twitter.IsSpecified() {
		organizationSetter.Twitter = omitnull.From(request.Body.Twitter.MustGet())
	}
	if request.Body.Youtube.IsSpecified() {
		organizationSetter.Youtube = omitnull.From(request.Body.Youtube.MustGet())
	}
	if request.Body.Instagram.IsSpecified() {
		organizationSetter.Instagram = omitnull.From(request.Body.Instagram.MustGet())
	}
	if request.Body.Pinterest.IsSpecified() {
		organizationSetter.Pinterest = omitnull.From(request.Body.Pinterest.MustGet())
	}
	organization, err := a.persistor.Organization().Update(ctx, request.ID, organizationSetter)
	if err != nil {
		return nil, fmt.Errorf("failed to update an organization: %w", err)
	}
	resp := api.UpdateOrganization201JSONResponse(dto.OrganizationToResponse(organization))
	return resp, nil
}

func (a *ApiHandler) DeleteOrganization(ctx context.Context, request api.DeleteOrganizationRequestObject) (api.DeleteOrganizationResponseObject, error) {
	_, err := a.persistor.Organization().Delete(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an organization by id: %w", err)
	}
	resp := api.DeleteOrganization204Response{}
	return resp, nil
}

func (a *ApiHandler) GetOrganization(ctx context.Context, request api.GetOrganizationRequestObject) (api.GetOrganizationResponseObject, error) {
	organization, err := a.persistor.Organization().Get(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.GetOrganization404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Organization not found"}}, nil
		}
		return nil, fmt.Errorf("failed to get an organization by id: %w", err)
	}
	resp := api.GetOrganization200JSONResponse(dto.OrganizationToResponse(organization))
	return resp, nil
}

func (a *ApiHandler) ListOrganizations(ctx context.Context, request api.ListOrganizationsRequestObject) (api.ListOrganizationsResponseObject, error) {
	organizations, err := a.persistor.Organization().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	resp := make(api.ListOrganizations200JSONResponse, len(organizations))
	for i, org := range organizations {
		resp[i] = dto.OrganizationToResponse(org)
	}
	return resp, nil
}
