package dbtype

import (
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob/types"
)

type ListMyAnimalsFilters struct {
	api.ListMyAnimalsParams
	PropertiesFilters []DynamicFilter
}

type ListMyFavoriteAnimalsFilters struct {
	api.ListMyFavoriteAnimalsParams
	PropertiesFilters []DynamicFilter
}

type ListMyAdoptionsFilters struct {
	api.ListMyAdoptionsParams
	PropertiesFilters []DynamicFilter
}

type ListAnimalsFilters struct {
	api.ListAnimalsParams
	PropertiesFilters []DynamicFilter
	UserID            *uuid.UUID
}

type GetAnimalByIDFilters struct {
	api.GetAnimalParams
	UserID *uuid.UUID
}

type AnimalWithJoinData struct {
	models.Animal
	Type                              models.AnimalType
	Specie                            models.AnimalSpecie
	Microchip                         *models.Microchip
	Breeds                            types.JSON[*[]AnimalBreedWithJoinData]
	Tags                              types.JSON[*[]models.AnimalTag]
	Photos                            types.JSON[*[]models.AnimalPhoto]
	Videos                            types.JSON[*[]models.AnimalVideo]
	Organization                      *models.Organization
	OrganizationWorkHour              *models.OrganizationWorkHour
	OrganizationContact               *models.OrganizationContact
	OrganizationContactAddress        *models.Address
	OrganizationContactAddressCountry *models.Country
	User                              *models.User
	Likes                             int64
	Liked                             null.Val[bool]
	AdoptionID                        null.Val[int64]
}

type AnimalCreateSetter struct {
	Animal    models.AnimalSetter
	Microchip omitnull.Val[models.MicrochipSetter]
	Breeds    omitnull.Val[[]models.AnimalBreedSetter]
	Tags      omitnull.Val[[]models.AnimalTagSetter]
	Photos    omitnull.Val[[]models.AnimalPhotoSetter]
	Videos    omitnull.Val[[]models.AnimalVideoSetter]
}

type AnimalUpdateSetter struct {
	Animal    omitnull.Val[models.AnimalSetter]
	Microchip omitnull.Val[models.MicrochipSetter]
	Breeds    omitnull.Val[[]models.AnimalBreedSetter]
	Tags      omitnull.Val[[]models.AnimalTagSetter]
	Photos    omitnull.Val[[]models.AnimalPhotoSetter]
	Videos    omitnull.Val[[]models.AnimalVideoSetter]
}
