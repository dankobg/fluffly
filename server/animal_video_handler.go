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
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListAnimalVideos(ctx context.Context, request api.ListAnimalVideosRequestObject) (api.ListAnimalVideosResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListAnimalVideos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_videos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalVideosFilters{ListAnimalVideosParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	orgVideos, err := a.persistor.Animal().ListAnimalVideos(ctx, request.ID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list animal videos: %w", err)
	}

	orgVideosData := make([]api.AnimalVideo, len(orgVideos.Data))
	for i, video := range orgVideos.Data {
		orgVideosData[i] = dto.AnimalVideoToResponse(video, a.uploader)
	}

	resp := api.ListAnimalVideos200JSONResponse{
		Data: orgVideosData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, orgVideos.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetAnimalVideo(ctx context.Context, request api.GetAnimalVideoRequestObject) (api.GetAnimalVideoResponseObject, error) {
	orgVideo, err := a.persistor.Animal().GetAnimalVideo(ctx, request.ID, request.VideoID)
	if err != nil {
		if errors.Is(err, postgres.ErrAnimalVideoNotFound) {
			return api.GetAnimalVideo404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("animal_video_not_found", "animal video not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an animal video by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetAnimalVideodefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "animal_video_permission", "permission denied")}, nil
	}

	resp := api.GetAnimalVideo200JSONResponse(dto.AnimalVideoToResponse(orgVideo, a.uploader))

	return resp, nil
}

func (a *ApiHandler) CreateAnimalVideos(ctx context.Context, request api.CreateAnimalVideosRequestObject) (api.CreateAnimalVideosResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateAnimalVideos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_videos_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	videosDocData := form.Value["data"]
	videosFileData := form.File["data"]

	var videosDataBytes []byte

	if len(videosDocData) == 0 && len(videosFileData) == 0 {
		return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide animal media json data document")}, nil
	}

	if len(videosDocData) > 0 {
		videosDataBytes = []byte(videosDocData[0])
	}

	if len(videosFileData) > 0 {
		orgFile, err := videosFileData[0].Open()
		if err != nil {
			return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		videosDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.CreateAnimalVideosBody
	if err := json.Unmarshal(videosDataBytes, &input); err != nil {
		return api.CreateAnimalVideosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_videos_internal", "failed to unmarshal non file data")}, nil
	}

	videoFileHeaders := form.File["videos"]

	var videoFileSources []uploadSource
	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
	}

	if input.Videos != nil {
		for _, photo := range *input.Videos {
			videoFileSources = append(videoFileSources, urlSource{c: a.httpc, url: photo.URL})
		}
	}

	var videoUploadResults []uploadResult

	if len(videoFileSources) > 0 {
		videoResults := a.uploadAnimalFiles(ctx, videoFileSources, 5)
		videoUploadResults = append(videoUploadResults, videoResults...)
	}

	var filesToDelete []string

	for _, res := range videoUploadResults {
		if res.Err != nil {
			filesToDelete = append(filesToDelete, res.Name)
		}
	}

	if len(filesToDelete) > 0 {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	animalVideoSetters := make([]models.AnimalVideoSetter, len(videoUploadResults))
	if len(videoUploadResults) > 0 {
		for i, videoRes := range videoUploadResults {
			// @TODO: optimize video size
			animalVideoSetters[i] = models.AnimalVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}
	}

	orgVideos, err := a.persistor.Animal().CreateAnimalVideos(ctx, request.ID, animalVideoSetters)
	if err != nil {
		msg := "could not create animal videos"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_videos_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.CreateAnimalVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_videos_save", msg, reason)}, nil
		}

		return api.CreateAnimalVideosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_videos_save", msg, reason)}, nil
	}

	var data []api.AnimalVideo
	for _, vi := range orgVideos {
		data = append(data, dto.AnimalVideoToResponse(vi, a.uploader))
	}

	resp := api.CreateAnimalVideos200JSONResponse(api.CreateAnimalVideos200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_photos_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalVideos(ctx context.Context, request api.DeleteAnimalVideosRequestObject) (api.DeleteAnimalVideosResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteAnimalVideos204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animals",
			Object:    "animals", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalVideos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_videos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListAnimalVideosFilters{ListAnimalVideosParams: api.ListAnimalVideosParams{ID: &request.Body.Ids}}

	oldVideosResults, err := a.persistor.Animal().ListAnimalVideos(ctx, request.ID, filters)
	if err != nil {
		a.Log.Error("failed to get old animal videos", slog.Int64("animal_id", request.ID), slog.Any("video_ids", request.Body.Ids), slog.Any("error", err))
	}

	var oldVideos []string

	for _, rez := range oldVideosResults.Data {
		if rez.ObjectRef != "" {
			oldVideos = append(oldVideos, rez.ObjectRef)
		}
	}

	if err := a.persistor.Animal().DeleteAnimalVideos(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a animal videos by ids: %w", err)
	}

	resp := api.DeleteAnimalVideos204Response{}

	go func(oldNames []string) {
		if len(oldNames) > 0 {
			_ = a.deleteUploadedFiles(context.Background(), oldNames, len(oldNames))
		}
	}(oldVideos)

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete animal videos relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 		return api.DeleteAnimalVideosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_videos_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateAnimalVideo(ctx context.Context, request api.UpdateAnimalVideoRequestObject) (api.UpdateAnimalVideoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateAnimalVideo403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_permission", "permission denied")}, nil
	}

	oldVideo, err := a.persistor.Animal().GetAnimalVideo(ctx, request.ID, request.VideoID)
	if err != nil {
		a.Log.Error("failed to get old animal video", slog.Int64("animal_id", request.ID), slog.Int64("video_id", request.VideoID), slog.Any("error", err))
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.UpdateAnimalVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	videoURLs := form.Value["video_url"]
	videoFileHeaders := form.File["video"]

	var videoFileSources []uploadSource

	if len(videoFileHeaders) == 0 && len(videoURLs) == 0 {
		return api.UpdateAnimalVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "provide video or video_url", err.Error())}, nil
	}

	if len(videoFileHeaders) > 0 {
		videoFileSources = append(videoFileSources, multipartSource{fh: videoFileHeaders[0]})
	}

	if len(videoURLs) > 0 && len(videoFileHeaders) == 0 {
		videoFileSources = append(videoFileSources, urlSource{c: a.httpc, url: videoURLs[0]})
	}

	videoUploadResult := a.uploadAnimalFiles(ctx, videoFileSources, 1)

	var fileToDelete string
	if videoUploadResult[0].Err != nil {
		fileToDelete = videoUploadResult[0].Name
	}

	if fileToDelete != "" {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), []string{fileToDelete}, 1)
		}()

		return api.UpdateAnimalVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload video")}, nil
	}

	// @TODO: optimize video size
	orgVideoSetter := models.AnimalVideoSetter{
		ObjectKind: omit.From(a.uploader.Kind()),
		ObjectRef:  omit.From(videoUploadResult[0].Name),
	}

	orgVideo, err := a.persistor.Animal().UpdateAnimalVideo(ctx, request.ID, request.VideoID, orgVideoSetter)
	if err != nil {
		msg := "could not update animal video"

		var (
			reason string
			e1     postgres.ErrAnimalUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateAnimalVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_video_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "animal integrity error"
			return api.UpdateAnimalVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "animal_video_update", msg, reason)}, nil
		}

		return api.UpdateAnimalVideodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_video_update", msg, reason)}, nil
	}

	resp := api.UpdateAnimalVideo200JSONResponse(dto.AnimalVideoToResponse(orgVideo, a.uploader))

	go func(oldName string) {
		if oldName != "" {
			deleteVideoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deleteVideoResult[0].Err != nil {
				a.Log.Error("failed to delete old animal video", slog.String("name", oldName), slog.Any("error", deleteVideoResult[0].Err))
			}
		}
	}(oldVideo.ObjectRef)

	// // @TODO: use outbox pattern
	// if err := createAnimalRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert animal relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("animal_id", resp.ID), slog.Any("error", err))
	// 	return api.UpdateAnimalVideodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteAnimalVideo(ctx context.Context, request api.DeleteAnimalVideoRequestObject) (api.DeleteAnimalVideoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Animal",
			Object:    shared.AuthzAnimalID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteAnimalVideo403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("animal_video_permission", "permission denied")}, nil
	}

	oldVideo, err := a.persistor.Animal().GetAnimalVideo(ctx, request.ID, request.VideoID)
	if err != nil {
		a.Log.Error("failed to get old animal video", slog.Int64("animal_id", request.ID), slog.Int64("video_id", request.VideoID), slog.Any("error", err))
	}

	if _, err := a.persistor.Animal().DeleteAnimalVideo(ctx, request.ID, request.VideoID); err != nil {
		return nil, fmt.Errorf("failed to delete a animal video: %w", err)
	}

	resp := api.DeleteAnimalVideo204Response{}

	go func(oldName string) {
		if oldName != "" {
			deleteVideoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deleteVideoResult[0].Err != nil {
				a.Log.Error("failed to delete old animal video", slog.String("name", oldName), slog.Any("error", deleteVideoResult[0].Err))
			}
		}
	}(oldVideo.ObjectRef)

	// @TODO: use outbox pattern
	// if err := deleteAnimalRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete animal video relation-tuple", slog.Int64("animal_id", id), slog.Any("error", err))
	// 	return api.DeleteAnimalPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "animal_video_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
