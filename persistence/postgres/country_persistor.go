package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/dbmodel"
	"github.com/dankobg/fluffly/db/queries"
	"github.com/dankobg/fluffly/persistence"
)

var _ persistence.CountryPersistor = (*PgCountryPersistor)(nil)

type PgCountryPersistor struct {
	*PgPersistor
}

func NewPgCountryPersistor(ps *PgPersistor) *PgCountryPersistor {
	return &PgCountryPersistor{
		PgPersistor: ps,
	}
}

func (p *PgCountryPersistor) Create(ctx context.Context, in dbmodel.CountrySetter) (*dbmodel.Country, error) {
	insertedCountry, err := dbmodel.Countries.Insert(&in).One(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to create a country: %w", err)
	}
	return insertedCountry, nil
}

func (p *PgCountryPersistor) Update(ctx context.Context, countryID int64, in dbmodel.CountrySetter) (*dbmodel.Country, error) {
	updatedCountry, err := dbmodel.Countries.Update(in.UpdateMod(), dbmodel.UpdateWhere.Countries.ID.EQ(countryID)).One(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to update a country: %w", err)
	}
	return updatedCountry, nil
}

func (p *PgCountryPersistor) Delete(ctx context.Context, countryID int64) (int64, error) {
	_, err := dbmodel.Countries.Delete(dbmodel.DeleteWhere.Countries.ID.EQ(countryID)).Exec(ctx, p.db)
	if err != nil {
		return 0, fmt.Errorf("failed to delete a country: %w", err)
	}
	return countryID, nil
}

func (p *PgCountryPersistor) Get(ctx context.Context, countryID int64) (*dbmodel.Country, error) {
	country, err := dbmodel.FindCountry(ctx, p.db, countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get a country: %w", err)
	}
	return country, nil
}

func (p *PgCountryPersistor) List(ctx context.Context) (dbmodel.CountrySlice, error) {
	countryRows, err := queries.ListCountries().All(ctx, p.db)
	countries := make([]*dbmodel.Country, len(countryRows))
	if err != nil {
		return countries, fmt.Errorf("failed to list countries: %w", err)
	}
	for i, country := range countryRows {
		countries[i] = &dbmodel.Country{
			ID:         country.ID,
			Name:       country.Name,
			IsoAlpha2:  country.IsoAlpha2,
			IsoAlpha3:  country.IsoAlpha3,
			IsoNumeric: country.IsoNumeric,
			CreatedAt:  country.CreatedAt,
			UpdatedAt:  country.UpdatedAt,
		}
	}
	return countries, nil
}
