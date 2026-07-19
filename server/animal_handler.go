package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/convert"
	"github.com/dankobg/fluffly/db/gen/enums"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/geocoding/nominatim"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	"github.com/google/uuid"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"github.com/stephenafamo/bob/types"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

const (
	createAnimalFileMaxMemory = 100 << 20
)

func (a *ApiHandler) SubmitAnimal(ctx context.Context, request api.SubmitAnimalRequestObject) (api.SubmitAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "submit",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.SubmitAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(createAnimalFileMaxMemory)
	if err != nil {
		return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	animalDocData := form.Value["data"]
	animalFileData := form.File["data"]

	var animalDataBytes []byte

	if len(animalDocData) == 0 && len(animalFileData) == 0 {
		return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide animal json data document")}, nil
	}

	if len(animalDocData) > 0 {
		animalDataBytes = []byte(animalDocData[0])
	}

	if len(animalFileData) > 0 {
		orgFile, err := animalFileData[0].Open()
		if err != nil {
			return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		animalDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.SubmitAnimalBody
	if err := json.Unmarshal(animalDataBytes, &input); err != nil {
		return api.SubmitAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_internal", "failed to unmarshal non file data")}, nil
	}

	mainImageFileHeaders := form.File["image"]
	if len(mainImageFileHeaders) == 0 && (input.ImageURL == nil || (input.ImageURL != nil && *input.ImageURL == "")) {
		return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_image_missing", "Provide image file or image_url for main animal photo", err.Error())}, nil
	}

	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]

	var (
		mainImageFileSources []uploadSource
		photoFileSources     []uploadSource
		videoFileSources     []uploadSource
	)

	for _, fh := range mainImageFileHeaders {
		mainImageFileSources = append(mainImageFileSources, multipartSource{fh: fh})
	}

	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
	}

	if input.ImageURL != nil {
		mainImageFileSources = append(mainImageFileSources, urlSource{c: a.httpc, url: *input.ImageURL})
	}

	if input.Photos != nil {
		for _, photo := range *input.Photos {
			photoFileSources = append(photoFileSources, urlSource{c: a.httpc, url: photo.URL})
		}
	}

	if input.Videos != nil {
		for _, video := range *input.Videos {
			videoFileSources = append(videoFileSources, urlSource{c: a.httpc, url: video.URL})
		}
	}

	var (
		mainImageUploadResults []uploadResult
		photoUploadResults     []uploadResult
		videoUploadResults     []uploadResult
	)

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
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	var animalSubmitSetter dbtype.AnimalCreateSetter

	var properties omitnull.Val[types.JSON[dbcustom.Properties]]
	if input.Properties != nil {
		properties.Set(types.NewJSON(*input.Properties))
	}

	animalSetter := models.AnimalSetter{
		UserID:               omit.FromPtr(input.UserID),
		OrganizationID:       omitnull.FromPtr(input.OrganizationID),
		Hermaphrodite:        omit.FromPtr(input.Hermaphrodite),
		Description:          omitnull.FromPtr(input.Description),
		Properties:           properties,
		TypeID:               omit.From(input.AnimalTypeID),
		SpecieID:             omit.From(input.AnimalSpecieID),
		Name:                 omit.From(input.Name),
		Age:                  omit.From(string(input.Age)),
		Size:                 omit.From(string(input.Size)),
		ImageObjectKind:      omit.From(a.uploader.Kind()),
		ImageObjectRefSmall:  omit.From(mainImageUploadResults[0].Name),
		ImageObjectRefMedium: omit.From(mainImageUploadResults[0].Name),
		ImageObjectRefLarge:  omit.From(mainImageUploadResults[0].Name),
		ImageObjectRefFull:   omit.From(mainImageUploadResults[0].Name),
	}
	if input.Gender != nil {
		gender := enums.GenderM
		if string(*input.Gender) == enums.GenderF.String() {
			gender = enums.GenderF
		}

		animalSetter.Gender = omitnull.From(gender)
	}

	animalSubmitSetter.Animal = animalSetter

	if input.AnimalBreeds != nil {
		animalBreedSetters := make([]models.AnimalBreedSetter, 0)
		for _, breed := range *input.AnimalBreeds {
			animalBreedSetters = append(animalBreedSetters, models.AnimalBreedSetter{
				BreedID: omit.From(breed.BreedID),
				Primary: omit.FromPtr(breed.Primary),
			})
		}

		animalSubmitSetter.Breeds = omitnull.From(animalBreedSetters)
	}

	if input.Microchip != nil {
		animalSubmitSetter.Microchip = omitnull.From(models.MicrochipSetter{
			Number:      omit.From(input.Microchip.Number),
			Brand:       omitnull.FromPtr(input.Microchip.Brand),
			Description: omitnull.FromPtr(input.Microchip.Description),
			Location:    omitnull.FromPtr(input.Microchip.Location),
		})
	}

	if input.Tags != nil {
		animalTagSetters := make([]models.AnimalTagSetter, 0)
		for _, tag := range *input.Tags {
			animalTagSetters = append(animalTagSetters, models.AnimalTagSetter{
				Name: omit.From(tag),
			})
		}

		animalSubmitSetter.Tags = omitnull.From(animalTagSetters)
	}

	if len(photoUploadResults) > 0 {
		animalPhotoSetters := make([]models.AnimalPhotoSetter, len(photoUploadResults))
		for i, photoRes := range photoUploadResults {
			// @TODO: optimize image size
			photoSetter := models.AnimalPhotoSetter{
				ObjectKind:      omit.From(a.uploader.Kind()),
				ObjectRefSmall:  omitnull.From(photoRes.Name),
				ObjectRefMedium: omitnull.From(photoRes.Name),
				ObjectRefLarge:  omitnull.From(photoRes.Name),
				ObjectRefFull:   omitnull.From(photoRes.Name),
			}
			animalPhotoSetters[i] = photoSetter
		}

		animalSubmitSetter.Photos = omitnull.From(animalPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		animalVideoSetters := make([]models.AnimalVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			animalVideoSetters[i] = models.AnimalVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}

		animalSubmitSetter.Videos = omitnull.From(animalVideoSetters)
	}

	animal, err := a.persistor.Animal().CreateAnimal(ctx, animalSubmitSetter)
	if err != nil {
		msg := "could not submit an animal"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.SubmitAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_save", msg, reason)}, nil
		}

		return api.SubmitAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_save", msg, reason)}, nil
	}

	resp := api.SubmitAnimal201JSONResponse(dto.AnimalToResponse(animal, a.uploader))

	// @TODO: use outbox pattern
	if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
		a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
		return api.SubmitAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) CreateAnimal(ctx context.Context, request api.CreateAnimalRequestObject) (api.CreateAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(createAnimalFileMaxMemory)
	if err != nil {
		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	animalDocData := form.Value["data"]
	animalFileData := form.File["data"]

	var animalDataBytes []byte

	if len(animalDocData) == 0 && len(animalFileData) == 0 {
		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide animal json data document")}, nil
	}

	if len(animalDocData) > 0 {
		animalDataBytes = []byte(animalDocData[0])
	}

	if len(animalFileData) > 0 {
		orgFile, err := animalFileData[0].Open()
		if err != nil {
			return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		animalDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.CreateAnimalBody
	if err := json.Unmarshal(animalDataBytes, &input); err != nil {
		return api.CreateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_internal", "failed to unmarshal non file data")}, nil
	}

	mainImageFileHeaders := form.File["image"]
	if len(mainImageFileHeaders) == 0 && (input.ImageURL == nil || (input.ImageURL != nil && *input.ImageURL == "")) {
		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_image_missing", "Provide image file or image_url for main animal photo", err.Error())}, nil
	}

	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]

	var (
		mainImageFileSources []uploadSource
		photoFileSources     []uploadSource
		videoFileSources     []uploadSource
	)

	for _, fh := range mainImageFileHeaders {
		mainImageFileSources = append(mainImageFileSources, multipartSource{fh: fh})
	}

	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
	}

	if input.ImageURL != nil {
		mainImageFileSources = append(mainImageFileSources, urlSource{c: a.httpc, url: *input.ImageURL})
	}

	if input.Photos != nil {
		for _, photo := range *input.Photos {
			photoFileSources = append(photoFileSources, urlSource{c: a.httpc, url: photo.URL})
		}
	}

	if input.Videos != nil {
		for _, video := range *input.Videos {
			videoFileSources = append(videoFileSources, urlSource{c: a.httpc, url: video.URL})
		}
	}

	var (
		mainImageUploadResults []uploadResult
		photoUploadResults     []uploadResult
		videoUploadResults     []uploadResult
	)

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
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	var animalCreateSetter dbtype.AnimalCreateSetter

	var properties omitnull.Val[types.JSON[dbcustom.Properties]]
	if input.Properties != nil {
		properties.Set(types.NewJSON(*input.Properties))
	}

	animalSetter := models.AnimalSetter{
		UserID:               omit.FromPtr(input.UserID),
		OrganizationID:       omitnull.FromPtr(input.OrganizationID),
		Hermaphrodite:        omit.FromPtr(input.Hermaphrodite),
		Description:          omitnull.FromPtr(input.Description),
		Properties:           properties,
		TypeID:               omit.From(input.AnimalTypeID),
		SpecieID:             omit.From(input.AnimalSpecieID),
		Name:                 omit.From(input.Name),
		Age:                  omit.From(string(input.Age)),
		Size:                 omit.From(string(input.Size)),
		Status:               omit.FromPtr((*string)(input.Status)),
		ImageObjectKind:      omit.From(a.uploader.Kind()),
		ImageObjectRefSmall:  omit.From(mainImageUploadResults[0].Name),
		ImageObjectRefMedium: omit.From(mainImageUploadResults[0].Name),
		ImageObjectRefLarge:  omit.From(mainImageUploadResults[0].Name),
		ImageObjectRefFull:   omit.From(mainImageUploadResults[0].Name),
	}
	if input.Gender != nil {
		gender := enums.GenderM
		if string(*input.Gender) == enums.GenderF.String() {
			gender = enums.GenderF
		}

		animalSetter.Gender = omitnull.From(gender)
	}

	animalCreateSetter.Animal = animalSetter

	if input.AnimalBreeds != nil {
		animalBreedSetters := make([]models.AnimalBreedSetter, 0)
		for _, breed := range *input.AnimalBreeds {
			animalBreedSetters = append(animalBreedSetters, models.AnimalBreedSetter{
				BreedID: omit.From(breed.BreedID),
				Primary: omit.FromPtr(breed.Primary),
			})
		}

		animalCreateSetter.Breeds = omitnull.From(animalBreedSetters)
	}

	if input.Microchip != nil {
		animalCreateSetter.Microchip = omitnull.From(models.MicrochipSetter{
			Number:      omit.From(input.Microchip.Number),
			Brand:       omitnull.FromPtr(input.Microchip.Brand),
			Description: omitnull.FromPtr(input.Microchip.Description),
			Location:    omitnull.FromPtr(input.Microchip.Location),
		})
	}

	if input.Tags != nil {
		animalTagSetters := make([]models.AnimalTagSetter, 0)
		for _, tag := range *input.Tags {
			animalTagSetters = append(animalTagSetters, models.AnimalTagSetter{
				Name: omit.From(tag),
			})
		}

		animalCreateSetter.Tags = omitnull.From(animalTagSetters)
	}

	if len(photoUploadResults) > 0 {
		animalPhotoSetters := make([]models.AnimalPhotoSetter, len(photoUploadResults))
		for i, photoRes := range photoUploadResults {
			// @TODO: optimize image size
			photoSetter := models.AnimalPhotoSetter{
				ObjectKind:      omit.From(a.uploader.Kind()),
				ObjectRefSmall:  omitnull.From(photoRes.Name),
				ObjectRefMedium: omitnull.From(photoRes.Name),
				ObjectRefLarge:  omitnull.From(photoRes.Name),
				ObjectRefFull:   omitnull.From(photoRes.Name),
			}
			animalPhotoSetters[i] = photoSetter
		}

		animalCreateSetter.Photos = omitnull.From(animalPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		animalVideoSetters := make([]models.AnimalVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			animalVideoSetters[i] = models.AnimalVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}

		animalCreateSetter.Videos = omitnull.From(animalVideoSetters)
	}

	animal, err := a.persistor.Animal().CreateAnimal(ctx, animalCreateSetter)
	if err != nil {
		msg := "could not create an animal"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_save", msg, reason)}, nil
		}

		return api.CreateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_save", msg, reason)}, nil
	}

	resp := api.CreateAnimal201JSONResponse(dto.AnimalToResponse(animal, a.uploader))

	// @TODO: use outbox pattern
	if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
		a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
		return api.CreateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) UpdateAnimal(ctx context.Context, request api.UpdateAnimalRequestObject) (api.UpdateAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	animalDocData := form.Value["data"]
	animalFileData := form.File["data"]

	var animalDataBytes []byte

	if len(animalDocData) == 0 && len(animalFileData) == 0 {
		return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide animal json data document")}, nil
	}

	if len(animalDocData) > 0 {
		animalDataBytes = []byte(animalDocData[0])
	}

	if len(animalFileData) > 0 {
		orgFile, err := animalFileData[0].Open()
		if err != nil {
			return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		animalDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.UpdateAnimalBody
	if err := json.Unmarshal(animalDataBytes, &input); err != nil {
		return api.UpdateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_internal", "failed to unmarshal non file data")}, nil
	}

	mainImageFileHeaders := form.File["image"]
	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]

	var (
		mainImageFileSources []uploadSource
		photoFileSources     []uploadSource
		videoFileSources     []uploadSource
	)

	for _, fh := range mainImageFileHeaders {
		mainImageFileSources = append(mainImageFileSources, multipartSource{fh: fh})
	}

	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
	}

	if input.ImageURL != nil {
		mainImageFileSources = append(mainImageFileSources, urlSource{c: a.httpc, url: *input.ImageURL})
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

	var (
		mainImageUploadResults []uploadResult
		photoUploadResults     []uploadResult
		videoUploadResults     []uploadResult
	)

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
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	var properties omitnull.Val[types.JSON[dbcustom.Properties]]
	if input.Properties.IsSpecified() && !input.Properties.IsNull() {
		properties.Set(types.NewJSON(input.Properties.MustGet()))
	}

	animalSetter := models.AnimalSetter{
		TypeID:         omit.FromPtr(input.AnimalTypeID),
		SpecieID:       omit.FromPtr(input.AnimalSpecieID),
		Name:           omit.FromPtr(input.Name),
		UserID:         convert.NullableToOmit(input.UserID),
		OrganizationID: convert.NullableToOmitNull(input.OrganizationID),
		Hermaphrodite:  convert.NullableToOmit(input.Hermaphrodite),
		Description:    convert.NullableToOmitNull(input.Description),
		Properties:     properties,
	}
	if len(mainImageUploadResults) > 0 {
		animalSetter.ImageObjectKind = omit.From(a.uploader.Kind())
		animalSetter.ImageObjectRefSmall = omit.From(mainImageUploadResults[0].Name)
		animalSetter.ImageObjectRefMedium = omit.From(mainImageUploadResults[0].Name)
		animalSetter.ImageObjectRefLarge = omit.From(mainImageUploadResults[0].Name)
		animalSetter.ImageObjectRefFull = omit.From(mainImageUploadResults[0].Name)
	}

	if input.Age.IsSpecified() && !input.Age.IsNull() {
		animalSetter.Age = omit.From(string(input.Age.MustGet()))
	}

	if input.Size.IsSpecified() && !input.Size.IsNull() {
		animalSetter.Size = omit.From(string(input.Size.MustGet()))
	}

	if input.Gender.IsSpecified() {
		if input.Gender.IsNull() {
			animalSetter.Gender = omitnull.FromPtr[enums.Gender](nil)
		} else {
			gender := enums.GenderM
			if string(input.Gender.MustGet()) == enums.GenderF.String() {
				gender = enums.GenderF
			}

			animalSetter.Gender = omitnull.From(gender)
		}
	}

	animalUpdateSetter := dbtype.AnimalUpdateSetter{
		Animal: omitnull.From(animalSetter),
	}

	if input.Microchip.IsSpecified() {
		if input.Microchip.IsNull() {
			animalUpdateSetter.Microchip = omitnull.FromPtr[models.MicrochipSetter](nil)
		} else {
			inMicrochip := input.Microchip.MustGet()
			animalUpdateSetter.Microchip = omitnull.From(models.MicrochipSetter{
				Number:      omit.FromPtr(inMicrochip.Number),
				Brand:       convert.NullableToOmitNull(inMicrochip.Brand),
				Description: convert.NullableToOmitNull(inMicrochip.Description),
				Location:    convert.NullableToOmitNull(inMicrochip.Location),
			})
		}
	}

	if len(photoUploadResults) > 0 {
		animalPhotoSetters := make([]models.AnimalPhotoSetter, len(photoUploadResults))
		for i, photoRes := range photoUploadResults {
			// @TODO: optimize image size
			animalPhotoSetters[i] = models.AnimalPhotoSetter{
				ObjectKind:      omit.From(a.uploader.Kind()),
				ObjectRefSmall:  omitnull.From(photoRes.Name),
				ObjectRefMedium: omitnull.From(photoRes.Name),
				ObjectRefLarge:  omitnull.From(photoRes.Name),
				ObjectRefFull:   omitnull.From(photoRes.Name),
			}
		}

		animalUpdateSetter.Photos = omitnull.From(animalPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		animalVideoSetters := make([]models.AnimalVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			animalVideoSetters[i] = models.AnimalVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}

		animalUpdateSetter.Videos = omitnull.From(animalVideoSetters)
	}

	if input.AnimalBreeds.IsSpecified() {
		if input.AnimalBreeds.IsNull() {
			animalUpdateSetter.Breeds = omitnull.FromPtr[[]models.AnimalBreedSetter](nil)
		} else {
			animalBreedsSetters := make([]models.AnimalBreedSetter, 0)

			if input.AnimalBreeds.IsSpecified() {
				for _, breed := range input.AnimalBreeds.MustGet() {
					animalBreedsSetters = append(animalBreedsSetters, models.AnimalBreedSetter{
						BreedID: omit.FromPtr(breed.BreedID),
						Primary: convert.NullableToOmit(breed.Primary),
					})
				}
			}

			animalUpdateSetter.Breeds = omitnull.From(animalBreedsSetters)
		}
	}

	if input.Tags.IsSpecified() {
		if input.Tags.IsNull() {
			// animalUpdateSetter.Tags = nullable.NewNullNullable[[]models.AnimalTagSetter]()
			animalUpdateSetter.Tags = omitnull.FromPtr[[]models.AnimalTagSetter](nil)
		} else {
			animalTagsSetters := make([]models.AnimalTagSetter, 0)

			if input.Tags.IsSpecified() {
				for _, tag := range input.Tags.MustGet() {
					animalTagsSetters = append(animalTagsSetters, models.AnimalTagSetter{
						Name: omit.FromPtr(tag.Name),
					})
				}
			}

			animalUpdateSetter.Tags = omitnull.From(animalTagsSetters)
		}
	}

	var oldMainImage string

	if len(mainImageUploadResults) > 0 {
		oldAnimalData, err := a.persistor.Animal().GetAnimalMinimalByID(ctx, request.ID)
		if err != nil {
			return api.UpdateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_get_old", "failed to get old animal")}, nil
		}

		oldMainImage = oldAnimalData.ImageObjectRefFull
	}

	animal, err := a.persistor.Animal().UpdateAnimal(ctx, request.ID, animalUpdateSetter)
	if err != nil {
		msg := "could not update an animal"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_update", msg, reason)}, nil
		}

		return api.UpdateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_edit", msg, reason)}, nil
	}

	resp := api.UpdateAnimal201JSONResponse(dto.AnimalToResponse(animal, a.uploader))

	go func(oldImage string) {
		if oldImage != "" {
			_ = a.deleteUploadedFiles(context.Background(), []string{oldImage}, 1)
		}
	}(oldMainImage)

	// @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.UpdateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteAnimal(ctx context.Context, request api.DeleteAnimalRequestObject) (api.DeleteAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	_, err := a.persistor.Animal().DeleteAnimalByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an animal by id: %w", err)
	}

	resp := api.DeleteAnimal204Response{}

	// @TODO: use outbox pattern
	if err := deleteAnimalRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", request.ID), slog.Any("error", err))
		return api.DeleteAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) DeleteAnimals(ctx context.Context, request api.DeleteAnimalsRequestObject) (api.DeleteAnimalsResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimals204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.Body.Ids[0]), // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimals403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animals_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().DeleteAnimals(ctx, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete animals by ids: %w", err)
	}

	resp := api.DeleteAnimals204Response{}

	// @TODO: use outbox pattern
	for _, id := range request.Body.Ids {
		if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
			a.Log.Error("failed to delete animals relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
			return api.DeleteAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animals_permissions", "failed to delete permissions")}, nil
		}
	}

	return resp, nil
}

func (a *ApiHandler) ListAnimals(ctx context.Context, request api.ListAnimalsRequestObject) (api.ListAnimalsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAnimalsdefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("animal_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalsFilters{ListAnimalsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page

	filters.PageSize = &paginationParams.PageSize
	if sessionValid(sess) {
		filters.UserID = new(uuid.MustParse(sess.Identity.Id))
	}

	if request.Params.Properties != nil && request.Params.AnimalSpecieID != nil && len(*request.Params.AnimalSpecieID) > 0 {
		specieSchema, err := a.LoadSpeciesPropertiesJsonSchema((*request.Params.AnimalSpecieID)[0])
		if err != nil {
			return api.ListAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("schema_error", "invalid json schema")}, nil
		}

		propertiesFilters, err := dbtype.CreateDynamicPropertiesFilter(specieSchema, *request.Params.Properties)
		if err != nil {
			return api.ListAnimalsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newUnauthorizedErr("filters_err", "failed to create dynamic filters")}, nil
		}

		filters.PropertiesFilters = propertiesFilters
	}

	if request.Params.Location != nil && (filters.Lat == nil || filters.Lon == nil) {
		normalizedSearchTerm := nominatim.NormalizeSearchQuery(*request.Params.Location)

		cachedGeocodingResult, err := a.persistor.Geocoding().GetGeocodingResultByQuery(ctx, normalizedSearchTerm)
		if err != nil {
			ctm, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			geocodingResult, err := a.geocoder.ForwardGeocode(ctm, normalizedSearchTerm)
			if err != nil {
				a.Log.Error("failed to forward geocode", slog.String("query", *request.Params.Location), slog.Any("error", err))

				filters.Lat = nil
				filters.Lon = nil
				filters.RadiusM = nil
			} else {
				filters.Lat = &geocodingResult.Lat
				filters.Lon = &geocodingResult.Lon

				point := geom.NewPoint(geom.XY).SetSRID(dbcustom.SRID).MustSetCoords(geom.Coord{geocodingResult.Lon, geocodingResult.Lat})
				coords := &ewkb.Point{Point: point}

				if _, err := a.persistor.Geocoding().CreateGeocodingResult(ctx, models.GeocodingResultSetter{
					Query:  omit.From(normalizedSearchTerm),
					Coords: omit.From(coords),
				}); err != nil {
					a.Log.Error("create geocoding result cache", slog.String("query", *request.Params.Location), slog.Any("error", err))
				}
			}
		} else {
			cachedGeocodingResult.Coords.FlatCoords()

			flatCoords := cachedGeocodingResult.Coords.FlatCoords()
			if len(flatCoords) == 2 {
				filters.Lon = &flatCoords[0]
				filters.Lat = &flatCoords[1]
			}
		}
	}

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
	sess := GetSession(ctx)

	filters := dbtype.GetAnimalByIDFilters{GetAnimalParams: request.Params}
	if sessionValid(sess) {
		filters.UserID = new(uuid.MustParse(sess.Identity.Id))
	}

	animalWithJoinData, err := a.persistor.Animal().GetAnimalByID(ctx, request.ID, filters)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalNotFound) {
			return api.GetAnimal404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_not_found", "animal not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an animal by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnimaldefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "animal_permission", "permission denied")}, nil
	}

	resp := api.GetAnimal200JSONResponse(dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader))

	return resp, nil
}

func (a *ApiHandler) UnlikeAnimal(ctx context.Context, request api.UnlikeAnimalRequestObject) (api.UnlikeAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "unlike",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UnlikeAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)
	if err := a.persistor.Animal().UnlikeAnimal(ctx, userID, request.ID); err != nil {
		return nil, fmt.Errorf("failed to unlike an animal: %w", err)
	}

	return api.EmptyResponseResponse{}, nil
}

func (a *ApiHandler) ApproveAnimal(ctx context.Context, request api.ApproveAnimalRequestObject) (api.ApproveAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ApproveAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().ApproveAnimal(ctx, request.ID); err != nil {
		return api.ApproveAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_approve", err.Error())}, nil
	}

	return api.ApproveAnimal204Response{}, nil
}

func (a *ApiHandler) RejectAnimal(ctx context.Context, request api.RejectAnimalRequestObject) (api.RejectAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.RejectAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	if err := a.persistor.Animal().RejectAnimal(ctx, request.ID); err != nil {
		return api.RejectAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_reject", err.Error())}, nil
	}

	return api.ApproveOrganization204Response{}, nil
}

func createAnimalRelationTuples(ctx context.Context, c *keto.Client, identityID string, animalID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Animal",
					Object:    shared.AuthzAnimalID(animalID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Animal",
					Object:    shared.AuthzAnimalID(animalID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Animals", "animals", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert animal relation tuples: %w", err)
	}

	return nil
}

func deleteAnimalRelationTuples(ctx context.Context, c *keto.Client, animalID int64) error {
	ownersResp, err := c.Read.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: new("Animal"),
			Object:    new(shared.AuthzAnimalID(animalID)),
			Relation:  new("owners"),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list animal relation tuples: %w", err)
	}

	tuplesToDelete := []*rts.RelationTupleDelta{
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Animal",
				Object:    shared.AuthzAnimalID(animalID),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Animals", "animals", ""),
			},
		},
	}

	for _, tuple := range ownersResp.RelationTuples {
		subject := tuple.GetSubject()
		if subject == nil {
			continue
		}

		tuplesToDelete = append(tuplesToDelete, &rts.RelationTupleDelta{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Animal",
				Object:    shared.AuthzAnimalID(animalID),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(subject.GetId()),
			},
		})
	}

	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuplesToDelete,
	}); err != nil {
		return fmt.Errorf("failed to delete animal relation tuples: %w", err)
	}

	return nil
}
