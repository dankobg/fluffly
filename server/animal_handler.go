package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/ptr"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
)

func (a *ApiHandler) CreateAnimal(ctx context.Context, request api.CreateAnimalRequestObject) (api.CreateAnimalResponseObject, error) {
	form, err := request.Body.ReadForm(50 << 20)
	if err != nil {
		return api.CreateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: 400, Message: err.Error()}}, nil
	}
	defer form.RemoveAll()
	animalData := form.Value["data"][0]
	var input api.CreateAnimalBody
	if err := json.Unmarshal([]byte(animalData), &input); err != nil {
		return api.CreateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: 400, Message: err.Error()}}, nil
	}

	mainImageFileHeaders := form.File["image"]
	if len(mainImageFileHeaders) == 0 && (input.ImageURL.IsNull() || (!input.ImageURL.IsNull() && input.ImageURL.MustGet() == "")) {
		return api.CreateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: 400, Message: "Provide image file or image_url as main animal photo"}}, nil
	}
	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]
	var mainImageFileSources []uploadSource
	var photoFileSources []uploadSource
	var videoFileSources []uploadSource
	for _, fh := range mainImageFileHeaders {
		mainImageFileSources = append(mainImageFileSources, multipartSource{fh: fh})
	}
	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}
	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
	}
	if input.ImageURL.IsSpecified() && !input.ImageURL.IsNull() {
		mainImageFileSources = append(mainImageFileSources, urlSource{c: a.httpc, url: input.ImageURL.MustGet()})
	}
	if input.Photos.IsSpecified() && !input.Photos.IsNull() {
		for _, photo := range input.Photos.MustGet() {
			photoFileSources = append(photoFileSources, urlSource{c: a.httpc, url: photo.URL})
		}
	}
	if input.Videos.IsSpecified() && !input.Videos.IsNull() {
		for _, video := range input.Videos.MustGet() {
			videoFileSources = append(videoFileSources, urlSource{c: a.httpc, url: video.URL})
		}
	}
	var mainImageUploadResults []uploadResult
	var photoUploadResults []uploadResult
	var videoUploadResults []uploadResult
	if len(mainImageFileSources) > 0 {
		mainImageResults := a.uploadAnimalFiles(ctx, mainImageFileSources, 1)
		mainImageUploadResults = append(mainImageUploadResults, mainImageResults...)
	}
	if len(photoFileSources) > 0 {
		photoResults := a.uploadAnimalFiles(ctx, photoFileSources, 5)
		photoUploadResults = append(photoUploadResults, photoResults...)
	}
	if len(videoFileSources) > 0 {
		videoResults := a.uploadAnimalFiles(ctx, videoFileSources, 5)
		videoUploadResults = append(videoUploadResults, videoResults...)
	}
	var filesToDelete []string
	for _, res := range append(mainImageUploadResults, append(photoUploadResults, videoUploadResults...)...) {
		if res.Err != nil {
			filesToDelete = append(filesToDelete, res.Name)
		}
	}
	if len(filesToDelete) > 0 {
		go func() {
			_ = a.deleteUploadedFiles(ctx, filesToDelete, 5)
		}()
		return api.CreateAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: 400, Message: "failed to upload files"}}, nil
	}

	var animalCreateSetter dbtype.AnimalCreateSetter

	animalCreateSetter.Animal = dbtype.AnimalSetter{
		TypeID:               nullable.NewNullableWithValue(input.AnimalTypeID),
		SpeciesID:            nullable.NewNullableWithValue(input.AnimalSpeciesID),
		Name:                 nullable.NewNullableWithValue(input.Name),
		Age:                  nullable.NewNullableWithValue(string(input.Age)),
		Size:                 nullable.NewNullableWithValue(string(input.Size)),
		ImageObjectKind:      nullable.NewNullableWithValue(a.uploader.Kind()),
		ImageObjectRefSmall:  nullable.NewNullableWithValue(mainImageUploadResults[0].Name),
		ImageObjectRefMedium: nullable.NewNullableWithValue(mainImageUploadResults[0].Name),
		ImageObjectRefLarge:  nullable.NewNullableWithValue(mainImageUploadResults[0].Name),
		ImageObjectRefFull:   nullable.NewNullableWithValue(mainImageUploadResults[0].Name),
		UserID:               input.UserID,
		OrganizationID:       input.OrganizationID,
		Hermaphrodite:        input.Hermaphrodite,
		Description:          input.Description,
		Properties:           input.Properties,
	}
	if input.Gender.IsSpecified() && !input.Gender.IsNull() {
		gender := model.Gender_M
		if string(input.Gender.MustGet()) == model.Gender_F.String() {
			gender = model.Gender_F
		}
		animalCreateSetter.Animal.Gender = nullable.NewNullableWithValue(gender)
	}

	if input.Breeds.IsSpecified() && !input.Breeds.IsNull() {
		animalBreedSetters := make([]dbtype.AnimalBreedSetter, 0)
		for _, breed := range input.Breeds.MustGet() {
			animalBreedSetter := dbtype.AnimalBreedSetter{
				BreedID: nullable.NewNullableWithValue(breed.BreedID),
			}
			if breed.Primary != nil {
				animalBreedSetter.Primary = nullable.NewNullableWithValue(*breed.Primary)
			}
			animalBreedSetters = append(animalBreedSetters, animalBreedSetter)
		}
		animalCreateSetter.Breeds = nullable.NewNullableWithValue(animalBreedSetters)
	}

	if input.Microchip != nil {
		microchipSetter := dbtype.MicrochipSetter{
			Number: nullable.NewNullableWithValue(input.Microchip.Number),
		}
		if input.Microchip.Brand != nil {
			microchipSetter.Brand = nullable.NewNullableWithValue(*input.Microchip.Brand)
		}
		if input.Microchip.Description != nil {
			microchipSetter.Description = nullable.NewNullableWithValue(*input.Microchip.Description)
		}
		if input.Microchip.Location != nil {
			microchipSetter.Location = nullable.NewNullableWithValue(*input.Microchip.Location)
		}
		animalCreateSetter.Microchip = nullable.NewNullableWithValue(microchipSetter)
	}

	if input.Tags.IsSpecified() && !input.Tags.IsNull() {
		animalTagSetters := make([]dbtype.AnimalTagSetter, 0)
		for _, tag := range input.Tags.MustGet() {
			animalTagSetters = append(animalTagSetters, dbtype.AnimalTagSetter{
				Name: nullable.NewNullableWithValue(tag),
			})
		}
		animalCreateSetter.Tags = nullable.NewNullableWithValue(animalTagSetters)
	}

	if len(photoUploadResults) > 0 {
		animalPhotoSetters := make([]dbtype.AnimalPhotoSetter, len(photoUploadResults))
		for i, photoRes := range photoUploadResults {
			// @TODO: optimize image size
			photoSetter := dbtype.AnimalPhotoSetter{
				ObjectKind:      nullable.NewNullableWithValue(a.uploader.Kind()),
				ObjectRefSmall:  nullable.NewNullableWithValue(photoRes.Name),
				ObjectRefMedium: nullable.NewNullableWithValue(photoRes.Name),
				ObjectRefLarge:  nullable.NewNullableWithValue(photoRes.Name),
				ObjectRefFull:   nullable.NewNullableWithValue(photoRes.Name),
			}
			animalPhotoSetters[i] = photoSetter
		}
		animalCreateSetter.Photos = nullable.NewNullableWithValue(animalPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		animalVideoSetters := make([]dbtype.AnimalVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			videoSetter := dbtype.AnimalVideoSetter{
				ObjectKind: nullable.NewNullableWithValue(a.uploader.Kind()),
				ObjectRef:  nullable.NewNullableWithValue(videoRes.Name),
			}
			animalVideoSetters[i] = videoSetter
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
	animalSetter := dbtype.AnimalSetter{
		UserID:         request.Body.UserID,
		OrganizationID: request.Body.OrganizationID,
		Hermaphrodite:  request.Body.Hermaphrodite,
		Description:    request.Body.Description,
		// Properties:     request.Body.Properties,
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
	var filters dbtype.AnimalFilters
	filters.Pagination = ptr.Of(getPaginationParams(request.Params.Page, request.Params.PageSize))
	animals, err := a.persistor.Animal().ListAnimals(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animals: %w", err)
	}
	animalsData := make([]api.Animal, len(animals.Data))
	for i, animalWithJoinData := range animals.Data {
		animalsData[i] = dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader)
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
	resp := api.GetAnimal200JSONResponse(dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader))
	return resp, nil
}

func (a *ApiHandler) LikeAnimal(ctx context.Context, request api.LikeAnimalRequestObject) (api.LikeAnimalResponseObject, error) {
	// @TODO: remove hardcoded value
	userID := uuid.MustParse("6e482ec1-64c4-4f34-93ba-392f4e473444")
	if err := a.persistor.Animal().LikeAnimal(ctx, userID, request.ID); err != nil {
		msg := "failed to like an animal"
		var e1 postgres.ErrAnimalUniqueViolation
		if errors.As(err, &e1) {
			msg += ", duplicate " + e1.Name
			return api.LikeAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			msg += ", animal integrity error"
			return api.LikeAnimal400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: msg}}, nil
		}
		return nil, fmt.Errorf("failed to like an animal")
	}
	return api.EmptyResponseResponse{}, nil
}

func (a *ApiHandler) UnlikeAnimal(ctx context.Context, request api.UnlikeAnimalRequestObject) (api.UnlikeAnimalResponseObject, error) {
	// @TODO: remove hardcoded value
	userID := uuid.MustParse("6e482ec1-64c4-4f34-93ba-392f4e473444")
	if err := a.persistor.Animal().UnlikeAnimal(ctx, userID, request.ID); err != nil {
		return nil, fmt.Errorf("failed to unlike an animal: %w", err)
	}
	return api.EmptyResponseResponse{}, nil
}
