package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListAnimalPhotos(ctx context.Context, request api.ListAnimalPhotosRequestObject) (api.ListAnimalPhotosResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAnimalPhotos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_photos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalPhotosFilters{ListAnimalPhotosParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	orgPhotos, err := a.persistor.Animal().ListAnimalPhotos(ctx, request.ID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animal photos: %w", err)
	}

	orgPhotosData := make([]api.AnimalPhoto, len(orgPhotos.Data))
	for i, photo := range orgPhotos.Data {
		orgPhotosData[i] = dto.AnimalPhotoToResponse(photo, a.uploader)
	}

	resp := api.ListAnimalPhotos200JSONResponse{
		Data: orgPhotosData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, orgPhotos.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetAnimalPhoto(ctx context.Context, request api.GetAnimalPhotoRequestObject) (api.GetAnimalPhotoResponseObject, error) {
	orgPhoto, err := a.persistor.Animal().GetAnimalPhoto(ctx, request.ID, request.PhotoID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalPhotoNotFound) {
			return api.GetAnimalPhoto404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_photo_not_found", "animal photo not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an animal photo by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnimalPhotodefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "animal_photo_permission", "permission denied")}, nil
	}

	resp := api.GetAnimalPhoto200JSONResponse(dto.AnimalPhotoToResponse(orgPhoto, a.uploader))

	return resp, nil
}

func (a *ApiHandler) CreateAnimalPhotos(ctx context.Context, request api.CreateAnimalPhotosRequestObject) (api.CreateAnimalPhotosResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateAnimalPhotos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_photos_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	photosDocData := form.Value["data"]
	photosFileData := form.File["data"]

	var photosDataBytes []byte

	if len(photosDocData) == 0 && len(photosFileData) == 0 {
		return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide animal media json data document")}, nil
	}

	if len(photosDocData) > 0 {
		photosDataBytes = []byte(photosDocData[0])
	}

	if len(photosFileData) > 0 {
		orgFile, err := photosFileData[0].Open()
		if err != nil {
			return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		photosDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.CreateAnimalPhotosBody
	if err := json.Unmarshal(photosDataBytes, &input); err != nil {
		return api.CreateAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photos_internal", "failed to unmarshal non file data")}, nil
	}

	photoFileHeaders := form.File["photos"]

	var photoFileSources []uploadSource
	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	if input.Photos != nil {
		for _, photo := range *input.Photos {
			photoFileSources = append(photoFileSources, urlSource{c: a.httpc, url: photo.URL})
		}
	}

	var photoUploadResults []uploadResult

	if len(photoFileSources) > 0 {
		photoResults := a.uploadAnimalFiles(ctx, photoFileSources, 5)
		photoUploadResults = append(photoUploadResults, photoResults...)
	}

	var filesToDelete []string

	for _, res := range photoUploadResults {
		if res.Err != nil {
			filesToDelete = append(filesToDelete, res.Name)
		}
	}

	if len(filesToDelete) > 0 {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	animalPhotoSetters := make([]models.AnimalPhotoSetter, len(photoUploadResults))
	if len(photoUploadResults) > 0 {
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
	}

	orgPhotos, err := a.persistor.Animal().CreateAnimalPhotos(ctx, request.ID, animalPhotoSetters)
	if err != nil {
		msg := "could not create animal photos"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_photos_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.CreateAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_photos_save", msg, reason)}, nil
		}

		return api.CreateAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photos_save", msg, reason)}, nil
	}

	var data []api.AnimalPhoto
	for _, ph := range orgPhotos {
		data = append(data, dto.AnimalPhotoToResponse(ph, a.uploader))
	}

	resp := api.CreateAnimalPhotos200JSONResponse(api.CreateAnimalPhotos200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photos_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalPhotos(ctx context.Context, request api.DeleteAnimalPhotosRequestObject) (api.DeleteAnimalPhotosResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimalPhotos204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalPhotos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_photos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalPhotosFilters{ListAnimalPhotosParams: api.ListAnimalPhotosParams{ID: &request.Body.Ids}}

	oldPhotosResults, err := a.persistor.Animal().ListAnimalPhotos(ctx, request.ID, filters)
	if err != nil {
		a.Log.Error("failed to get old animal photos", slog.Int64("animal_id", request.ID), slog.Any("photo_ids", request.Body.Ids), slog.Any("error", err))
		return api.DeleteAnimalPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to list old photos", err.Error())}, nil
	}

	var oldPhotos []string

	for _, rez := range oldPhotosResults.Data {
		if !rez.ObjectRefFull.IsNull() {
			oldPhotos = append(oldPhotos, rez.ObjectRefFull.MustGet())
		}
	}

	if err := a.persistor.Animal().DeleteAnimalPhotos(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a animal photos by ids: %w", err)
	}

	resp := api.DeleteAnimalPhotos204Response{}

	go func(oldNames []string) {
		if len(oldNames) > 0 {
			_ = a.deleteUploadedFiles(context.Background(), oldNames, len(oldNames))
		}
	}(oldPhotos)

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete animal photos relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 		return api.DeleteAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photos_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateAnimalPhoto(ctx context.Context, request api.UpdateAnimalPhotoRequestObject) (api.UpdateAnimalPhotoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimalPhoto403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	oldPhoto, err := a.persistor.Animal().GetAnimalPhoto(ctx, request.ID, request.PhotoID)
	if err != nil {
		a.Log.Error("failed to get old animal photo", slog.Int64("animal_id", request.ID), slog.Int64("photo_id", request.PhotoID), slog.Any("error", err))
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.UpdateAnimalPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	photoURLs := form.Value["photo_url"]
	photoFileHeaders := form.File["photo"]

	var photoFileSources []uploadSource

	if len(photoFileHeaders) == 0 && len(photoURLs) == 0 {
		return api.UpdateAnimalPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "provide photo or photo_url", err.Error())}, nil
	}

	if len(photoFileHeaders) > 0 {
		photoFileSources = append(photoFileSources, multipartSource{fh: photoFileHeaders[0]})
	}

	if len(photoURLs) > 0 && len(photoFileHeaders) == 0 {
		photoFileSources = append(photoFileSources, urlSource{c: a.httpc, url: photoURLs[0]})
	}

	photoUploadResult := a.uploadAnimalFiles(ctx, photoFileSources, 1)

	var fileToDelete string
	if photoUploadResult[0].Err != nil {
		fileToDelete = photoUploadResult[0].Name
	}

	if fileToDelete != "" {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), []string{fileToDelete}, 1)
		}()

		return api.UpdateAnimalPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload photo")}, nil
	}

	// @TODO: optimize image size
	orgPhotoSetter := models.AnimalPhotoSetter{
		ObjectKind:      omit.From(a.uploader.Kind()),
		ObjectRefSmall:  omitnull.From(photoUploadResult[0].Name),
		ObjectRefMedium: omitnull.From(photoUploadResult[0].Name),
		ObjectRefLarge:  omitnull.From(photoUploadResult[0].Name),
		ObjectRefFull:   omitnull.From(photoUploadResult[0].Name),
	}

	orgPhoto, err := a.persistor.Animal().UpdateAnimalPhoto(ctx, request.ID, request.PhotoID, orgPhotoSetter)
	if err != nil {
		msg := "could not update animal photo"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimalPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_photo_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.UpdateAnimalPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_photo_update", msg, reason)}, nil
		}

		return api.UpdateAnimalPhotodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photo_update", msg, reason)}, nil
	}

	resp := api.UpdateAnimalPhoto200JSONResponse(dto.AnimalPhotoToResponse(orgPhoto, a.uploader))

	// // @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.UpdateAnimalPhotodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_permissions", "failed to create permissions")}, nil
	// }

	go func(oldName string) {
		if oldName != "" {
			deletePhotoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deletePhotoResult[0].Err != nil {
				a.Log.Error("failed to delete old animal photo", slog.String("name", oldName), slog.Any("error", deletePhotoResult[0].Err))
			}
		}
	}(oldPhoto.ObjectRefFull.MustGet())

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalPhoto(ctx context.Context, request api.DeleteAnimalPhotoRequestObject) (api.DeleteAnimalPhotoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalPhoto403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_photo_permission", "permission denied")}, nil
	}

	oldPhoto, err := a.persistor.Animal().GetAnimalPhoto(ctx, request.ID, request.PhotoID)
	if err != nil {
		a.Log.Error("failed to get old animal photo", slog.Int64("animal_id", request.ID), slog.Int64("photo_id", request.PhotoID), slog.Any("error", err))
	}

	if _, err := a.persistor.Animal().DeleteAnimalPhoto(ctx, request.ID, request.PhotoID); err != nil {
		return nil, fmt.Errorf("failed to delete a animal photo: %w", err)
	}

	resp := api.DeleteAnimalPhoto204Response{}

	go func(oldName string) {
		if oldName != "" {
			deletePhotoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deletePhotoResult[0].Err != nil {
				a.Log.Error("failed to delete old animal photo", slog.String("name", oldName), slog.Any("error", deletePhotoResult[0].Err))
			}
		}
	}(oldPhoto.ObjectRefFull.MustGet())

	// @TODO: use outbox pattern
	// if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete animal photos relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 	return api.DeleteAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photos_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
