package nominatim

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dankobg/fluffly/geocoding"
	"golang.org/x/time/rate"
)

const baseURL = "https://nominatim.openstreetmap.org"

type NominatimGeocoder struct {
	c       *http.Client
	baseURL string
	limiter *rate.Limiter
}

func NewNominatimGeocoder(c *http.Client) (*NominatimGeocoder, error) {
	if c == nil {
		c = http.DefaultClient
	}

	return &NominatimGeocoder{
		c:       c,
		baseURL: baseURL,
		limiter: rate.NewLimiter(rate.Every(time.Second+time.Millisecond*50), 1),
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

func (n *NominatimGeocoder) ForwardGeocode(ctx context.Context, search string) (geocoding.ForwardGeocodeResult, error) {
	if search == "" {
		return geocoding.ForwardGeocodeResult{}, errors.New("search query can't be empty")
	}

	if err := n.limiter.Wait(ctx); err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("rate limiter wait: %w", err)
	}

	u, _ := url.Parse(n.baseURL)
	searchURL := u.JoinPath("search")
	q := searchURL.Query()
	q.Add("format", "jsonv2")
	q.Add("q", search)
	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to create a search request: %w", err)
	}

	resp, err := n.c.Do(req)
	if err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to get geocoding results: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

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

	lon, e2 := strconv.ParseFloat(bestMatch.Lon, 64)
	if err := errors.Join(e1, e2); err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to parse latitude/longitude: %w", err)
	}

	return geocoding.ForwardGeocodeResult{Lat: lat, Lon: lon}, nil
}

func (n *NominatimGeocoder) ForwardGeocodeStructured(ctx context.Context, search geocoding.ForwardGeocodeStructuredInput) (geocoding.ForwardGeocodeResult, error) {
	if search.Country == nil && search.City == nil && search.StreetAddress == nil && search.PostalCode == nil {
		return geocoding.ForwardGeocodeResult{}, errors.New("search can't be empty")
	}

	if err := n.limiter.Wait(ctx); err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("rate limiter wait: %w", err)
	}

	u, _ := url.Parse(n.baseURL)
	searchURL := u.JoinPath("search")
	q := searchURL.Query()
	q.Add("format", "jsonv2")

	if search.Country != nil && strings.TrimSpace(*search.Country) != "" {
		q.Add("country", *search.Country)
	}

	if search.City != nil && strings.TrimSpace(*search.City) != "" {
		q.Add("city", *search.City)
	}

	if search.StreetAddress != nil && strings.TrimSpace(*search.StreetAddress) != "" {
		q.Add("street", *search.StreetAddress)
	}

	if search.PostalCode != nil && strings.TrimSpace(*search.PostalCode) != "" {
		q.Add("postalcode", *search.PostalCode)
	}

	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to create a search request: %w", err)
	}

	resp, err := n.c.Do(req)
	if err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to get geocoding results: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

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

	lon, e2 := strconv.ParseFloat(bestMatch.Lon, 64)
	if err := errors.Join(e1, e2); err != nil {
		return geocoding.ForwardGeocodeResult{}, fmt.Errorf("failed to parse latitude/longitude: %w", err)
	}

	return geocoding.ForwardGeocodeResult{Lat: lat, Lon: lon}, nil
}

var (
	reNonNumeric  = regexp.MustCompile(`[^a-z0-9\s,]+`)
	reMultiSpaces = regexp.MustCompile(`\s+`)
	reMultiCommas = regexp.MustCompile(`,+`)
)

func NormalizeSearchQuery(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// remove everything except letters, digits, spaces, commas
	s = reNonNumeric.ReplaceAllString(s, "")

	// replace multiple spaces with a single space
	s = reMultiSpaces.ReplaceAllString(s, " ")

	// replace multiple commas with a single comma
	s = reMultiCommas.ReplaceAllString(s, ",")

	// remove spaces before commas
	s = strings.ReplaceAll(s, " ,", ",")

	// ensure exactly one space after each comma
	s = strings.ReplaceAll(s, ",", ", ")

	// collapse any accidental double spaces created by the previous step
	s = reMultiSpaces.ReplaceAllString(s, " ")

	// trim leading/trailing spaces
	s = strings.TrimSpace(s)

	return s
}
