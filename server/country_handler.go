package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/aarondl/opt/omit"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/dbmodel"
)

func (a *ApiHandler) CreateCountry(ctx context.Context, request api.CreateCountryRequestObject) (api.CreateCountryResponseObject, error) {
	countrySetter := dbmodel.CountrySetter{
		Name:       omit.From(request.Body.Name),
		IsoAlpha2:  omit.From(request.Body.IsoAlpha2),
		IsoAlpha3:  omit.From(request.Body.IsoAlpha3),
		IsoNumeric: omit.From(request.Body.IsoNumeric),
	}
	country, err := a.persistor.Country().Create(ctx, countrySetter)
	if err != nil {
		return nil, fmt.Errorf("failed to update an country: %w", err)
	}
	resp := api.CreateCountry201JSONResponse(api.Country{
		ID:         country.ID,
		Name:       country.Name,
		IsoAlpha2:  country.IsoAlpha2,
		IsoAlpha3:  country.IsoAlpha3,
		IsoNumeric: country.IsoNumeric,
		CreatedAt:  country.CreatedAt,
		UpdatedAt:  country.UpdatedAt,
	})
	return resp, nil
}

func (a *ApiHandler) UpdateCountry(ctx context.Context, request api.UpdateCountryRequestObject) (api.UpdateCountryResponseObject, error) {
	countrySetter := dbmodel.CountrySetter{}
	if request.Body.Name != nil {
		countrySetter.Name = omit.From(*request.Body.Name)
	}
	if request.Body.IsoAlpha2 != nil {
		countrySetter.IsoAlpha2 = omit.From(*request.Body.IsoAlpha2)
	}
	if request.Body.IsoAlpha3 != nil {
		countrySetter.IsoAlpha3 = omit.From(*request.Body.IsoAlpha3)
	}
	if request.Body.IsoNumeric != nil {
		countrySetter.IsoNumeric = omit.From(*request.Body.IsoNumeric)
	}
	country, err := a.persistor.Country().Update(ctx, request.ID, countrySetter)
	if err != nil {
		return nil, fmt.Errorf("failed to update an country: %w", err)
	}
	resp := api.UpdateCountry201JSONResponse(api.Country{
		ID:         country.ID,
		Name:       country.Name,
		IsoAlpha2:  country.IsoAlpha2,
		IsoAlpha3:  country.IsoAlpha3,
		IsoNumeric: country.IsoNumeric,
		CreatedAt:  country.CreatedAt,
		UpdatedAt:  country.UpdatedAt,
	})
	return resp, nil
}

func (a *ApiHandler) DeleteCountry(ctx context.Context, request api.DeleteCountryRequestObject) (api.DeleteCountryResponseObject, error) {
	_, err := a.persistor.Country().Delete(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete a country by id: %w", err)
	}
	resp := api.DeleteCountry204Response{}
	return resp, nil
}

func (a *ApiHandler) ListCountries(ctx context.Context, request api.ListCountriesRequestObject) (api.ListCountriesResponseObject, error) {
	countries, err := a.persistor.Country().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list countries: %w", err)
	}
	resp := make(api.ListCountries200JSONResponse, len(countries))
	for i, country := range countries {
		resp[i] = api.Country{
			ID:         country.ID,
			Name:       country.Name,
			IsoAlpha2:  country.IsoAlpha2,
			IsoAlpha3:  country.IsoAlpha3,
			IsoNumeric: country.IsoNumeric,
			CreatedAt:  country.CreatedAt,
			UpdatedAt:  country.UpdatedAt,
		}
	}
	return resp, nil
}

func (a *ApiHandler) GetCountry(ctx context.Context, request api.GetCountryRequestObject) (api.GetCountryResponseObject, error) {
	countryRow, err := a.persistor.Country().Get(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.GetCountry404JSONResponse{NotFoundErrorJSONResponse: api.NotFoundErrorJSONResponse{Code: http.StatusNotFound, Message: "Country not found"}}, nil
		}
		return nil, fmt.Errorf("failed to get an country by id: %w", err)
	}
	resp := api.GetCountry200JSONResponse(api.Country{
		ID:         countryRow.ID,
		Name:       countryRow.Name,
		IsoAlpha2:  countryRow.IsoAlpha2,
		IsoAlpha3:  countryRow.IsoAlpha3,
		IsoNumeric: countryRow.IsoNumeric,
		CreatedAt:  countryRow.CreatedAt,
		UpdatedAt:  countryRow.UpdatedAt,
	})
	return resp, nil
}
