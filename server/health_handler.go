package server

import (
	"context"

	api "github.com/dankobg/fluffly/api/gen"
)

func (a *ApiHandler) GetHealthAlive(ctx context.Context, request api.GetHealthAliveRequestObject) (api.GetHealthAliveResponseObject, error) {
	return api.GetHealthAlive200JSONResponse{Alive: true}, nil
}

func (a *ApiHandler) GetHealthReady(ctx context.Context, request api.GetHealthReadyRequestObject) (api.GetHealthReadyResponseObject, error) {
	return api.GetHealthReady200JSONResponse{Ready: true}, nil
}
