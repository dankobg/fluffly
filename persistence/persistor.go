package persistence

import (
	"context"

	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/google/uuid"
)

type UserPersistor interface {
	ListUsers(ctx context.Context, filters dbtype.ListUsersFilters) (dbtype.PagedResult[models.User], error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	CreateUser(ctx context.Context, in models.UserSetter) (models.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, in models.UserSetter) (models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	ListMyAnimals(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyAnimalsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error)
	ListMyFavoriteAnimals(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyFavoriteAnimalsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error)
	ListMyAdoptions(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyAdoptionsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error)
	ListMyOrganizations(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyOrganizationsFilters) (dbtype.PagedResult[dbtype.OrganizationWithJoinData], error)
}

type CountryPersistor interface {
	ListCountries(ctx context.Context, filters dbtype.ListCountriesFilters) (dbtype.PagedResult[models.Country], error)
	GetCountryByID(ctx context.Context, countryID int64) (models.Country, error)
	CreateCountry(ctx context.Context, in models.CountrySetter) (models.Country, error)
	UpdateCountry(ctx context.Context, countryID int64, in models.CountrySetter) (models.Country, error)
	DeleteCountryByID(ctx context.Context, countryID int64) (int64, error)
	DeleteCountries(ctx context.Context, countryIDs []int64) error
}

type OrganizationPersistor interface {
	OrganizationSocialPersistor
	OrganizationPhotoPersistor
	OrganizationVideoPersistor
	OrganizationMembershipPersistor
	OrganizationInvitationPersistor

	GetOrganizationByID(ctx context.Context, organizationID int64, filters dbtype.GetOrganizationByIDFilters) (dbtype.OrganizationWithJoinData, error)
	ListOrganizations(ctx context.Context, filters dbtype.ListOrganizationsFilters) (dbtype.PagedResult[dbtype.OrganizationWithJoinData], error)
	ApplyForOrganization(ctx context.Context, userID uuid.UUID, in dbtype.OrganizationApplyForSetter) (models.Organization, error)
	CreateOrganization(ctx context.Context, in dbtype.OrganizationCreateSetter) (models.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID int64, in dbtype.OrganizationUpdateSetter) (models.Organization, error)
	DeleteOrganizationByID(ctx context.Context, organizationID int64) (int64, error)
	DeleteOrganizations(ctx context.Context, organizationIDs []int64) error
	ApproveOrganization(ctx context.Context, organizationID int64) error
	RejectOrganization(ctx context.Context, organizationID int64) error
}

type OrganizationSocialPersistor interface {
	GetOrganizationSocial(ctx context.Context, organizationID, socialID int64) (models.OrganizationSocial, error)
	ListOrganizationSocials(ctx context.Context, organizationID int64, filters dbtype.ListOrganizationSocialsFilters) (dbtype.PagedResult[models.OrganizationSocial], error)
	UpdateOrganizationSocial(ctx context.Context, organizationID, socialID int64, in models.OrganizationSocialSetter) (models.OrganizationSocial, error)
	DeleteOrganizationSocial(ctx context.Context, organizationID, socialID int64) (int64, error)
	CreateOrganizationSocials(ctx context.Context, organizationID int64, in []models.OrganizationSocialSetter) ([]models.OrganizationSocial, error)
	DeleteOrganizationSocials(ctx context.Context, organizationID int64, socialIDs []int64) error
}

type OrganizationPhotoPersistor interface {
	GetOrganizationPhoto(ctx context.Context, organizationID, photoID int64) (models.OrganizationPhoto, error)
	ListOrganizationPhotos(ctx context.Context, organizationID int64, filters dbtype.ListOrganizationPhotosFilters) (dbtype.PagedResult[models.OrganizationPhoto], error)
	UpdateOrganizationPhoto(ctx context.Context, organizationID, photoID int64, in models.OrganizationPhotoSetter) (models.OrganizationPhoto, error)
	DeleteOrganizationPhoto(ctx context.Context, organizationID, photoID int64) (int64, error)
	CreateOrganizationPhotos(ctx context.Context, organizationID int64, in []models.OrganizationPhotoSetter) ([]models.OrganizationPhoto, error)
	DeleteOrganizationPhotos(ctx context.Context, organizationID int64, photoIDs []int64) error
}

type OrganizationVideoPersistor interface {
	GetOrganizationVideo(ctx context.Context, organizationID, videoID int64) (models.OrganizationVideo, error)
	ListOrganizationVideos(ctx context.Context, organizationID int64, filters dbtype.ListOrganizationVideosFilters) (dbtype.PagedResult[models.OrganizationVideo], error)
	UpdateOrganizationVideo(ctx context.Context, organizationID, videoID int64, in models.OrganizationVideoSetter) (models.OrganizationVideo, error)
	DeleteOrganizationVideo(ctx context.Context, organizationID, videoID int64) (int64, error)
	CreateOrganizationVideos(ctx context.Context, organizationID int64, in []models.OrganizationVideoSetter) ([]models.OrganizationVideo, error)
	DeleteOrganizationVideos(ctx context.Context, organizationID int64, videoIds []int64) error
}

type OrganizationMembershipPersistor interface {
	CreateMember(ctx context.Context, organizationID int64, userID uuid.UUID, in models.OrganizationMembershipSetter) (models.OrganizationMembership, error)
	DeleteMember(ctx context.Context, organizationID int64, userID uuid.UUID) error
	UpdateMembership(ctx context.Context, organizationID int64, userID uuid.UUID, in models.OrganizationMembershipSetter) (models.OrganizationMembership, error)
}

type OrganizationInvitationPersistor interface {
	GetInvitationByID(ctx context.Context, invitationID int64) (models.OrganizationInvitation, error)
	CreateInvitation(ctx context.Context, in models.OrganizationInvitationSetter) (models.OrganizationInvitation, error)
	DeleteInvitation(ctx context.Context, invitationID int64) error
	UpdateInvitation(ctx context.Context, invitationID int64, in models.OrganizationInvitationSetter) (models.OrganizationInvitation, error)
}

type AdoptionPersistor interface {
	GetAdoptionByID(ctx context.Context, adoptionID int64, filters dbtype.GetAdoptionByIDFilters) (dbtype.AdoptionWithJoinData, error)
	ListAdoptions(ctx context.Context, filters dbtype.ListAdoptionsFilters) (dbtype.PagedResult[dbtype.AdoptionWithJoinData], error)
}

type AnimalPersistor interface {
	AnimalBreedPersistor
	AnimalTagPersistor
	AnimalPhotoPersistor
	AnimalVideoPersistor
	AnimalTypePersistor
	AnimalSpeciesPersistor
	BreedPersistor

	GetAnimalMinimalByID(ctx context.Context, animalID int64) (models.Animal, error)
	ListAnimals(ctx context.Context, filters dbtype.ListAnimalsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error)
	GetAnimalByID(ctx context.Context, animalID int64, filters dbtype.GetAnimalByIDFilters) (dbtype.AnimalWithJoinData, error)
	CreateAnimal(ctx context.Context, in dbtype.AnimalCreateSetter) (models.Animal, error)
	UpdateAnimal(ctx context.Context, animalID int64, in dbtype.AnimalUpdateSetter) (models.Animal, error)
	DeleteAnimalByID(ctx context.Context, animalID int64) (int64, error)
	DeleteAnimals(ctx context.Context, animalIDs []int64) error
	LikeAnimal(ctx context.Context, userID uuid.UUID, animalID int64) error
	UnlikeAnimal(ctx context.Context, userID uuid.UUID, animalID int64) error
	ApplyForAdoption(ctx context.Context, animalID int64, userID uuid.UUID, organizationID *int64) (models.Adoption, error)
	ApproveAdoption(ctx context.Context, adoptionID int64) error
	RejectAdoption(ctx context.Context, adoptionID int64) error
	ApproveAnimal(ctx context.Context, animalID int64) error
	RejectAnimal(ctx context.Context, animalID int64) error
}

type AnimalBreedPersistor interface {
	GetAnimalBreed(ctx context.Context, animalID, breedID int64) (dbtype.AnimalBreedWithJoinData, error)
	ListAnimalBreeds(ctx context.Context, animalID int64, filters dbtype.ListAnimalBreedsFilters) (dbtype.PagedResult[dbtype.AnimalBreedWithJoinData], error)
	UpdateAnimalBreed(ctx context.Context, animalID, breedID int64, in models.AnimalBreedSetter) (models.AnimalBreed, error)
	DeleteAnimalBreed(ctx context.Context, animalID, breedID int64) (int64, error)
	CreateAnimalBreeds(ctx context.Context, animalID int64, in []models.AnimalBreedSetter) ([]models.AnimalBreed, error)
	DeleteAnimalBreeds(ctx context.Context, animalID int64, breedIDs []int64) error
}

type AnimalTagPersistor interface {
	GetAnimalTag(ctx context.Context, animalID, tagID int64) (models.AnimalTag, error)
	ListAnimalTags(ctx context.Context, animalID int64, filters dbtype.ListAnimalTagsFilters) (dbtype.PagedResult[models.AnimalTag], error)
	UpdateAnimalTag(ctx context.Context, animalID, tagID int64, in models.AnimalTagSetter) (models.AnimalTag, error)
	DeleteAnimalTag(ctx context.Context, animalID, tagID int64) (int64, error)
	CreateAnimalTags(ctx context.Context, animalID int64, in []models.AnimalTagSetter) ([]models.AnimalTag, error)
	DeleteAnimalTags(ctx context.Context, animalID int64, tagIDs []int64) error
}

type AnimalPhotoPersistor interface {
	GetAnimalPhoto(ctx context.Context, animalID, photoID int64) (models.AnimalPhoto, error)
	ListAnimalPhotos(ctx context.Context, animalID int64, filters dbtype.ListAnimalPhotosFilters) (dbtype.PagedResult[models.AnimalPhoto], error)
	UpdateAnimalPhoto(ctx context.Context, animalID, photoID int64, in models.AnimalPhotoSetter) (models.AnimalPhoto, error)
	DeleteAnimalPhoto(ctx context.Context, animalID, photoID int64) (int64, error)
	CreateAnimalPhotos(ctx context.Context, animalID int64, in []models.AnimalPhotoSetter) ([]models.AnimalPhoto, error)
	DeleteAnimalPhotos(ctx context.Context, animalID int64, photoIDs []int64) error
}

type AnimalVideoPersistor interface {
	GetAnimalVideo(ctx context.Context, animalID, videoID int64) (models.AnimalVideo, error)
	ListAnimalVideos(ctx context.Context, animalID int64, filters dbtype.ListAnimalVideosFilters) (dbtype.PagedResult[models.AnimalVideo], error)
	UpdateAnimalVideo(ctx context.Context, animalID, videoID int64, in models.AnimalVideoSetter) (models.AnimalVideo, error)
	DeleteAnimalVideo(ctx context.Context, animalID, videoID int64) (int64, error)
	CreateAnimalVideos(ctx context.Context, animalID int64, in []models.AnimalVideoSetter) ([]models.AnimalVideo, error)
	DeleteAnimalVideos(ctx context.Context, animalID int64, videoIds []int64) error
}

type AnimalTypePersistor interface {
	ListAnimalTypes(ctx context.Context, filters dbtype.ListAnimalTypesFilters) (dbtype.PagedResult[models.AnimalType], error)
	GetAnimalTypeByID(ctx context.Context, animalTypeID int64) (models.AnimalType, error)
	CreateAnimalType(ctx context.Context, in models.AnimalTypeSetter) (models.AnimalType, error)
	UpdateAnimalType(ctx context.Context, animalTypeID int64, in models.AnimalTypeSetter) (models.AnimalType, error)
	DeleteAnimalTypeByID(ctx context.Context, animalTypeID int64) (int64, error)
	DeleteAnimalTypes(ctx context.Context, animalTypeIDs []int64) error
}

type AnimalSpeciesPersistor interface {
	ListAnimalSpecies(ctx context.Context, filters dbtype.ListAnimalSpeciesFilters) (dbtype.PagedResult[models.AnimalSpecie], error)
	GetAnimalSpecieByID(ctx context.Context, animalSpecieID int64) (models.AnimalSpecie, error)
	CreateAnimalSpecie(ctx context.Context, in models.AnimalSpecieSetter) (models.AnimalSpecie, error)
	UpdateAnimalSpecie(ctx context.Context, animalSpecieID int64, in models.AnimalSpecieSetter) (models.AnimalSpecie, error)
	DeleteAnimalSpecieByID(ctx context.Context, animalSpecieID int64) (int64, error)
	DeleteAnimalSpecies(ctx context.Context, animalSpecieIDs []int64) error
}

type BreedPersistor interface {
	ListBreeds(ctx context.Context, filters dbtype.ListBreedsFilters) (dbtype.PagedResult[models.Breed], error)
	GetBreedByID(ctx context.Context, breedID int64) (models.Breed, error)
	CreateBreed(ctx context.Context, in models.BreedSetter) (models.Breed, error)
	UpdateBreed(ctx context.Context, breedID int64, in models.BreedSetter) (models.Breed, error)
	DeleteBreedByID(ctx context.Context, breedID int64) (int64, error)
	DeleteBreeds(ctx context.Context, breedIDs []int64) error
}

type GeocodingResultPersistor interface {
	ListGeocodingResults(ctx context.Context, filters dbtype.ListGeocodingResultsFilters) (dbtype.PagedResult[models.GeocodingResult], error)
	GetGeocodingResultByID(ctx context.Context, geoResID int64) (models.GeocodingResult, error)
	GetGeocodingResultByQuery(ctx context.Context, query string) (models.GeocodingResult, error)
	CreateGeocodingResult(ctx context.Context, in models.GeocodingResultSetter) (models.GeocodingResult, error)
	UpdateGeocodingResult(ctx context.Context, geoResID int64, in models.GeocodingResultSetter) (models.GeocodingResult, error)
	DeleteGeocodingResultByID(ctx context.Context, geoResID int64) (int64, error)
	DeleteGeocodingResults(ctx context.Context, geoResIDs []int64) error
}

type AnalyticsPersistor interface {
	GetAnalyticsStats(ctx context.Context) (dbcustom.AnalyticsStats, error)
	GetMyAnalyticsStats(ctx context.Context, userID uuid.UUID) (dbcustom.MyAnalyticsStats, error)
}

type Persistor interface {
	User() UserPersistor
	Organization() OrganizationPersistor
	Country() CountryPersistor
	Animal() AnimalPersistor
	Geocoding() GeocodingResultPersistor
	Analytics() AnalyticsPersistor
	Adoption() AdoptionPersistor
}
