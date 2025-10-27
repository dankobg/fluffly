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
	"github.com/dankobg/fluffly/dto"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/ptr"
	"github.com/oapi-codegen/nullable"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

const (
	createOrganizationFileMaxMemory = 50 << 20
)

func (a *ApiHandler) CreateOrganization(ctx context.Context, request api.CreateOrganizationRequestObject) (api.CreateOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organizations",
			Object:    "organizations",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.CreateOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "invalid permission")}, nil
	}

	form, err := request.Body.ReadForm(createOrganizationFileMaxMemory)
	if err != nil {
		return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "failed to read a form data", err.Error())}, nil
	}
	defer form.RemoveAll()
	orgData := form.Value["data"][0]
	var input api.CreateOrganizationBody
	if err := json.Unmarshal([]byte(orgData), &input); err != nil {
		return api.CreateOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_internal", "failed to unmarshal non file data")}, nil
	}

	photoFileHeaders := form.File["photos"]
	videoFileHeaders := form.File["videos"]
	var photoFileSources []uploadSource
	var videoFileSources []uploadSource
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
	var photoUploadResults []uploadResult
	var videoUploadResults []uploadResult
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
			_ = a.deleteUploadedFiles(ctx, filesToDelete, 5)
		}()
		return api.CreateOrganization400JSONResponse{GenericErrorResponseJSONResponse: newGenericResp(http.StatusBadRequest, "file_upload", "failed to upload files")}, nil
	}

	var organizationCreateSetter dbtype.OrganizationCreateSetter

	organizationCreateSetter.Organization = dbtype.OrganizationSetter{
		Name:             nullable.NewNullableWithValue(input.Name),
		Website:          input.Website,
		MissionStatement: input.MissionStatement,
		AdoptionPolicy:   input.AdoptionPolicy,
		AdoptionURL:      input.AdoptionURL,
		Distance:         input.Distance,
	}

	organizationCreateSetter.Contact = dbtype.OrganizationContactSetter{
		Phone: nullable.NewNullableWithValue(input.Contact.Phone),
		Email: nullable.NewNullableWithValue(string(input.Contact.Email)),
	}

	organizationCreateSetter.Address = dbtype.AddressSetter{
		CountryID:     nullable.NewNullableWithValue(input.Contact.Address.CountryID),
		UnitNumber:    nullable.NewNullableWithValue(*input.Contact.Address.UnitNumber),
		StreetNumber:  nullable.NewNullableWithValue(*input.Contact.Address.StreetNumber),
		StreetAddress: nullable.NewNullableWithValue(input.Contact.Address.StreetAddress),
		City:          nullable.NewNullableWithValue(input.Contact.Address.City),
		Region:        nullable.NewNullableWithValue(*input.Contact.Address.Region),
		PostalCode:    nullable.NewNullableWithValue(*input.Contact.Address.PostalCode),
		Lat:           nullable.NewNullableWithValue(float64(*input.Contact.Address.Lat)),
		Lng:           nullable.NewNullableWithValue(float64(*input.Contact.Address.Lng)),
		Note:          nullable.NewNullableWithValue(*input.Contact.Address.Note),
	}

	if input.WorkHour != nil {
		organizationCreateSetter.WorkHour = nullable.NewNullableWithValue(dbtype.OrganizationWorkHourSetter{
			Monday:    input.WorkHour.Monday,
			Tuesday:   input.WorkHour.Tuesday,
			Wednesday: input.WorkHour.Wednesday,
			Thursday:  input.WorkHour.Thursday,
			Friday:    input.WorkHour.Friday,
			Saturday:  input.WorkHour.Saturday,
			Sunday:    input.WorkHour.Sunday,
		})
	}

	if len(photoUploadResults) > 0 {
		organizationPhotoSetters := make([]dbtype.OrganizationPhotoSetter, len(photoUploadResults))
		for i, photoRes := range photoUploadResults {
			// @TODO: optimize image size
			photoSetter := dbtype.OrganizationPhotoSetter{
				ObjectKind:      nullable.NewNullableWithValue(a.uploader.Kind()),
				ObjectRefSmall:  nullable.NewNullableWithValue(photoRes.Name),
				ObjectRefMedium: nullable.NewNullableWithValue(photoRes.Name),
				ObjectRefLarge:  nullable.NewNullableWithValue(photoRes.Name),
				ObjectRefFull:   nullable.NewNullableWithValue(photoRes.Name),
			}
			organizationPhotoSetters[i] = photoSetter
		}
		organizationCreateSetter.Photos = nullable.NewNullableWithValue(organizationPhotoSetters)
	}

	if len(videoUploadResults) > 0 {
		organizationVideoSetters := make([]dbtype.OrganizationVideoSetter, len(videoUploadResults))
		for i, videoRes := range videoUploadResults {
			videoSetter := dbtype.OrganizationVideoSetter{
				ObjectKind: nullable.NewNullableWithValue(a.uploader.Kind()),
				ObjectRef:  nullable.NewNullableWithValue(videoRes.Name),
			}
			organizationVideoSetters[i] = videoSetter
		}
		organizationCreateSetter.Videos = nullable.NewNullableWithValue(organizationVideoSetters)
	}

	if input.Socials.IsSpecified() && !input.Socials.IsNull() {
		organizationSocialsSetters := make([]dbtype.OrganizationSocialSetter, 0)
		if input.Socials.IsSpecified() {
			for _, social := range input.Socials.MustGet() {
				organizationSocialsSetters = append(organizationSocialsSetters, dbtype.OrganizationSocialSetter{
					Platform: nullable.NewNullableWithValue(social.Platform),
					URL:      nullable.NewNullableWithValue(social.URL),
				})
			}
		}
		organizationCreateSetter.Socials = nullable.NewNullableWithValue(organizationSocialsSetters)
	}

	organization, err := a.persistor.Organization().CreateOrganization(ctx, organizationCreateSetter)
	if err != nil {
		msg := "could not create an organzation"
		var reason string
		var e1 postgres.ErrOrganizationUniqueViolation
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
			Object:    AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.UpdateOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "invalid permission")}, nil
	}

	organizationSetter := dbtype.OrganizationSetter{
		Name:             request.Body.Name,
		Website:          request.Body.Website,
		MissionStatement: request.Body.MissionStatement,
		AdoptionPolicy:   request.Body.AdoptionPolicy,
		AdoptionURL:      request.Body.AdoptionURL,
		Distance:         request.Body.Distance,
	}

	organization, err := a.persistor.Organization().UpdateOrganization(ctx, request.ID, organizationSetter)
	if err != nil {
		msg := "could not update an organzation"
		var reason string
		var e1 postgres.ErrOrganizationUniqueViolation
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
	return resp, nil
}

func (a *ApiHandler) DeleteOrganization(ctx context.Context, request api.DeleteOrganizationRequestObject) (api.DeleteOrganizationResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    AuthzOrganizationID(request.ID),
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DeleteOrganization403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("organization_permission", "invalid permission")}, nil
	}

	_, err := a.persistor.Organization().DeleteOrganizationByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete an organization by id: %w", err)
	}
	resp := api.DeleteOrganization204Response{}

	// @TODO: use outbox pattern
	if err := deleteOrganizationRelationTuples(ctx, a.Keto, request.ID); err != nil {
		a.Log.Error("failed to delete organization relation-tuple", slog.String("identity_id", sess.Identity.Id), slog.Int64("organization_id", request.ID), slog.Any("error", err))
		return api.DeleteOrganizationdefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: newGenericErr(http.StatusInternalServerError, "organization_permissions", "failed to delete permissions")}, nil
	}

	return resp, nil
}

func (a *ApiHandler) GetOrganization(ctx context.Context, request api.GetOrganizationRequestObject) (api.GetOrganizationResponseObject, error) {
	organizationWithJoinData, err := a.persistor.Organization().GetOrganizationByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrOrganizationNotFound) {
			return api.GetOrganization404JSONResponse{NotFoundErrorResponseJSONResponse: newNotFoundResp("organization_not_found", "organization not found")}, nil
		}
		return nil, fmt.Errorf("failed to get an organization by id: %w", err)
	}
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Organization",
			Object:    AuthzOrganizationID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetOrganizationdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_permission", "invalid permission")}, nil
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
	}); err != nil || !checkResp.Allowed {
		return api.ListOrganizationsdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: newGenericErr(http.StatusUnauthorized, "organization_permission", "invalid permission")}, nil
	}

	var filters dbtype.OrganizationFilters
	filters.Pagination = ptr.Of(getPaginationParams(request.Params.Page, request.Params.PageSize))
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
					Object:    AuthzOrganizationID(organizationID),
					Relation:  "owners",
					Subject:   rts.NewSubjectID(AuthzIdentityID(identityID)),
				},
			},
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: "Organization",
					Object:    AuthzOrganizationID(organizationID),
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
			Namespace: ptr.Of("Organization"),
			Object:    ptr.Of(AuthzOrganizationID(organizationID)),
			Relation:  ptr.Of("owners"),
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
				Object:    AuthzOrganizationID(organizationID),
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
				Object:    AuthzOrganizationID(organizationID),
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
