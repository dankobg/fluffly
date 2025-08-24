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
	"github.com/dankobg/fluffly/db/dbmodel"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence"
	"github.com/shopspring/decimal"
)

func (a *ApiHandler) CreateOrganization(ctx context.Context, request api.CreateOrganizationRequestObject) (api.CreateOrganizationResponseObject, error) {
	organizationSetter := &dbmodel.OrganizationSetter{
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

	organizationWorkHourSetter := &dbmodel.OrganizationWorkHourSetter{}
	if request.Body.WorkHour != nil {
		if request.Body.WorkHour.Monday.IsSpecified() {
			organizationWorkHourSetter.Monday = omitnull.From(request.Body.WorkHour.Monday.MustGet())
		}
		if request.Body.WorkHour.Tuesday.IsSpecified() {
			organizationWorkHourSetter.Tuesday = omitnull.From(request.Body.WorkHour.Tuesday.MustGet())
		}
		if request.Body.WorkHour.Wednesday.IsSpecified() {
			organizationWorkHourSetter.Wednesday = omitnull.From(request.Body.WorkHour.Wednesday.MustGet())
		}
		if request.Body.WorkHour.Thursday.IsSpecified() {
			organizationWorkHourSetter.Thursday = omitnull.From(request.Body.WorkHour.Thursday.MustGet())
		}
		if request.Body.WorkHour.Friday.IsSpecified() {
			organizationWorkHourSetter.Friday = omitnull.From(request.Body.WorkHour.Friday.MustGet())
		}
		if request.Body.WorkHour.Saturday.IsSpecified() {
			organizationWorkHourSetter.Saturday = omitnull.From(request.Body.WorkHour.Saturday.MustGet())
		}
		if request.Body.WorkHour.Sunday.IsSpecified() {
			organizationWorkHourSetter.Sunday = omitnull.From(request.Body.WorkHour.Sunday.MustGet())
		}
	}

	addressSetter := &dbmodel.AddressSetter{
		CountryID:     omit.From[int64](198),
		UnitNumber:    omitnull.From(*request.Body.Contact.Address.UnitNumber),
		StreetNumber:  omitnull.From(*request.Body.Contact.Address.StreetNumber),
		StreetAddress: omit.From(request.Body.Contact.Address.StreetAddress),
		City:          omit.From(request.Body.Contact.Address.City),
		Region:        omitnull.From(*request.Body.Contact.Address.Region),
		PostalCode:    omitnull.From(*request.Body.Contact.Address.PostalCode),
		Lat:           omitnull.From(decimal.NewFromFloat32(*request.Body.Contact.Address.Lat)),
		LNG:           omitnull.From(decimal.NewFromFloat32(*request.Body.Contact.Address.Lng)),
		Note:          omitnull.From(*request.Body.Contact.Address.Note),
	}

	organizationContactSetter := &dbmodel.OrganizationContactSetter{
		Phone: omit.From(request.Body.Contact.Phone),
		Email: omit.From(request.Body.Contact.Email),
	}

	organizationPhotoSetters := make([]*dbmodel.OrganizationPhotoSetter, 0)
	if request.Body.Photos.IsSpecified() {
		for _, photo := range request.Body.Photos.MustGet() {
			photoSetter := &dbmodel.OrganizationPhotoSetter{}
			if photo.Small.IsSpecified() {
				photoSetter.Small = omitnull.From(photo.Small.MustGet())
			}
			if photo.Medium.IsSpecified() {
				photoSetter.Medium = omitnull.From(photo.Medium.MustGet())
			}
			if photo.Large.IsSpecified() {
				photoSetter.Large = omitnull.From(photo.Large.MustGet())
			}
			if photo.Full.IsSpecified() {
				photoSetter.Full = omitnull.From(photo.Full.MustGet())
			}
			organizationPhotoSetters = append(organizationPhotoSetters, photoSetter)
		}
	}

	organizationSocialsSetters := make([]*dbmodel.OrganizationSocialSetter, 0)
	if request.Body.Socials.IsSpecified() {
		for _, social := range request.Body.Socials.MustGet() {
			organizationSocialsSetters = append(organizationSocialsSetters, &dbmodel.OrganizationSocialSetter{
				Platform: omit.From(social.Platform),
				URL:      omit.From(social.URL),
			})
		}
		organizationSetter.Website = omitnull.From(request.Body.Website.MustGet())
	}

	organization, err := a.persistor.Organization().Create(ctx, persistence.OrganizationCreate{
		Org:      organizationSetter,
		Address:  addressSetter,
		Contact:  organizationContactSetter,
		WorkHour: organizationWorkHourSetter,
		Photos:   organizationPhotoSetters,
		Socials:  organizationSocialsSetters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create an organization: %w", err)
	}

	resp := api.CreateOrganization201JSONResponse(dto.OrganizationToResponse(organization))
	return resp, nil
}

func (a *ApiHandler) UpdateOrganization(ctx context.Context, request api.UpdateOrganizationRequestObject) (api.UpdateOrganizationResponseObject, error) {
	organizationSetter := dbmodel.OrganizationSetter{}
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
	organizationRow, err := a.persistor.Organization().Get(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.GetOrganization404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Organization not found"}}, nil
		}
		return nil, fmt.Errorf("failed to get an organization by id: %w", err)
	}
	resp := api.GetOrganization200JSONResponse(dto.GetOrganizationByIdRowToResponse(organizationRow))
	return resp, nil
}

func (a *ApiHandler) ListOrganizations(ctx context.Context, request api.ListOrganizationsRequestObject) (api.ListOrganizationsResponseObject, error) {
	organizations, err := a.persistor.Organization().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	resp := make(api.ListOrganizations200JSONResponse, len(organizations))
	for i, org := range organizations {
		resp[i] = dto.ListOrganizationsRowToResponse(org)
	}
	return resp, nil
}
