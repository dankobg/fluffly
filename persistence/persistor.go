package persistence

import (
	"context"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/google/uuid"
)

type UserPersistor interface {
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error)
	CreateUser(ctx context.Context, in UserSetter) (model.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, in UserSetter) (model.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type OrganizationPersistor interface {
	ListOrganizations(ctx context.Context) ([]OrganizationWithJoinData, error)
	GetOrganizationByID(ctx context.Context, organizationID int64) (OrganizationWithJoinData, error)
	CreateOrganization(ctx context.Context, in OrganizationCreateSetter) (model.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID int64, in OrganizationSetter) (model.Organization, error)
	DeleteOrganizationByID(ctx context.Context, organizationID int64) (int64, error)
}

type CountryPersistor interface {
	ListCountries(ctx context.Context) ([]model.Country, error)
	GetCountryByID(ctx context.Context, countryID int64) (model.Country, error)
	CreateCountry(ctx context.Context, in CountrySetter) (model.Country, error)
	UpdateCountry(ctx context.Context, countryID int64, in CountrySetter) (model.Country, error)
	DeleteCountryByID(ctx context.Context, countryID int64) (int64, error)
}

type Persistor interface {
	User() UserPersistor
	Organization() OrganizationPersistor
	Country() CountryPersistor
}
