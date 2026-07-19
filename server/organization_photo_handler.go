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

func (a *ApiHandler) ListOrganizationPhotos(ctx context.Context, request api.ListOrganizationPhotosRequestObject) (api.ListOrganizationPhotosResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListOrganizationPhotos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_photos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListOrganizationPhotosFilters{ListOrganizationPhotosParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	orgPhotos, err := a.persistor.Organization().ListOrganizationPhotos(ctx, request.ID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization photos: %w", err)
	}

	orgPhotosData := make([]api.OrganizationPhoto, len(orgPhotos.Data))
	for i, photo := range orgPhotos.Data {
		orgPhotosData[i] = dto.OrganizationPhotoToResponse(photo, a.uploader)
	}

	resp := api.ListOrganizationPhotos200JSONResponse{
		Data: orgPhotosData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, orgPhotos.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetOrganizationPhoto(ctx context.Context, request api.GetOrganizationPhotoRequestObject) (api.GetOrganizationPhotoResponseObject, error) {
	orgPhoto, err := a.persistor.Organization().GetOrganizationPhoto(ctx, request.ID, request.PhotoID)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationPhotoNotFound) {
			return api.GetOrganizationPhoto404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("organization_photo_not_found", "organization photo not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an organization photo by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetOrganizationPhotodefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_photo_permission", "permission denied")}, nil
	}

	resp := api.GetOrganizationPhoto200JSONResponse(dto.OrganizationPhotoToResponse(orgPhoto, a.uploader))

	return resp, nil
}

func (a *ApiHandler) CreateOrganizationPhotos(ctx context.Context, request api.CreateOrganizationPhotosRequestObject) (api.CreateOrganizationPhotosResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateOrganizationPhotos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_photos_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	photosDocData := form.Value["data"]
	photosFileData := form.File["data"]

	var photosDataBytes []byte

	if len(photosDocData) == 0 && len(photosFileData) == 0 {
		return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide organization media json data document")}, nil
	}

	if len(photosDocData) > 0 {
		photosDataBytes = []byte(photosDocData[0])
	}

	if len(photosFileData) > 0 {
		orgFile, err := photosFileData[0].Open()
		if err != nil {
			return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		photosDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.CreateOrganizationPhotosBody
	if err := json.Unmarshal(photosDataBytes, &input); err != nil {
		return api.CreateOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photos_internal", "failed to unmarshal non file data")}, nil
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
		photoResults := a.uploadOrganizationFiles(ctx, photoFileSources, 5)
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

		return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	organizationPhotoSetters := make([]models.OrganizationPhotoSetter, len(photoUploadResults))
	if len(photoUploadResults) > 0 {
		for i, photoRes := range photoUploadResults {
			// @TODO: optimize image size
			organizationPhotoSetters[i] = models.OrganizationPhotoSetter{
				ObjectKind:      omit.From(a.uploader.Kind()),
				ObjectRefSmall:  omitnull.From(photoRes.Name),
				ObjectRefMedium: omitnull.From(photoRes.Name),
				ObjectRefLarge:  omitnull.From(photoRes.Name),
				ObjectRefFull:   omitnull.From(photoRes.Name),
			}
		}
	}

	orgPhotos, err := a.persistor.Organization().CreateOrganizationPhotos(ctx, request.ID, organizationPhotoSetters)
	if err != nil {
		msg := "could not create organization photos"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_photos_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.CreateOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_photos_save", msg, reason)}, nil
		}

		return api.CreateOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photos_save", msg, reason)}, nil
	}

	var data []api.OrganizationPhoto
	for _, ph := range orgPhotos {
		data = append(data, dto.OrganizationPhotoToResponse(ph, a.uploader))
	}

	resp := api.CreateOrganizationPhotos200JSONResponse(api.CreateOrganizationPhotos200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photos_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizationPhotos(ctx context.Context, request api.DeleteOrganizationPhotosRequestObject) (api.DeleteOrganizationPhotosResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteOrganizationPhotos204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizationPhotos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_photos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListOrganizationPhotosFilters{ListOrganizationPhotosParams: api.ListOrganizationPhotosParams{ID: &request.Body.Ids}}

	oldPhotosResults, err := a.persistor.Organization().ListOrganizationPhotos(ctx, request.ID, filters)
	if err != nil {
		a.Log.Error("failed to get old organization photos", slog.Int64("organization_id", request.ID), slog.Any("photo_ids", request.Body.Ids), slog.Any("error", err))
		return api.DeleteOrganizationPhotos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to list old photos", err.Error())}, nil
	}

	var oldPhotos []string

	for _, rez := range oldPhotosResults.Data {
		if !rez.ObjectRefFull.IsNull() {
			oldPhotos = append(oldPhotos, rez.ObjectRefFull.MustGet())
		}
	}

	if err := a.persistor.Organization().DeleteOrganizationPhotos(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a organization photos by ids: %w", err)
	}

	resp := api.DeleteOrganizationPhotos204Response{}

	go func(oldNames []string) {
		if len(oldNames) > 0 {
			_ = a.deleteUploadedFiles(context.Background(), oldNames, len(oldNames))
		}
	}(oldPhotos)

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete organization photos relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
	// 		return api.DeleteOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photos_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateOrganizationPhoto(ctx context.Context, request api.UpdateOrganizationPhotoRequestObject) (api.UpdateOrganizationPhotoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateOrganizationPhoto403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	oldPhoto, err := a.persistor.Organization().GetOrganizationPhoto(ctx, request.ID, request.PhotoID)
	if err != nil {
		a.Log.Error("failed to get old organization photo", slog.Int64("organization_id", request.ID), slog.Int64("photo_id", request.PhotoID), slog.Any("error", err))
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.UpdateOrganizationPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	photoURLs := form.Value["photo_url"]
	photoFileHeaders := form.File["photo"]

	var photoFileSources []uploadSource

	if len(photoFileHeaders) == 0 && len(photoURLs) == 0 {
		return api.UpdateOrganizationPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "provide photo or photo_url", err.Error())}, nil
	}

	if len(photoFileHeaders) > 0 {
		photoFileSources = append(photoFileSources, multipartSource{fh: photoFileHeaders[0]})
	}

	if len(photoURLs) > 0 && len(photoFileHeaders) == 0 {
		photoFileSources = append(photoFileSources, urlSource{c: a.httpc, url: photoURLs[0]})
	}

	photoUploadResult := a.uploadOrganizationFiles(ctx, photoFileSources, 1)

	var fileToDelete string
	if photoUploadResult[0].Err != nil {
		fileToDelete = photoUploadResult[0].Name
	}

	if fileToDelete != "" {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), []string{fileToDelete}, 1)
		}()

		return api.UpdateOrganizationPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload photo")}, nil
	}

	// @TODO: optimize image size
	orgPhotoSetter := models.OrganizationPhotoSetter{
		ObjectKind:      omit.From(a.uploader.Kind()),
		ObjectRefSmall:  omitnull.From(photoUploadResult[0].Name),
		ObjectRefMedium: omitnull.From(photoUploadResult[0].Name),
		ObjectRefLarge:  omitnull.From(photoUploadResult[0].Name),
		ObjectRefFull:   omitnull.From(photoUploadResult[0].Name),
	}

	orgPhoto, err := a.persistor.Organization().UpdateOrganizationPhoto(ctx, request.ID, request.PhotoID, orgPhotoSetter)
	if err != nil {
		msg := "could not update organization photo"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateOrganizationPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_photo_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.UpdateOrganizationPhoto400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_photo_update", msg, reason)}, nil
		}

		return api.UpdateOrganizationPhotodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photo_update", msg, reason)}, nil
	}

	resp := api.UpdateOrganizationPhoto200JSONResponse(dto.OrganizationPhotoToResponse(orgPhoto, a.uploader))

	// // @TODO: use outbox pattern
	// if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
	// 	return api.UpdateOrganizationPhotodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to create permissions")}, nil
	// }

	go func(oldName string) {
		if oldName != "" {
			deletePhotoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deletePhotoResult[0].Err != nil {
				a.Log.Error("failed to delete old organization photo", slog.String("name", oldName), slog.Any("error", deletePhotoResult[0].Err))
			}
		}
	}(oldPhoto.ObjectRefFull.MustGet())

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizationPhoto(ctx context.Context, request api.DeleteOrganizationPhotoRequestObject) (api.DeleteOrganizationPhotoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizationPhoto403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_photo_permission", "permission denied")}, nil
	}

	oldPhoto, err := a.persistor.Organization().GetOrganizationPhoto(ctx, request.ID, request.PhotoID)
	if err != nil {
		a.Log.Error("failed to get old organization photo", slog.Int64("organization_id", request.ID), slog.Int64("photo_id", request.PhotoID), slog.Any("error", err))
	}

	if _, err := a.persistor.Organization().DeleteOrganizationPhoto(ctx, request.ID, request.PhotoID); err != nil {
		return nil, fmt.Errorf("failed to delete a organization photo: %w", err)
	}

	resp := api.DeleteOrganizationPhoto204Response{}

	go func(oldName string) {
		if oldName != "" {
			deletePhotoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deletePhotoResult[0].Err != nil {
				a.Log.Error("failed to delete old organization photo", slog.String("name", oldName), slog.Any("error", deletePhotoResult[0].Err))
			}
		}
	}(oldPhoto.ObjectRefFull.MustGet())

	// @TODO: use outbox pattern
	// if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete organization photos relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
	// 	return api.DeleteOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photos_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
