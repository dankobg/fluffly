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
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/convert"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/geocoding"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/shared"
	"github.com/google/uuid"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

const (
	maxFileMemory = 100 << 20
)

func (a *ApiHandler) ApplyForOrganization(ctx context.Context, request api.ApplyForOrganizationRequestObject) (api.ApplyForOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "apply",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ApplyForOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	geocodingResult, err := a.geocoder.ForwardGeocodeStructured(ctx, geocoding.ForwardGeocodeStructuredInput{})
	if err != nil {
		return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to forward geocode structured", err.Error())}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	orgDocData := form.Value["data"]
	orgFileData := form.File["data"]

	var orgDataBytes []byte

	if len(orgDocData) == 0 && len(orgFileData) == 0 {
		return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide organization json data document")}, nil
	}

	if len(orgDocData) > 0 {
		orgDataBytes = []byte(orgDocData[0])
	}

	if len(orgFileData) > 0 {
		orgFile, err := orgFileData[0].Open()
		if err != nil {
			return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		orgDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.ApplyForOrganizationBody
	if err := json.Unmarshal(orgDataBytes, &input); err != nil {
		return api.ApplyForOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_internal", "failed to unmarshal non file data")}, nil
	}

	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]

	var (
		photoFileSources []uploadSource
		videoFileSources []uploadSource
	)

	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
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
		photoUploadResults []uploadResult
		videoUploadResults []uploadResult
	)

	if len(photoFileSources) > 0 {
		photoResults := a.uploadOrganizationFiles(ctx, photoFileSources, 5)
		photoUploadResults = append(photoUploadResults, photoResults...)
	}

	if len(videoFileSources) > 0 {
		videoResults := a.uploadOrganizationFiles(ctx, videoFileSources, 5)
		videoUploadResults = append(videoUploadResults, videoResults...)
	}

	var filesToDelete []string

	for _, res := range append(photoUploadResults, videoUploadResults...) {
		if res.Err != nil {
			filesToDelete = append(filesToDelete, res.Name)
		}
	}

	if len(filesToDelete) > 0 {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	var organizationApplyForSetter dbtype.OrganizationApplyForSetter

	organizationApplyForSetter.Organization = models.OrganizationSetter{
		Name:             omit.From(input.Name),
		Website:          omitnull.FromPtr(input.Website),
		MissionStatement: omitnull.FromPtr(input.MissionStatement),
		AdoptionPolicy:   omitnull.FromPtr(input.AdoptionPolicy),
		AdoptionURL:      omitnull.FromPtr(input.AdoptionURL),
		Distance:         omitnull.FromPtr(input.Distance),
	}

	organizationApplyForSetter.Contact = models.OrganizationContactSetter{
		Phone: omit.From(input.Contact.Phone),
		Email: omit.From(string(input.Contact.Email)),
	}

	point := geom.NewPoint(geom.XY).SetSRID(dbcustom.SRID).MustSetCoords(geom.Coord{geocodingResult.Lon, geocodingResult.Lat})
	coords := &ewkb.Point{Point: point}

	organizationApplyForSetter.Address = models.AddressSetter{
		CountryID:     omit.From(input.Contact.Address.CountryID),
		StreetAddress: omit.From(input.Contact.Address.StreetAddress),
		City:          omit.From(input.Contact.Address.City),
		Coords:        omitnull.From(coords),
		Note:          omitnull.FromPtr(input.Contact.Address.Note),
		Region:        omitnull.FromPtr(input.Contact.Address.Region),
		PostalCode:    omitnull.FromPtr(input.Contact.Address.PostalCode),
		UnitNumber:    omitnull.FromPtr(input.Contact.Address.UnitNumber),
		StreetNumber:  omitnull.FromPtr(input.Contact.Address.StreetNumber),
	}

	if input.WorkHour != nil {
		if (input.WorkHour.Monday != nil) ||
			(input.WorkHour.Tuesday != nil) ||
			(input.WorkHour.Wednesday != nil) ||
			(input.WorkHour.Thursday != nil) ||
			(input.WorkHour.Friday != nil) ||
			(input.WorkHour.Saturday != nil) ||
			(input.WorkHour.Sunday != nil) {
			organizationApplyForSetter.WorkHour = omitnull.From(models.OrganizationWorkHourSetter{
				Monday:    omitnull.FromPtr(input.WorkHour.Monday),
				Tuesday:   omitnull.FromPtr(input.WorkHour.Tuesday),
				Wednesday: omitnull.FromPtr(input.WorkHour.Wednesday),
				Thursday:  omitnull.FromPtr(input.WorkHour.Thursday),
				Friday:    omitnull.FromPtr(input.WorkHour.Friday),
				Saturday:  omitnull.FromPtr(input.WorkHour.Saturday),
				Sunday:    omitnull.FromPtr(input.WorkHour.Sunday),
			})
		}
	}

	if len(photoUploadResults) > 0 {
		organizationPhotoSetters := make([]models.OrganizationPhotoSetter, len(photoUploadResults))
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

		organizationApplyForSetter.Photos = omitnull.From(organizationPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		organizationVideoSetters := make([]models.OrganizationVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			organizationVideoSetters[i] = models.OrganizationVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}

		organizationApplyForSetter.Videos = omitnull.From(organizationVideoSetters)
	}

	if input.Socials != nil {
		organizationSocialsSetters := make([]models.OrganizationSocialSetter, 0)
		for _, social := range *input.Socials {
			organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
				Platform: omit.From(social.Platform),
				URL:      omit.From(social.URL),
			})
		}

		organizationApplyForSetter.Socials = omitnull.From(organizationSocialsSetters)
	}

	applicantID := uuid.MustParse(sess.Identity.Id)

	organization, err := a.persistor.Organization().ApplyForOrganization(ctx, applicantID, organizationApplyForSetter)
	if err != nil {
		msg := "could not apply for organization"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.ApplyForOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_save", msg, reason)}, nil
		}

		return api.ApplyForOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_save", msg, reason)}, nil
	}

	resp := api.ApplyForOrganization201JSONResponse(dto.OrganizationToResponse(organization))

	// @TODO: use outbox pattern
	if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
		a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
		return api.ApplyForOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) ApproveOrganization(ctx context.Context, request api.ApproveOrganizationRequestObject) (api.ApproveOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ApproveOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	if err := a.persistor.Organization().ApproveOrganization(ctx, request.ID); err != nil {
		return api.ApproveOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_approve", err.Error())}, nil
	}

	return api.ApproveOrganization204Response{}, nil
}

func (a *ApiHandler) RejectOrganization(ctx context.Context, request api.RejectOrganizationRequestObject) (api.RejectOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.RejectOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	if err := a.persistor.Organization().RejectOrganization(ctx, request.ID); err != nil {
		return api.RejectOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_Reject", err.Error())}, nil
	}

	return api.ApproveOrganization204Response{}, nil
}

func (a *ApiHandler) CreateOrganization(ctx context.Context, request api.CreateOrganizationRequestObject) (api.CreateOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.CreateOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	geocodingResult, err := a.geocoder.ForwardGeocodeStructured(ctx, geocoding.ForwardGeocodeStructuredInput{})
	if err != nil {
		return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to forward geocode structured", err.Error())}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	orgDocData := form.Value["data"]
	orgFileData := form.File["data"]

	var orgDataBytes []byte

	if len(orgDocData) == 0 && len(orgFileData) == 0 {
		return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide organization json data document")}, nil
	}

	if len(orgDocData) > 0 {
		orgDataBytes = []byte(orgDocData[0])
	}

	if len(orgFileData) > 0 {
		orgFile, err := orgFileData[0].Open()
		if err != nil {
			return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		orgDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.CreateOrganizationBody
	if err := json.Unmarshal(orgDataBytes, &input); err != nil {
		return api.CreateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_internal", "failed to unmarshal non file data")}, nil
	}

	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]

	var (
		photoFileSources []uploadSource
		videoFileSources []uploadSource
	)

	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
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
		photoUploadResults []uploadResult
		videoUploadResults []uploadResult
	)

	if len(photoFileSources) > 0 {
		photoResults := a.uploadOrganizationFiles(ctx, photoFileSources, 5)
		photoUploadResults = append(photoUploadResults, photoResults...)
	}

	if len(videoFileSources) > 0 {
		videoResults := a.uploadOrganizationFiles(ctx, videoFileSources, 5)
		videoUploadResults = append(videoUploadResults, videoResults...)
	}

	var filesToDelete []string

	for _, res := range append(photoUploadResults, videoUploadResults...) {
		if res.Err != nil {
			filesToDelete = append(filesToDelete, res.Name)
		}
	}

	if len(filesToDelete) > 0 {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	var organizationCreateSetter dbtype.OrganizationCreateSetter

	organizationCreateSetter.Organization = models.OrganizationSetter{
		Name:             omit.From(input.Name),
		Website:          omitnull.FromPtr(input.Website),
		MissionStatement: omitnull.FromPtr(input.MissionStatement),
		AdoptionPolicy:   omitnull.FromPtr(input.AdoptionPolicy),
		AdoptionURL:      omitnull.FromPtr(input.AdoptionURL),
		Distance:         omitnull.FromPtr(input.Distance),
		Status:           omit.FromPtr((*string)(input.Status)),
	}

	organizationCreateSetter.Contact = models.OrganizationContactSetter{
		Phone: omit.From(input.Contact.Phone),
		Email: omit.From(string(input.Contact.Email)),
	}

	point := geom.NewPoint(geom.XY).SetSRID(dbcustom.SRID).MustSetCoords(geom.Coord{geocodingResult.Lon, geocodingResult.Lat})
	coords := &ewkb.Point{Point: point}

	organizationCreateSetter.Address = models.AddressSetter{
		CountryID:     omit.From(input.Contact.Address.CountryID),
		StreetAddress: omit.From(input.Contact.Address.StreetAddress),
		City:          omit.From(input.Contact.Address.City),
		Coords:        omitnull.From(coords),
		Note:          omitnull.FromPtr(input.Contact.Address.Note),
		Region:        omitnull.FromPtr(input.Contact.Address.Region),
		PostalCode:    omitnull.FromPtr(input.Contact.Address.PostalCode),
		UnitNumber:    omitnull.FromPtr(input.Contact.Address.UnitNumber),
		StreetNumber:  omitnull.FromPtr(input.Contact.Address.StreetNumber),
	}

	if input.WorkHour != nil {
		if (input.WorkHour.Monday != nil) ||
			(input.WorkHour.Tuesday != nil) ||
			(input.WorkHour.Wednesday != nil) ||
			(input.WorkHour.Thursday != nil) ||
			(input.WorkHour.Friday != nil) ||
			(input.WorkHour.Saturday != nil) ||
			(input.WorkHour.Sunday != nil) {
			organizationCreateSetter.WorkHour = omitnull.From(models.OrganizationWorkHourSetter{
				Monday:    omitnull.FromPtr(input.WorkHour.Monday),
				Tuesday:   omitnull.FromPtr(input.WorkHour.Tuesday),
				Wednesday: omitnull.FromPtr(input.WorkHour.Wednesday),
				Thursday:  omitnull.FromPtr(input.WorkHour.Thursday),
				Friday:    omitnull.FromPtr(input.WorkHour.Friday),
				Saturday:  omitnull.FromPtr(input.WorkHour.Saturday),
				Sunday:    omitnull.FromPtr(input.WorkHour.Sunday),
			})
		}
	}

	if len(photoUploadResults) > 0 {
		organizationPhotoSetters := make([]models.OrganizationPhotoSetter, len(photoUploadResults))
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

		organizationCreateSetter.Photos = omitnull.From(organizationPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		organizationVideoSetters := make([]models.OrganizationVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			organizationVideoSetters[i] = models.OrganizationVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}

		organizationCreateSetter.Videos = omitnull.From(organizationVideoSetters)
	}

	if input.Socials != nil {
		organizationSocialsSetters := make([]models.OrganizationSocialSetter, 0)
		for _, social := range *input.Socials {
			organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
				Platform: omit.From(social.Platform),
				URL:      omit.From(social.URL),
			})
		}

		organizationCreateSetter.Socials = omitnull.From(organizationSocialsSetters)
	}

	organization, err := a.persistor.Organization().CreateOrganization(ctx, organizationCreateSetter)
	if err != nil {
		msg := "could not create an organization"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_save", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_save", msg, reason)}, nil
		}

		return api.CreateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_save", msg, reason)}, nil
	}

	resp := api.CreateOrganization201JSONResponse(dto.OrganizationToResponse(organization))

	// @TODO: use outbox pattern
	if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
		a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
		return api.CreateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to create permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) UpdateOrganization(ctx context.Context, request api.UpdateOrganizationRequestObject) (api.UpdateOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.UpdateOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	form, err := request.Body.ReadForm(maxFileMemory)
	if err != nil {
		return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}

	defer func() { _ = form.RemoveAll() }()

	orgDocData := form.Value["data"]
	orgFileData := form.File["data"]

	var orgDataBytes []byte

	if len(orgDocData) == 0 && len(orgFileData) == 0 {
		return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "data is missing", "provide organization json data document")}, nil
	}

	if len(orgDocData) > 0 {
		orgDataBytes = []byte(orgDocData[0])
	}

	if len(orgFileData) > 0 {
		orgFile, err := orgFileData[0].Open()
		if err != nil {
			return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to open data file", "provide valid data json file")}, nil
		}

		orgDataBytes, err = io.ReadAll(orgFile)
		if err != nil {
			return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read data file", "provide valid data json file")}, nil
		}
	}

	var input api.UpdateOrganizationBody
	if err := json.Unmarshal(orgDataBytes, &input); err != nil {
		return api.UpdateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_internal", "failed to unmarshal non file data")}, nil
	}

	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]

	var (
		photoFileSources []uploadSource
		videoFileSources []uploadSource
	)

	for _, fh := range photoFileHeaders {
		photoFileSources = append(photoFileSources, multipartSource{fh: fh})
	}

	for _, fh := range videoFileHeaders {
		videoFileSources = append(videoFileSources, multipartSource{fh: fh})
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
		photoUploadResults []uploadResult
		videoUploadResults []uploadResult
	)

	if len(photoFileSources) > 0 {
		photoResults := a.uploadOrganizationFiles(ctx, photoFileSources, 5)
		photoUploadResults = append(photoUploadResults, photoResults...)
	}

	if len(videoFileSources) > 0 {
		videoResults := a.uploadOrganizationFiles(ctx, videoFileSources, 5)
		videoUploadResults = append(videoUploadResults, videoResults...)
	}

	var filesToDelete []string

	for _, res := range append(photoUploadResults, videoUploadResults...) {
		if res.Err != nil {
			filesToDelete = append(filesToDelete, res.Name)
		}
	}

	if len(filesToDelete) > 0 {
		go func() {
			_ = a.deleteUploadedFiles(context.Background(), filesToDelete, 5)
		}()

		return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	organizationUpdateSetter := dbtype.OrganizationUpdateSetter{
		Organization: omitnull.From(models.OrganizationSetter{
			Name:             convert.NullableToOmit(input.Name),
			Website:          convert.NullableToOmitNull(input.Website),
			MissionStatement: convert.NullableToOmitNull(input.MissionStatement),
			AdoptionPolicy:   convert.NullableToOmitNull(input.AdoptionPolicy),
			AdoptionURL:      convert.NullableToOmitNull(input.AdoptionURL),
			Distance:         convert.NullableToOmitNull(input.Distance),
		}),
	}
	if input.Contact != nil {
		contactSetter := models.OrganizationContactSetter{
			Phone: convert.NullableToOmit(input.Contact.Phone),
		}
		if input.Contact.Email.IsSpecified() {
			if input.Contact.Email.IsNull() {
				contactSetter.Email = omit.FromPtr[string](nil)
			} else {
				contactSetter.Email = omit.From(string(input.Contact.Email.MustGet()))
			}
		}

		organizationUpdateSetter.Contact = omitnull.From(contactSetter)

		if input.Contact.Address.IsSpecified() && !input.Contact.Address.IsNull() {
			addr := input.Contact.Address.MustGet()
			organizationUpdateSetter.Address.Set(models.AddressSetter{
				// CountryID:    nullable.NewNullableWithValue[int64](),
				UnitNumber:   convert.NullableToOmitNull(addr.UnitNumber),
				StreetNumber: convert.NullableToOmitNull(addr.StreetNumber),
				Region:       convert.NullableToOmitNull(addr.Region),
				PostalCode:   convert.NullableToOmitNull(addr.PostalCode),
				// Coords: xxx, // use geocoder if any addr info of importance changes
				Note:          convert.NullableToOmitNull(addr.Note),
				City:          omit.FromPtr(addr.City),
				StreetAddress: omit.FromPtr(addr.StreetAddress),
			})
		}
	}

	if input.WorkHour.IsSpecified() {
		if input.WorkHour.IsNull() {
			organizationUpdateSetter.WorkHour = omitnull.FromPtr[models.OrganizationWorkHourSetter](nil)
		} else {
			inWorkHour := input.WorkHour.MustGet()
			organizationUpdateSetter.WorkHour = omitnull.From(models.OrganizationWorkHourSetter{
				Monday:    convert.NullableToOmitNull(inWorkHour.Monday),
				Tuesday:   convert.NullableToOmitNull(inWorkHour.Tuesday),
				Wednesday: convert.NullableToOmitNull(inWorkHour.Wednesday),
				Thursday:  convert.NullableToOmitNull(inWorkHour.Thursday),
				Friday:    convert.NullableToOmitNull(inWorkHour.Friday),
				Saturday:  convert.NullableToOmitNull(inWorkHour.Saturday),
				Sunday:    convert.NullableToOmitNull(inWorkHour.Sunday),
			})
		}
	}

	if len(photoUploadResults) > 0 {
		organizationPhotoSetters := make([]models.OrganizationPhotoSetter, len(photoUploadResults))
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

		organizationUpdateSetter.Photos = omitnull.From(organizationPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		organizationVideoSetters := make([]models.OrganizationVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			organizationVideoSetters[i] = models.OrganizationVideoSetter{
				ObjectKind: omit.From(a.uploader.Kind()),
				ObjectRef:  omit.From(videoRes.Name),
			}
		}

		organizationUpdateSetter.Videos = omitnull.From(organizationVideoSetters)
	}

	if input.Socials.IsSpecified() {
		if input.Socials.IsNull() {
			organizationUpdateSetter.Socials = omitnull.FromPtr[[]models.OrganizationSocialSetter](nil)
		} else {
			organizationSocialsSetters := make([]models.OrganizationSocialSetter, 0)

			if input.Socials.IsSpecified() {
				for _, social := range input.Socials.MustGet() {
					organizationSocialsSetters = append(organizationSocialsSetters, models.OrganizationSocialSetter{
						Platform: omit.From(social.Platform),
						URL:      omit.From(social.URL),
					})
				}
			}

			organizationUpdateSetter.Socials = omitnull.From(organizationSocialsSetters)
		}
	}

	organization, err := a.persistor.Organization().UpdateOrganization(ctx, request.ID, organizationUpdateSetter)
	if err != nil {
		msg := "could not update an organization"

		var (
			reason string
			e1     postgres.ErrOrganizationUniqueViolation
		)

		if errors.As(err, &e1) {
			reason = "duplicate " + e1.Name
			return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_edit", msg, reason)}, nil
		}

		var e2 postgres.IntegrityViolationError
		if errors.As(err, &e2) {
			reason = "organization integrity error"
			return api.UpdateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "organization_edit", msg, reason)}, nil
		}

		return api.UpdateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_edit", msg, reason)}, nil
	}

	resp := api.UpdateOrganization201JSONResponse(dto.OrganizationToResponse(organization))

	// @TODO: use outbox pattern
	// if err := createOrganizationRelationTuples(ctx, a.Keto, sess.Identity.Id, resp.ID); err != nil {
	// 	a.Log.Error("failed to insert organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", resp.ID), slog.Any("error", err))
	// 	return api.UpdateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to create permissions")}, nil
	// }

	return resp, nil
}

func (a *ApiHandler) DeleteOrganization(ctx context.Context, request api.DeleteOrganizationRequestObject) (api.DeleteOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "permission denied")}, nil
	}

	oldPhotosResults, err := a.persistor.Organization().ListOrganizationPhotos(ctx, request.ID, dbtype.ListOrganizationPhotosFilters{})
	if err != nil {
		a.Log.Error("failed to get old organization photos", slog.Int64("organization_id", request.ID), slog.Any("error", err))
	}

	oldVideoResults, err := a.persistor.Organization().ListOrganizationVideos(ctx, request.ID, dbtype.ListOrganizationVideosFilters{})
	if err != nil {
		a.Log.Error("failed to get old organization videos", slog.Int64("organization_id", request.ID), slog.Any("error", err))
	}

	var oldMedia []string

	for _, rez := range oldPhotosResults.Data {
		if !rez.ObjectRefFull.IsNull() {
			oldMedia = append(oldMedia, rez.ObjectRefFull.MustGet())
		}
	}

	for _, rez := range oldVideoResults.Data {
		if rez.ObjectRef != "" {
			oldMedia = append(oldMedia, rez.ObjectRef)
		}
	}

	if _, err := a.persistor.Organization().DeleteOrganizationByID(ctx, request.ID); err != nil {
		return nil, fmt.Errorf("failed to delete an organization by id: %w", err)
	}

	resp := api.DeleteOrganization204Response{}

	go func(oldNames []string) {
		if len(oldNames) > 0 {
			_ = a.deleteUploadedFiles(context.Background(), oldNames, len(oldNames))
		}
	}(oldMedia)

	// @TODO: use outbox pattern
	if err := deleteOrganizationRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", request.ID), slog.Any("error", err))
		return api.DeleteOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) DeleteOrganizations(ctx context.Context, request api.DeleteOrganizationsRequestObject) (api.DeleteOrganizationsResponseObject, error) {
	sess := GetSession(ctx)

	if len(request.Body.Ids) == 0 {
		return api.DeleteOrganizations204Response{}, nil
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations", // @TODO: bulk check permissions instead of looping for each id
			Relation:  "manage",
			Subject:   rts.NewSubjectID(shared.AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.DeleteOrganizations403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organizations_permission", "permission denied")}, nil
	}

	if err := a.persistor.Organization().DeleteOrganizations(ctx, request.Body.Ids); err != nil {
		return nil, fmt.Errorf("failed to delete a organizations by ids: %w", err)
	}

	resp := api.DeleteOrganizations204Response{}

	// @TODO: use outbox pattern
	for _, id := range request.Body.Ids {
		if err := deleteOrganizationRelationTuples(ctx, a.Keto, id); err != nil {
			a.Log.Error("failed to delete organizations relation-tuple", slog.Int64("organization_id", id), slog.Any("error", err))
			return api.DeleteOrganizationsdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organizations_permissions", "failed to delete permissions")}, nil
		}
	}

	return resp, nil
}

func (a *ApiHandler) GetOrganization(ctx context.Context, request api.GetOrganizationRequestObject) (api.GetOrganizationResponseObject, error) {
	filters := dbtype.GetOrganizationByIDFilters{GetOrganizationParams: request.Params}

	organizationWithJoinData, err := a.persistor.Organization().GetOrganizationByID(ctx, request.ID, filters)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationNotFound) {
			return api.GetOrganization404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("organization_not_found", "organization not found")}, nil
		}

		return nil, fmt.Errorf("failed to get an organization by id: %w", err)
	}

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    shared.AuthzOrganizationID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.GetOrganizationdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_permission", "permission denied")}, nil
	}

	resp := api.GetOrganization200JSONResponse(dto.OrganizationWithJoinDataToResponse(organizationWithJoinData, a.uploader))

	return resp, nil
}

func (a *ApiHandler) ListOrganizations(ctx context.Context, request api.ListOrganizationsRequestObject) (api.ListOrganizationsResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.GetAllowed() {
		return api.ListOrganizationsdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_permission", "permission denied")}, nil
	}

	filters := dbtype.ListOrganizationsFilters{ListOrganizationsParams: request.Params}
	paginationParams := getPaginationParams(request.Params.Page, request.Params.PageSize)
	filters.Page = &paginationParams.Page
	filters.PageSize = &paginationParams.PageSize

	organizations, err := a.persistor.Organization().ListOrganizations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	organizationsData := make([]api.Organization, len(organizations.Data))
	for i, organizationWithJoinData := range organizations.Data {
		organizationsData[i] = dto.OrganizationWithJoinDataToResponse(organizationWithJoinData, a.uploader)
	}

	resp := api.ListOrganizations200JSONResponse{
		Data: organizationsData,
		Meta: getPaginationMeta(request.Params.Page, request.Params.PageSize, organizations.TotalCount),
	}

	return resp, nil
}

func createOrganizationRelationTuples(ctx context.Context, c *keto.Client, identityID string, organizationID int64) error {
	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Organization",
					Object:    shared.AuthzOrganizationID(organizationID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(shared.AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Organization",
					Object:    shared.AuthzOrganizationID(organizationID),
					Relation:  "parents",
					Subject:   rts.NewSubjectSet("Organizations", "organizations", ""),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to insert organization relation tuples: %w", err)
	}

	return nil
}

func deleteOrganizationRelationTuples(ctx context.Context, c *keto.Client, organizationID int64) error {
	ownersResp, err := c.Read.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		RelationQuery: &rts.RelationQuery{
			Namespace: new("Organization"),
			Object:    new(shared.AuthzOrganizationID(organizationID)),
			Relation:  new("owners"),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list organization relation tuples: %w", err)
	}

	tuplesToDelete := []*rts.RelationTupleDelta{
		{
			Action: rts.RelationTupleDelta_ACTION_DELETE,
			RelationTuple: &rts.RelationTuple{
				Namespace: "Organization",
				Object:    shared.AuthzOrganizationID(organizationID),
				Relation:  "parents",
				Subject:   rts.NewSubjectSet("Organizations", "organizations", ""),
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
				Namespace: "Organization",
				Object:    shared.AuthzOrganizationID(organizationID),
				Relation:  "owners",
				Subject:   rts.NewSubjectID(subject.GetId()),
			},
		})
	}

	if _, err := c.Write.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuplesToDelete,
	}); err != nil {
		return fmt.Errorf("failed to delete organization relation tuples: %w", err)
	}

	return nil
}
