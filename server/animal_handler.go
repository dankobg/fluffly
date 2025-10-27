package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/ptr"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

const (
	createAnimalFileMaxMemory = 50 << 20
)

func (a *ApiHandler) CreateAnimal(ctx context.Context, request api.CreateAnimalRequestObject) (api.CreateAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.CreateAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "invalid permission")}, nil
	}

	form, err := request.Body.ReadForm(createAnimalFileMaxMemory)
	if err != nil {
		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}
	defer form.RemoveAll()
	animalData := form.Value["data"][0]
	var input api.CreateAnimalBody
	if err := json.Unmarshal([]byte(animalData), &input); err != nil {
		return api.CreateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_internal", "failed to unmarshal non file data")}, nil
	}

	mainImageFileHeaders := form.File["image"]
	if len(mainImageFileHeaders) == 0 && (input.ImageURL.IsNull() || (!input.ImageURL.IsNull() && input.ImageURL.MustGet() == "")) {
		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_image_missing", "Provide image file or image_url for main animal photo", err.Error())}, nil
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
		return api.CreateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
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
		var reason string
		var e1 postgres.ErrAnimalUniqueViolation
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
	resp := api.CreateAnimal201JSONResponse(dto.AnimalToResponse(animal))

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
			Object:    AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.UpdateAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "invalid permission")}, nil
	}

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
		var reason string
		var e1 postgres.ErrAnimalUniqueViolation
		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_edit", msg, reason)}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.UpdateAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_edit", msg, reason)}, nil
		}
		return api.UpdateAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_edit", msg, reason)}, nil
	}
	resp := api.UpdateAnimal201JSONResponse(dto.AnimalToResponse(animal))
	return resp, nil
}

func (a *ApiHandler) DeleteAnimal(ctx context.Context, request api.DeleteAnimalRequestObject) (api.DeleteAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DeleteAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "invalid permission")}, nil
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

func (a *ApiHandler) ListAnimals(ctx context.Context, request api.ListAnimalsRequestObject) (api.ListAnimalsResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.Allowed {
		fmt.Println("err", err, checkResp.Allowed)
		return api.ListAnimalsdefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("animal_permission", "invalid permission")}, nil
	}

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
			return api.GetAnimal404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_not_found", "animal not found")}, nil
		}
		return nil, fmt.Errorf("failed to get an animal by id: %w", err)
	}
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    AuthzAnimalID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetAnimaldefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "animal_permission", "invalid permission")}, nil
	}
	resp := api.GetAnimal200JSONResponse(dto.AnimalWithJoinDataToResponse(animalWithJoinData, a.uploader))
	return resp, nil
}

func (a *ApiHandler) LikeAnimal(ctx context.Context, request api.LikeAnimalRequestObject) (api.LikeAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.LikeAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "invalid permission")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)
	if err := a.persistor.Animal().LikeAnimal(ctx, userID, request.ID); err != nil {
		msg := "failed to like an animal"
		var reason string
		var e1 postgres.ErrAnimalUniqueViolation
		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.LikeAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_like", msg, reason)}, nil
		}
		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.LikeAnimal400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_like", msg, reason)}, nil
		}
		return api.LikeAnimaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_like", msg, reason)}, nil
	}
	return api.EmptyResponseResponse{}, nil
}

func (a *ApiHandler) UnlikeAnimal(ctx context.Context, request api.UnlikeAnimalRequestObject) (api.UnlikeAnimalResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.UnlikeAnimal403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "invalid permission")}, nil
	}

	userID := uuid.MustParse(sess.Identity.Id)
	if err := a.persistor.Animal().UnlikeAnimal(ctx, userID, request.ID); err != nil {
		return nil, fmt.Errorf("failed to unlike an animal: %w", err)
	}
	return api.EmptyResponseResponse{}, nil
}

func createAnimalRelationTuples(ctx context.Context, c *keto.Client, identityID string, animalID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Animal",
					Object:    AuthzAnimalID(animalID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Animal",
					Object:    AuthzAnimalID(animalID),
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
			Namespace: ptr.Of("Animal"),
			Object:    ptr.Of(AuthzAnimalID(animalID)),
			Relation:  ptr.Of("owners"),
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
				Object:    AuthzAnimalID(animalID),
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
				Object:    AuthzAnimalID(animalID),
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
