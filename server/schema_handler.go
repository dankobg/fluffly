package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListIdentitySchemas(ctx context.Context, request api.ListIdentitySchemasRequestObject) (api.ListIdentitySchemasResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Schemas",
			Object:    "schemas",
			Relation:  "view",
			Subject:   rts.NewSubjectID(sess.Identity.Id),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ListIdentitySchemasdefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: api.Error{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.ListIdentitySchemas(ctx)
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}
	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}
	schemaContainers, schemaContainersResp, err := req.Execute()
	if err != nil {
		return api.ListIdentitySchemasdefaultJSONResponse{Body: api.Error{Code: http.StatusBadRequest, Message: err.Error()}}, nil
	}
	defer schemaContainersResp.Body.Close()
	resp := make(api.ListIdentitySchemas200JSONResponse, 0, len(schemaContainers))
	for _, sc := range schemaContainers {
		res, err := dto.SchemaContainerToResponse(sc)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func (a *ApiHandler) GetIdentitySchema(ctx context.Context, request api.GetIdentitySchemaRequestObject) (api.GetIdentitySchemaResponseObject, error) {
	sess := GetSession(ctx)
	req := a.Kratos.Admin.IdentityAPI.GetIdentitySchema(ctx, request.ID)
	identitySchema, identitySchemaResp, err := req.Execute()
	if err != nil {
		return api.GetIdentitySchema404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "schema not found"}}, nil
	}
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Schema",
			Object:    authzSchemaID(request.ID),
			Relation:  "view",
			Subject:   rts.NewSubjectID(sess.Identity.Id),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetIdentitySchemadefaultJSONResponse{StatusCode: http.StatusUnauthorized, Body: api.Error{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}}, nil
	}
	defer identitySchemaResp.Body.Close()
	return api.GetIdentitySchema200JSONResponse(identitySchema), nil
}
