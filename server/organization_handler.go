package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/ptr"
	"github.com/oapi-codegen/nullable"
)

func (a *ApiHandler) CreateOrganization(ctx context.Context, request api.CreateOrganizationRequestObject) (api.CreateOrganizationResponseObject, error) {
	var organizationCreateSetter dbtype.OrganizationCreateSetter

	organizationCreateSetter.Organization = dbtype.OrganizationSetter{
		Name:             nullable.NewNullableWithValue(request.Body.Name),
		Website:          request.Body.Website,
		MissionStatement: request.Body.MissionStatement,
		AdoptionPolicy:   request.Body.AdoptionPolicy,
		AdoptionURL:      request.Body.AdoptionURL,
		Distance:         request.Body.Distance,
	}

	organizationCreateSetter.Contact = dbtype.OrganizationContactSetter{
		Phone: nullable.NewNullableWithValue(request.Body.Contact.Phone),
		Email: nullable.NewNullableWithValue(string(request.Body.Contact.Email)),
	}

	organizationCreateSetter.Address = dbtype.AddressSetter{
		CountryID:     nullable.NewNullableWithValue(request.Body.Contact.Address.CountryID),
		UnitNumber:    nullable.NewNullableWithValue(*request.Body.Contact.Address.UnitNumber),
		StreetNumber:  nullable.NewNullableWithValue(*request.Body.Contact.Address.StreetNumber),
		StreetAddress: nullable.NewNullableWithValue(request.Body.Contact.Address.StreetAddress),
		City:          nullable.NewNullableWithValue(request.Body.Contact.Address.City),
		Region:        nullable.NewNullableWithValue(*request.Body.Contact.Address.Region),
		PostalCode:    nullable.NewNullableWithValue(*request.Body.Contact.Address.PostalCode),
		Lat:           nullable.NewNullableWithValue(float64(*request.Body.Contact.Address.Lat)),
		Lng:           nullable.NewNullableWithValue(float64(*request.Body.Contact.Address.Lng)),
		Note:          nullable.NewNullableWithValue(*request.Body.Contact.Address.Note),
	}

	if request.Body.WorkHour != nil {
		organizationCreateSetter.WorkHour = nullable.NewNullableWithValue(dbtype.OrganizationWorkHourSetter{
			Monday:    request.Body.WorkHour.Monday,
			Tuesday:   request.Body.WorkHour.Tuesday,
			Wednesday: request.Body.WorkHour.Wednesday,
			Thursday:  request.Body.WorkHour.Thursday,
			Friday:    request.Body.WorkHour.Friday,
			Saturday:  request.Body.WorkHour.Saturday,
			Sunday:    request.Body.WorkHour.Sunday,
		})
	}

	if request.Body.Photos.IsSpecified() && !request.Body.Photos.IsNull() {
		organizationPhotoSetters := make([]dbtype.OrganizationPhotoSetter, 0)
		for _, photo := range request.Body.Photos.MustGet() {
			photoSetter := dbtype.OrganizationPhotoSetter{}
			if photo.Small.IsSpecified() {
				photoSetter.Small = photo.Small
			}
			if photo.Medium.IsSpecified() {
				photoSetter.Medium = photo.Medium
			}
			if photo.Large.IsSpecified() {
				photoSetter.Large = photo.Large
			}
			if photo.Full.IsSpecified() {
				photoSetter.Full = photo.Full
			}
			organizationPhotoSetters = append(organizationPhotoSetters, photoSetter)
		}
		organizationCreateSetter.Photos = nullable.NewNullableWithValue(organizationPhotoSetters)
	}

	if request.Body.Socials.IsSpecified() && !request.Body.Socials.IsNull() {
		organizationSocialsSetters := make([]dbtype.OrganizationSocialSetter, 0)
		if request.Body.Socials.IsSpecified() {
			for _, social := range request.Body.Socials.MustGet() {
				organizationSocialsSetters = append(organizationSocialsSetters, dbtype.OrganizationSocialSetter{
					Platform: nullable.NewNullableWithValue(social.Platform),
					URL:      nullable.NewNullableWithValue(social.URL),
				})
			}
		}
		organizationCreateSetter.Socials = nullable.NewNullableWithValue(organizationSocialsSetters)
	}

	organization, err := a.persistor.Organization().CreateOrganization(ctx, organizationCreateSetter)
	if err != nil {
		msg := "could not create an organzation"
		var e1 postgres.ErrOrganizationUniqueViolation
		if errors.As(err, &e1) {
			msg += ", duplicate " + e1.Name
			return api.CreateOrganization400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			msg += ", organization integrity error"
			return api.CreateOrganization400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		return nil, fmt.Errorf("failed to create an organzation")
	}
	resp := api.CreateOrganization201JSONResponse(dto.OrganizationToResponse(organization))
	return resp, nil
}

func (a *ApiHandler) UpdateOrganization(ctx context.Context, request api.UpdateOrganizationRequestObject) (api.UpdateOrganizationResponseObject, error) {
	organizationSetter := dbtype.OrganizationSetter{
		Name:             request.Body.Name,
		Website:          request.Body.Website,
		MissionStatement: request.Body.MissionStatement,
		AdoptionPolicy:   request.Body.AdoptionPolicy,
		AdoptionURL:      request.Body.AdoptionURL,
		Distance:         request.Body.Distance,
	}

	organization, err := a.persistor.Organization().UpdateOrganization(ctx, request.ID, organizationSetter)
	if err != nil {
		msg := "could not update an organzation"
		var e1 postgres.ErrOrganizationUniqueViolation
		if errors.As(err, &e1) {
			msg += ", duplicate " + e1.Name
			return api.UpdateOrganization400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			msg += ", organization integrity error"
			return api.UpdateOrganization400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		return nil, fmt.Errorf("failed to update an organzation")
	}
	resp := api.UpdateOrganization201JSONResponse(dto.OrganizationToResponse(organization))
	return resp, nil
}

func (a *ApiHandler) DeleteOrganization(ctx context.Context, request api.DeleteOrganizationRequestObject) (api.DeleteOrganizationResponseObject, error) {
	_, err := a.persistor.Organization().DeleteOrganizationByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an organization by id: %w", err)
	}
	resp := api.DeleteOrganization204Response{}
	return resp, nil
}

func (a *ApiHandler) GetOrganization(ctx context.Context, request api.GetOrganizationRequestObject) (api.GetOrganizationResponseObject, error) {
	organizationWithJoinData, err := a.persistor.Organization().GetOrganizationByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationNotFound) {
			return api.GetOrganization404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Organization not found"}}, nil
		}
		return nil, fmt.Errorf("failed to get an organization by id: %w", err)
	}
	resp := api.GetOrganization200JSONResponse(dto.OrganizationWithJoinDataToResponse(organizationWithJoinData))
	return resp, nil
}

func (a *ApiHandler) ListOrganizations(ctx context.Context, request api.ListOrganizationsRequestObject) (api.ListOrganizationsResponseObject, error) {
	var filters dbtype.OrganizationFilters
	filters.Pagination = ptr.Of(getPaginationParams(request.Params.Page, request.Params.PageSize))
	organizations, err := a.persistor.Organization().ListOrganizations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	organizationsData := make([]api.Organization, len(organizations.Data))
	for i, organizationWithJoinData := range organizations.Data {
		organizationsData[i] = dto.OrganizationWithJoinDataToResponse(organizationWithJoinData)
	}
	resp := api.ListOrganizations200JSONResponse{
		Data: organizationsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, organizations.TotalCount),
	}
	return resp, nil
}
