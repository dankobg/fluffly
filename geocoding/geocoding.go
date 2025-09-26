package geocoding

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
	ForwardGeocode(search string) (ForwardGeocodeResult, error)
	ForwardGeocodeStructured(search ForwardGeocodeStructuredInput) (ForwardGeocodeResult, error)
}
