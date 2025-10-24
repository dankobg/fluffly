package server

import (
	"context"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) GetHealthAlive(ctx context.Context, request api.GetHealthAliveRequestObject) (api.GetHealthAliveResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Health",
			Object:    "alive",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetHealthAlivedefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("health_permission", "invalid permission")}, nil
	}

	return api.GetHealthAlive200JSONResponse{Alive: true}, nil
}

func (a *ApiHandler) GetHealthReady(ctx context.Context, request api.GetHealthReadyRequestObject) (api.GetHealthReadyResponseObject, error) {
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Health",
			Object:    "ready",
			Relation:  "view",
			Subject:   rts.NewSubjectID("*"),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetHealthReadydefaultJSONResponse{StatusCode: http.StatusForbidden, Body: newUnauthorizedErr("health_permission", "invalid permission")}, nil
	}

	return api.GetHealthReady200JSONResponse{Ready: true}, nil
}
