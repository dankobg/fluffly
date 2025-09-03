package server

import (
	"context"
	"fmt"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/dbg"
)

func (a *ApiHandler) CreateAnimal(ctx context.Context, request api.CreateAnimalRequestObject) (api.CreateAnimalResponseObject, error) {
	dbg.PrintJSON(request.Body)
	return api.CreateAnimal201JSONResponse{}, nil
}

func (a *ApiHandler) UpdateAnimal(ctx context.Context, request api.UpdateAnimalRequestObject) (api.UpdateAnimalResponseObject, error) {
	dbg.PrintJSON(request.Body)
	return api.UpdateAnimal201JSONResponse{}, nil
}

func (a *ApiHandler) DeleteAnimal(ctx context.Context, request api.DeleteAnimalRequestObject) (api.DeleteAnimalResponseObject, error) {
	fmt.Println(request.ID)
	return api.DeleteAnimal204Response{}, nil
}

func (a *ApiHandler) ListAnimals(ctx context.Context, request api.ListAnimalsRequestObject) (api.ListAnimalsResponseObject, error) {
	return api.ListAnimals200JSONResponse{}, nil
}

func (a *ApiHandler) GetAnimal(ctx context.Context, request api.GetAnimalRequestObject) (api.GetAnimalResponseObject, error) {
	return api.GetAnimal200JSONResponse{}, nil
}
