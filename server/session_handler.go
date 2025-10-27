package server

import (
	"context"
	"log/slog"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func (a *ApiHandler) ListSessions(ctx context.Context, request api.ListSessionsRequestObject) (api.ListSessionsResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ListSessions403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.ListSessions(ctx)
	if request.Params.Active != nil {
		req = req.Active(*request.Params.Active)
	}
	if request.Params.PageSize != nil {
		req = req.PageSize(*request.Params.PageSize)
	}
	if request.Params.PageToken != nil && *request.Params.PageToken != "1" {
		req = req.PageToken(*request.Params.PageToken)
	}
	if request.Params.Expand != nil {
		expands := make([]string, 0, len(*request.Params.Expand))
		for _, x := range *request.Params.Expand {
			expands = append(expands, string(x))
		}
		req = req.Expand(expands)
	}
	sessions, sessionsResp, err := req.Execute()
	if err != nil {
		return api.ListSessionsdefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "sessions_list", "failed to list sessions")}, nil
	}
	defer sessionsResp.Body.Close()
	resp := make(api.ListSessions200JSONResponse, 0, len(sessions))
	for _, session := range sessions {
		res, err := dto.SessionToResponse(session)
		if err != nil {
			a.Log.Error("failed to convert session to response", slog.Any("error", err))
			return nil, err
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func (a *ApiHandler) GetSession(ctx context.Context, request api.GetSessionRequestObject) (api.GetSessionResponseObject, error) {
	sess := GetSession(ctx)
	req := a.Kratos.Admin.IdentityAPI.GetSession(ctx, request.ID)
	if request.Params.Expand != nil {
		expands := make([]string, 0, len(*request.Params.Expand))
		for _, x := range *request.Params.Expand {
			expands = append(expands, string(x))
		}
		req = req.Expand(expands)
	}
	session, sessionResp, err := req.Execute()
	if err != nil {
		return api.GetSessiondefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "session_get", "failed to get session")}, nil
	}
	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "Session",
			// Object:    AuthzSessionID(request.ID),
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "view",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.GetSession403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "invalid permission")}, nil
	}
	defer sessionResp.Body.Close()
	resp, err := dto.SessionToResponse(*session)
	if err != nil {
		return nil, err
	}
	return api.GetSession200JSONResponse(resp), nil
}

func (a *ApiHandler) DisableSession(ctx context.Context, request api.DisableSessionRequestObject) (api.DisableSessionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "Session",
			// Object:    AuthzSessionID(request.ID),
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.DisableSession403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.DisableSession(ctx, request.ID)
	disableSessionResp, err := req.Execute()
	if err != nil {
		return api.DisableSessiondefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "session_disable", "failed to disable session")}, nil
	}
	defer disableSessionResp.Body.Close()
	return api.DisableSession204Response{}, nil
}

func (a *ApiHandler) ExtendSession(ctx context.Context, request api.ExtendSessionRequestObject) (api.ExtendSessionResponseObject, error) {
	sess := GetSession(ctx)

	if checkResp, err := a.Keto.Check.Check(ctx, &rts.CheckRequest{
		Tuple: &rts.RelationTuple{
			// Namespace: "Session",
			// Object:    AuthzSessionID(request.ID),
			Namespace: "Sessions",
			Object:    "sessions",
			Relation:  "manage",
			Subject:   rts.NewSubjectID(AuthzIdentityID(sess.Identity.Id)),
		},
	}); err != nil || !checkResp.Allowed {
		return api.ExtendSession403JSONResponse{UnauthorizedErrorResponseJSONResponse: newUnauthorizedResp("session_permission", "invalid permission")}, nil
	}

	req := a.Kratos.Admin.IdentityAPI.ExtendSession(ctx, request.ID)
	session, sessionResp, err := req.Execute()
	if err != nil {
		return api.ExtendSessiondefaultJSONResponse{StatusCode: http.StatusServiceUnavailable, Body: newGenericErr(http.StatusServiceUnavailable, "session_extend", "failed to extend session")}, nil
	}
	defer sessionResp.Body.Close()
	resp, err := dto.SessionToResponse(*session)
	if err != nil {
		return nil, err
	}
	return api.ExtendSession200JSONResponse(resp), nil
}
