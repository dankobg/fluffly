package fakegeo

import (
	"context"

	"github.com/dankobg/fluffly/geocoding"
	"github.com/go-faker/faker/v4"
)

type FakeGeocoder struct{}

func (f FakeGeocoder) ForwardGeocode(ctx context.Context, search string) (geocoding.ForwardGeocodeResult, error) {
	return geocoding.ForwardGeocodeResult{Lat: faker.Latitude(), Lon: faker.Longitude()}, nil
}

func (f FakeGeocoder) ForwardGeocodeStructured(ctx context.Context, search geocoding.ForwardGeocodeStructuredInput) (geocoding.ForwardGeocodeResult, error) {
	return geocoding.ForwardGeocodeResult{Lat: faker.Latitude(), Lon: faker.Longitude()}, nil
}
