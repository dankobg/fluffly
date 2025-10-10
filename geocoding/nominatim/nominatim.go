package nominatim

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/dankobg/fluffly/geocoding"
)

const baseURL = "https://nominatim.openstreetmap.org"

type NominatimGeocoder struct {
	c       *http.Client
	baseURL string
}

func NewNominatimGeocoder(c *http.Client) (*NominatimGeocoder, error) {
	if c == nil {
		c = http.DefaultClient
	}
	return &NominatimGeocoder{
		c:       c,
		baseURL: baseURL,
	}, nil
}

type nominatimResult struct {
	PlaceID     int     `json:"place_id,omitempty"`
	Licence     string  `json:"licence,omitempty"`
	OsmType     string  `json:"osm_type,omitempty"`
	OsmID       int     `json:"osm_id,omitempty"`
	Lat         string  `json:"lat,omitempty"`
	Lon         string  `json:"lon,omitempty"`
	Category    string  `json:"category,omitempty"`
	Type        string  `json:"type,omitempty"`
	PlaceRank   int     `json:"place_rank,omitempty"`
	Importance  float64 `json:"importance,omitempty"`
	Addresstype string  `json:"addresstype,omitempty"`
	Name        string  `json:"name,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	Address     struct {
		Shop         string `json:"shop,omitempty"`
		Road         string `json:"road,omitempty"`
		Quarter      string `json:"quarter,omitempty"`
		Suburb       string `json:"suburb,omitempty"`
		Borough      string `json:"borough,omitempty"`
		City         string `json:"city,omitempty"`
		ISO31662Lvl4 string `json:"ISO3166-2-lvl4,omitempty"`
		Postcode     string `json:"postcode,omitempty"`
		Country      string `json:"country,omitempty"`
		CountryCode  string `json:"country_code,omitempty"`
	} `json:"address,omitzero"`
	Boundingbox []string `json:"boundingbox,omitempty"`
}

func (n *NominatimGeocoder) ForwardGeocode(search string) (geocoding.ForwardGeocodeResult, error) {
	if search == "" {
		return geocoding.ForwardGeocodeResult{}, errors.New("search query can't be empty")
	}

	u, _ := url.Parse(n.baseURL)
	searchURL := u.JoinPath("search")
	q := searchURL.Query()
	q.Add("q", search)
	q.Add("format", "jsonv2")
	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to create a search request: %w", err)
	}
	resp, err := n.c.Do(req)
	if err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to get geocoding results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("received non ok status code: %d", resp.StatusCode)
	}

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to decode geocoding results: %w", err)
	}
	if len(results) == 0 {
		return geocoding.ForwardGeocodeResult{}, errors.New("no results found")
	}

	bestMatch := slices.MaxFunc(results, func(a, b nominatimResult) int {
		return cmp.Compare(a.Importance, b.Importance)
	})

	lat, e1 := strconv.ParseFloat(bestMatch.Lat, 64)
	lng, e2 := strconv.ParseFloat(bestMatch.Lon, 64)
	if err := errors.Join(e1, e2); err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to parse latitude/longitude: %w", err)
	}

	return geocoding.ForwardGeocodeResult{Lat: lat, Lng: lng}, nil
}

func (n *NominatimGeocoder) ForwardGeocodeStructured(search geocoding.ForwardGeocodeStructuredInput) (geocoding.ForwardGeocodeResult, error) {
	panic("not implemented")
}
