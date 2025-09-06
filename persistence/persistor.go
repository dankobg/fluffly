package persistence

import (
	"context"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/google/uuid"
)

type UserPersistor interface {
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error)
	CreateUser(ctx context.Context, in dbtype.UserSetter) (model.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, in dbtype.UserSetter) (model.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type CountryPersistor interface {
	ListCountries(ctx context.Context, filters dbtype.CountryFilters) (dbtype.PagedResult[model.Country], error)
	GetCountryByID(ctx context.Context, countryID int64) (model.Country, error)
	CreateCountry(ctx context.Context, in dbtype.CountrySetter) (model.Country, error)
	UpdateCountry(ctx context.Context, countryID int64, in dbtype.CountrySetter) (model.Country, error)
	DeleteCountryByID(ctx context.Context, countryID int64) (int64, error)
}

type OrganizationPersistor interface {
	ListOrganizations(ctx context.Context, filters dbtype.OrganizationFilters) (dbtype.PagedResult[dbtype.OrganizationWithJoinData], error)
	GetOrganizationByID(ctx context.Context, organizationID int64) (dbtype.OrganizationWithJoinData, error)
	CreateOrganization(ctx context.Context, in dbtype.OrganizationCreateSetter) (model.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID int64, in dbtype.OrganizationSetter) (model.Organization, error)
	DeleteOrganizationByID(ctx context.Context, organizationID int64) (int64, error)
}

type AnimalPersistor interface {
	ListAnimals(ctx context.Context, filters dbtype.AnimalFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error)
	GetAnimalByID(ctx context.Context, animalID int64) (dbtype.AnimalWithJoinData, error)
	CreateAnimal(ctx context.Context, in dbtype.AnimalCreateSetter) (model.Animal, error)
	UpdateAnimal(ctx context.Context, animalID int64, in dbtype.AnimalSetter) (model.Animal, error)
	DeleteAnimalByID(ctx context.Context, animalID int64) (int64, error)
	LikeAnimal(ctx context.Context, userID uuid.UUID, animalID int64) error
	UnlikeAnimal(ctx context.Context, userID uuid.UUID, animalID int64) error
}

type Persistor interface {
	User() UserPersistor
	Organization() OrganizationPersistor
	Country() CountryPersistor
	Animal() AnimalPersistor
}
