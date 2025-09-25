package server

import (
	"context"
	"log/slog"
	"net/http"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dto"
)

func (a *ApiHandler) ListSessions(ctx context.Context, request api.ListSessionsRequestObject) (api.ListSessionsResponseObject, error) {
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
	sessions, _, err := req.Execute()
	if err != nil {
		return api.ListSessions400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: 400, Message: err.Error()}}, nil
	}
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

func (a *ApiHandler) DisableSession(ctx context.Context, request api.DisableSessionRequestObject) (api.DisableSessionResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.DisableSession(ctx, request.ID)
	_, err := req.Execute()
	if err != nil {
		return api.DisableSession400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: err.Error()}}, nil
	}
	return api.DisableSession204Response{}, nil
}

func (a *ApiHandler) GetSession(ctx context.Context, request api.GetSessionRequestObject) (api.GetSessionResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.GetSession(ctx, request.ID)
	if request.Params.Expand != nil {
		expands := make([]string, 0, len(*request.Params.Expand))
		for _, x := range *request.Params.Expand {
			expands = append(expands, string(x))
		}
		req = req.Expand(expands)
	}
	session, _, err := req.Execute()
	if err != nil {
		return api.GetSession400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: err.Error()}}, nil
	}
	resp, err := dto.SessionToResponse(*session)
	if err != nil {
		return nil, err
	}
	return api.GetSession200JSONResponse(resp), nil
}

func (a *ApiHandler) ExtendSession(ctx context.Context, request api.ExtendSessionRequestObject) (api.ExtendSessionResponseObject, error) {
	req := a.Kratos.Admin.IdentityAPI.ExtendSession(ctx, request.ID)
	session, _, err := req.Execute()
	if err != nil {
		return api.ExtendSession400JSONResponse{GenericErrorJSONResponse: api.GenericErrorJSONResponse{Code: http.StatusBadRequest, Message: err.Error()}}, nil
	}
	resp, err := dto.SessionToResponse(*session)
	if err != nil {
		return nil, err
	}
	return api.ExtendSession200JSONResponse(resp), nil
}
