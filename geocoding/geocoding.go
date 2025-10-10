package geocoding

import "context"

type ForwardGeocodeResult struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ForwardGeocodeStructuredInput struct {
	Country       *string
	City          *string
	StreetAddress *string
	PostalCode    *string
}

type Geocoder interface {
	ForwardGeocode(ctx context.Context, search string) (ForwardGeocodeResult, error)
	ForwardGeocodeStructured(ctx context.Context, search ForwardGeocodeStructuredInput) (ForwardGeocodeResult, error)
}
