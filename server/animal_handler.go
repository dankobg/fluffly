package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/ptr"
	"github.com/oapi-codegen/nullable"
)

func (a *ApiHandler) CreateAnimal(ctx context.Context, request api.CreateAnimalRequestObject) (api.CreateAnimalResponseObject, error) {
	var animalCreateSetter persistence.AnimalCreateSetter

	animalCreateSetter.Animal = persistence.AnimalSetter{
		TypeID:    nullable.NewNullableWithValue(request.Body.AnimalTypeID),
		SpeciesID: nullable.NewNullableWithValue(request.Body.AnimalSpeciesID),
		Name:      nullable.NewNullableWithValue(request.Body.Name),
		Age:       nullable.NewNullableWithValue(string(request.Body.Age)),
		Size:      nullable.NewNullableWithValue(string(request.Body.Size)),
		ImageURL:  nullable.NewNullableWithValue(request.Body.ImageURL),
	}
	if request.Body.Gender != nil {
		gender := model.Gender_M
		if string(*request.Body.Gender) == model.Gender_F.String() {
			gender = model.Gender_F
		}
		animalCreateSetter.Animal.Gender = nullable.NewNullableWithValue(gender)
	}
	if request.Body.UserID != nil {
		animalCreateSetter.Animal.UserID = nullable.NewNullableWithValue(*request.Body.UserID)
	}
	if request.Body.OrganizationID != nil {
		animalCreateSetter.Animal.OrganizationID = nullable.NewNullableWithValue(*request.Body.OrganizationID)
	}
	if request.Body.Hermaphrodite != nil {
		animalCreateSetter.Animal.Hermaphrodite = nullable.NewNullableWithValue(*request.Body.Hermaphrodite)
	}
	if request.Body.Description != nil {
		animalCreateSetter.Animal.Description = nullable.NewNullableWithValue(*request.Body.Description)
	}
	if request.Body.Attributes != nil {
		animalCreateSetter.Animal.Attributes = nullable.NewNullableWithValue(*request.Body.Attributes)
	}

	if request.Body.Breeds != nil {
		animalBreedSetters := make([]persistence.AnimalBreedSetter, 0)
		for _, breed := range *request.Body.Breeds {
			animalBreedSetter := persistence.AnimalBreedSetter{
				BreedID: nullable.NewNullableWithValue(breed.BreedID),
			}
			if breed.Primary != nil {
				animalBreedSetter.Primary = nullable.NewNullableWithValue(*breed.Primary)
			}
			animalBreedSetters = append(animalBreedSetters, animalBreedSetter)
		}
		animalCreateSetter.Breeds = nullable.NewNullableWithValue(animalBreedSetters)
	}

	if request.Body.Microchip != nil {
		microchipSetter := persistence.MicrochipSetter{
			Number: nullable.NewNullableWithValue(request.Body.Microchip.Number),
		}
		if request.Body.Microchip.Brand != nil {
			microchipSetter.Brand = nullable.NewNullableWithValue(*request.Body.Microchip.Brand)
		}
		if request.Body.Microchip.Description != nil {
			microchipSetter.Description = nullable.NewNullableWithValue(*request.Body.Microchip.Description)
		}
		if request.Body.Microchip.Location != nil {
			microchipSetter.Location = nullable.NewNullableWithValue(*request.Body.Microchip.Location)
		}
		animalCreateSetter.Microchip = nullable.NewNullableWithValue(microchipSetter)
	}

	if request.Body.Tags != nil {
		animalTagSetters := make([]persistence.AnimalTagSetter, 0)
		for _, tag := range *request.Body.Tags {
			tagSetter := persistence.AnimalTagSetter{}
			if tag != "" {
				tagSetter.Name = nullable.NewNullableWithValue(tag)
			}
			animalTagSetters = append(animalTagSetters, tagSetter)
		}
		animalCreateSetter.Tags = nullable.NewNullableWithValue(animalTagSetters)
	}

	if request.Body.Photos != nil {
		animalPhotoSetters := make([]persistence.AnimalPhotoSetter, 0)
		for _, photo := range *request.Body.Photos {
			photoSetter := persistence.AnimalPhotoSetter{}
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
			animalPhotoSetters = append(animalPhotoSetters, photoSetter)
		}
		animalCreateSetter.Photos = nullable.NewNullableWithValue(animalPhotoSetters)
	}

	if request.Body.Videos != nil {
		animalVideoSetters := make([]persistence.AnimalVideoSetter, 0)
		for _, video := range *request.Body.Videos {
			videoSetter := persistence.AnimalVideoSetter{}
			if video.URL != "" {
				videoSetter.URL = nullable.NewNullableWithValue(video.URL)
			}
			animalVideoSetters = append(animalVideoSetters, videoSetter)
		}
		animalCreateSetter.Videos = nullable.NewNullableWithValue(animalVideoSetters)
	}

	animal, err := a.persistor.Animal().CreateAnimal(ctx, animalCreateSetter)
	if err != nil {
		msg := "could not create an animal"
		var e1 postgres.ErrAnimalUniqueViolation
		if errors.As(err, &e1) {
			msg += ", duplicate " + e1.Name
			return api.CreateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			msg += ", animal integrity error"
			return api.CreateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		return nil, fmt.Errorf("failed to create an animal")
	}
	resp := api.CreateAnimal201JSONResponse(dto.AnimalToResponse(animal))
	return resp, nil
}

func (a *ApiHandler) UpdateAnimal(ctx context.Context, request api.UpdateAnimalRequestObject) (api.UpdateAnimalResponseObject, error) {
	animalSetter := persistence.AnimalSetter{
		UserID:         request.Body.UserID,
		OrganizationID: request.Body.OrganizationID,
		Hermaphrodite:  request.Body.Hermaphrodite,
		Description:    request.Body.Description,
		Attributes:     request.Body.Attributes,
		// TypeID:          request.Body.AnimalTypeID,
		// SpeciesID:       request.Body.AnimalSpeciesID,
		// Name:            request.Body.Name,
		// Gender:          request.Body.Gender,
		// Age:             request.Body.Age,
		// Size:            request.Body.Size,
		// ImageURL:        request.Body.ImageURL,
		// Distance:        request.Body.Distance,
		// Status:          request.Body.Status,
		// StatusChangedAt: request.Body.StatusChangedAt,
		// AdoptedAt:       request.Body.AdoptedAt,
	}

	animal, err := a.persistor.Animal().UpdateAnimal(ctx, request.ID, animalSetter)
	if err != nil {
		msg := "could not update an animal"
		var e1 postgres.ErrAnimalUniqueViolation
		if errors.As(err, &e1) {
			msg += ", duplicate " + e1.Name
			return api.UpdateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			msg += ", animal integrity error"
			return api.UpdateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		return nil, fmt.Errorf("failed to update an animal")
	}
	resp := api.UpdateAnimal201JSONResponse(dto.AnimalToResponse(animal))
	return resp, nil
}

func (a *ApiHandler) DeleteAnimal(ctx context.Context, request api.DeleteAnimalRequestObject) (api.DeleteAnimalResponseObject, error) {
	_, err := a.persistor.Animal().DeleteAnimalByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an animal by id: %w", err)
	}
	resp := api.DeleteAnimal204Response{}
	return resp, nil
}

func (a *ApiHandler) ListAnimals(ctx context.Context, request api.ListAnimalsRequestObject) (api.ListAnimalsResponseObject, error) {
	var filters persistence.AnimalFilters
	filters.Pagination = ptr.Of(getPaginationParams(request.Params.Page, request.Params.PageSize))
	animals, err := a.persistor.Animal().ListAnimals(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animals: %w", err)
	}
	animalsData := make([]api.Animal, len(animals.Data))
	for i, animalWithJoinData := range animals.Data {
		animalsData[i] = dto.AnimalWithJoinDataToResponse(animalWithJoinData)
	}
	resp := api.ListAnimals200JSONResponse{
		Data: animalsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, animals.TotalCount),
	}
	return resp, nil
}

func (a *ApiHandler) GetAnimal(ctx context.Context, request api.GetAnimalRequestObject) (api.GetAnimalResponseObject, error) {
	animalWithJoinData, err := a.persistor.Animal().GetAnimalByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalNotFound) {
			return api.GetAnimal404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Animal not found"}}, nil
		}
		return nil, fmt.Errorf("failed to get an animal by id: %w", err)
	}
	resp := api.GetAnimal200JSONResponse(dto.AnimalWithJoinDataToResponse(animalWithJoinData))
	return resp, nil
}
