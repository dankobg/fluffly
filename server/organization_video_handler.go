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

func (a *ApiHandler) ListOrganizationVideos(ctx context.Context, request api.ListOrganizationVideosRequestObject) (api.ListOrganizationVideosResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListOrganizationVideos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_videos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListOrganizationVideosFilters{ListOrganizationVideosParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	orgVideos, err := a.persistor.Organization().ListOrganizationVideos(ctx, request.ID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization videos: %w", err)
	}

	orgVideosData := make([]api.OrganizationVideo, len(orgVideos.Data))
	for i, video := range orgVideos.Data {
		orgVideosData[i] = dto.OrganizationVideoToResponse(video, a.uploader)
	}

	resp := api.ListOrganizationVideos200JSONResponse{
		Data: orgVideosData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, orgVideos.TotalCount),
	}

	return resp, nil
}

func (a *ApiHandler) GetOrganizationVideo(ctx context.Context, request api.GetOrganizationVideoRequestObject) (api.GetOrganizationVideoResponseObject, error) {
	orgVideo, err := a.persistor.Organization().GetOrganizationVideo(ctx, request.ID, request.VideoID)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationVideoNotFound) {
			return api.GetOrganizationVideo404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("organization_video_not_found", "organization video not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an organization video by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetOrganizationVideodefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_video_permission", "permission denied")}, nil
	}

	resp := api.GetOrganizationVideo200JSONResponse(dto.OrganizationVideoToResponse(orgVideo, a.uploader))

	return resp, nil
}

func (a *ApiHandler) CreateOrganizationVideos(ctx context.Context, request api.CreateOrganizationVideosRequestObject) (api.CreateOrganizationVideosResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateOrganizationVideos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_videos_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	videosDocData := form.Value["data"]
	videosFileData := form.File["data"]

	var videosDataBytes []byte

	if len(videosDocData) == 0 && len(videosFileData) == 0 {
		return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide organization media json data document")}, nil
	}

	if len(videosDocData) > 0 {
		videosDataBytes = []byte(videosDocData[0])
	}

	if len(videosFileData) > 0 {
		orgFile, err := videosFileData[0].Open()
		if err != nil {
			return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		videosDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.CreateOrganizationVideosBody
	if err := json.Unmarshal(videosDataBytes, &input); err != nil {
		return api.CreateOrganizationVideosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_videos_internal", "failed to unmarshal non file data")}, nil
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
		videoResults := a.uploadOrganizationFiles(ctx, videoFileSources, 5)
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

		return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	organizationVideoSetters := make([]models.OrganizationVideoSetter, len(videoUploadResults))
	if len(videoUploadResults) > 0 {
		for i, videoRes := range videoUploadResults {
			// @TODO: optimize video size
			organizationVideoSetters[i] = models.OrganizationVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}
	}

	orgVideos, err := a.persistor.Organization().CreateOrganizationVideos(ctx, request.ID, organizationVideoSetters)
	if err != nil {
		msg := "could not create organization videos"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_videos_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.CreateOrganizationVideos400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_videos_save", msg, reason)}, nil
		}

		return api.CreateOrganizationVideosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_videos_save", msg, reason)}, nil
	}

	var data []api.OrganizationVideo
	for _, vi := range orgVideos {
		data = append(data, dto.OrganizationVideoToResponse(vi, a.uploader))
	}

	resp := api.CreateOrganizationVideos200JSONResponse(api.CreateOrganizationVideos200JSONResponse{Data: data})

	// @TODO: use outbox pattern
	// if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
	// 	return api.CreateOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_photos_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizationVideos(ctx context.Context, request api.DeleteOrganizationVideosRequestObject) (api.DeleteOrganizationVideosResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteOrganizationVideos204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizationVideos403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_videos_permission", "permission denied")}, nil
	}

	filters := dbtype.ListOrganizationVideosFilters{ListOrganizationVideosParams: api.ListOrganizationVideosParams{ID: &request.Body.Ids}}

	oldVideosResults, err := a.persistor.Organization().ListOrganizationVideos(ctx, request.ID, filters)
	if err != nil {
		a.Log.Error("failed to get old organization videos", slog.Int64("organization_id", request.ID), slog.Any("video_ids", request.Body.Ids), slog.Any("error", err))
	}

	var oldVideos []string

	for _, rez := range oldVideosResults.Data {
		if rez.ObjectRef != "" {
			oldVideos = append(oldVideos, rez.ObjectRef)
		}
	}

	if err := a.persistor.Organization().DeleteOrganizationVideos(ctx, request.ID, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a organization videos by ids: %w", err)
	}

	resp := api.DeleteOrganizationVideos204Response{}

	go func(oldNames []string) {
		if len(oldNames) > 0 {
			_ = a.deleteUploadedFiles(context.Background(), oldNames, len(oldNames))
		}
	}(oldVideos)

	// @TODO: use outbox pattern
	// for _, id := range request.Body.Ids {
	// 	if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
	// 		a.Log.Error("failed to delete organization videos relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
	// 		return api.DeleteOrganizationVideosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_videos_permissions", "failed to delete permissions")}, nil
	// 	}
	// }

	return resp, nil
}

func (a *ApiHandler) UpdateOrganizationVideo(ctx context.Context, request api.UpdateOrganizationVideoRequestObject) (api.UpdateOrganizationVideoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateOrganizationVideo403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	oldVideo, err := a.persistor.Organization().GetOrganizationVideo(ctx, request.ID, request.VideoID)
	if err != nil {
		a.Log.Error("failed to get old organization video", slog.Int64("organization_id", request.ID), slog.Int64("video_id", request.VideoID), slog.Any("error", err))
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.UpdateOrganizationVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	videoURLs := form.Value["video_url"]
	videoFileHeaders := form.File["video"]

	var videoFileSources []uploadSource

	if len(videoFileHeaders) == 0 && len(videoURLs) == 0 {
		return api.UpdateOrganizationVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "provide video or video_url", err.Error())}, nil
	}

	if len(videoFileHeaders) > 0 {
		videoFileSources = append(videoFileSources, multipartSource{fh: videoFileHeaders[0]})
	}

	if len(videoURLs) > 0 && len(videoFileHeaders) == 0 {
		videoFileSources = append(videoFileSources, urlSource{c: a.httpc, url: videoURLs[0]})
	}

	videoUploadResult := a.uploadOrganizationFiles(ctx, videoFileSources, 1)

	var fileToDelete string
	if videoUploadResult[0].Err != nil {
		fileToDelete = videoUploadResult[0].Name
	}

	if fileToDelete != "" {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), []string{fileToDelete}, 1)
		}()

		return api.UpdateOrganizationVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload video")}, nil
	}

	// @TODO: optimize video size
	orgVideoSetter := models.OrganizationVideoSetter{
		ObjectKind: omit.From(a.uploader.Kind()),
		ObjectRef:  omit.From(videoUploadResult[0].Name),
	}

	orgVideo, err := a.persistor.Organization().UpdateOrganizationVideo(ctx, request.ID, request.VideoID, orgVideoSetter)
	if err != nil {
		msg := "could not update organization video"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateOrganizationVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_video_update", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.UpdateOrganizationVideo400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_video_update", msg, reason)}, nil
		}

		return api.UpdateOrganizationVideodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_video_update", msg, reason)}, nil
	}

	resp := api.UpdateOrganizationVideo200JSONResponse(dto.OrganizationVideoToResponse(orgVideo, a.uploader))

	go func(oldName string) {
		if oldName != "" {
			deleteVideoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deleteVideoResult[0].Err != nil {
				a.Log.Error("failed to delete old organization video", slog.String("name", oldName), slog.Any("error", deleteVideoResult[0].Err))
			}
		}
	}(oldVideo.ObjectRef)

	// // @TODO: use outbox pattern
	// if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
	// 	return api.UpdateOrganizationVideodefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizationVideo(ctx context.Context, request api.DeleteOrganizationVideoRequestObject) (api.DeleteOrganizationVideoResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizationVideo403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_video_permission", "permission denied")}, nil
	}

	oldVideo, err := a.persistor.Organization().GetOrganizationVideo(ctx, request.ID, request.VideoID)
	if err != nil {
		a.Log.Error("failed to get old organization video", slog.Int64("organization_id", request.ID), slog.Int64("video_id", request.VideoID), slog.Any("error", err))
	}

	if _, err := a.persistor.Organization().DeleteOrganizationVideo(ctx, request.ID, request.VideoID); err != nil {
		return nil, fmt.Errorf("failed to delete a organization video: %w", err)
	}

	resp := api.DeleteOrganizationVideo204Response{}

	go func(oldName string) {
		if oldName != "" {
			deleteVideoResult := a.deleteUploadedFiles(context.Background(), []string{oldName}, 1)
			if deleteVideoResult[0].Err != nil {
				a.Log.Error("failed to delete old organization video", slog.String("name", oldName), slog.Any("error", deleteVideoResult[0].Err))
			}
		}
	}(oldVideo.ObjectRef)

	// @TODO: use outbox pattern
	// if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
	// 	a.Log.Error("failed to delete organization video relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
	// 	return api.DeleteOrganizationPhotosdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_video_permissions", "failed to delete permissions")}, nil
	// }

	return resp, nil
}
