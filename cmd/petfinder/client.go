package petfinder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dankobg/fluffly/config"
)

const (
	baseURL = "https://api.petfinder.com/v2"
)

type petfinderClient struct {
	c           *http.Client
	baseURL     string
	apiKey      string
	apiSecret   string
	accessToken string
}

func newPetfinderClient(c *http.Client, cfg config.PetfinderConfig) (*petfinderClient, error) {
	if c == nil {
		c = http.DefaultClient
	}
	return &petfinderClient{
		c:           c,
		baseURL:     baseURL,
		apiKey:      cfg.ApiKey,
		apiSecret:   cfg.ApiSecret,
		accessToken: cfg.AccessToken,
	}, nil
}

type authResult struct {
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

func (p *petfinderClient) authenticate(ctx context.Context) (authResult, error) {
	u, _ := url.Parse(p.baseURL)
	tokenURL := u.JoinPath("oauth2/token")

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", p.apiKey)
	form.Add("client_secret", p.apiSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return authResult{}, fmt.Errorf("failed to create a access token request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.c.Do(req)
	if err != nil {
		return authResult{}, fmt.Errorf("failed to get access token result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return authResult{}, fmt.Errorf("received non ok status code: %d", resp.StatusCode)
	}

	var result authResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return authResult{}, fmt.Errorf("failed to decode access token result: %w", err)
	}

	return result, nil
}

// listOrganizationsQueryParams defines the query parameters for the list organizations endpoint
type listOrganizationsQueryParams struct {
	// Return results matching organization name
	Name *string

	// Return results by location.
	//
	// city, state; latitude,longitude; or postal code
	Location *string

	// Return results within distance of location (in miles).
	//
	// requires location to be set (default: 100, max: 500)
	Distance *int

	// Filter results by state
	//
	// Accepts two-letter abbreviations, e.g. AL, WY
	State *string

	// Filter results by country.
	//
	// Accepts two-letter abbreviations, e.g. US, CA
	Country *string

	// Search on name, city, state (Search that includes synonyms, varying punctuation, etc.)
	Query *string

	// Field to sort by; leading dash requests a reverse-order sort
	//
	// distance, -distance, name, -name, country, -country, state, -state
	Sort *string

	// Maximum number of results to return
	//
	// (default: 20, max: 100)
	Limit *int

	// Page of results to return
	//
	// (default: 1)
	Page *int
}

func (p *petfinderClient) listOrganizations(ctx context.Context, params listOrganizationsQueryParams) (*http.Response, error) {
	u, _ := url.Parse(p.baseURL)
	orgsURL := u.JoinPath("organizations")
	q := orgsURL.Query()

	if params.Name != nil {
		q.Add("name", *params.Name)
	}
	if params.Location != nil {
		q.Add("location", *params.Location)
	}
	if params.Distance != nil {
		q.Add("distance", strconv.Itoa(*params.Distance))
	}
	if params.State != nil {
		q.Add("state", *params.State)
	}
	if params.Country != nil {
		q.Add("country", *params.Country)
	}
	if params.Query != nil {
		q.Add("query", *params.Query)
	}
	if params.Sort != nil {
		q.Add("sort", *params.Sort)
	}
	if params.Limit != nil {
		q.Add("limit", strconv.Itoa(*params.Limit))
	}
	if params.Page != nil {
		q.Add("page", strconv.Itoa(*params.Page))
	}
	orgsURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, orgsURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create organizations request request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+p.accessToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := p.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token result: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("received non ok status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (p *petfinderClient) downloadAnimals(ctx context.Context) error {
	u, _ := url.Parse(p.baseURL)
	orgsURL := u.JoinPath("animals")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, orgsURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create animals request request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+p.accessToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := p.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get access token result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non ok status code: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode animals result: %w", err)
	}

	return nil
}

// listAnimalsQueryParams defines the query parameters for the list animals endpoint
type listAnimalsQueryParams struct {
	// Return results matching animal type.
	//
	// Possible values may be looked up via Get Animal Types.
	Type *string

	// Return results matching animal breed.
	//
	// Accepts multiple values (comma-separated), e.g. "pug,samoyed".
	// Possible values may be looked up via Get Animal Breeds.
	Breed *string

	// Return results matching animal size.
	//
	// Values: "small", "medium", "large", "xlarge".
	// Accepts multiple values comma-separated.
	Size *string

	// Return results matching animal gender.
	//
	// Values: "male", "female", "unknown". Accepts multiple values comma-separated.
	Gender *string

	// Return results matching animal age.
	//
	// Values: "baby", "young", "adult", "senior". Accepts multiple values comma-separated.
	Age *string

	// Return results matching animal color.
	//
	// Possible values may be looked up via Get Animal Types.
	Color *string

	// Return results matching animal coat.
	//
	// Values: "short", "medium", "long", "wire", "hairless", "curly".
	// Accepts multiple values comma-separated.
	Coat *string

	// Return results matching adoption status.
	//
	// Values: "adoptable", "adopted", "found". Default: "adoptable".
	Status *string

	// Return results matching animal name (includes partial matches).
	Name *string

	// Return results associated with specific organization ID(s).
	//
	// Accepts multiple values comma-separated.
	Organization *string

	// Return results that are good with children.
	//
	// Boolean (true/false) or numeric (1/0).
	GoodWithChildren *bool

	// Return results that are good with dogs.
	//
	// Boolean (true/false) or numeric (1/0).
	GoodWithDogs *bool

	// Return results that are good with cats.
	//
	// Boolean (true/false) or numeric (1/0).
	GoodWithCats *bool

	// Return results that are house-trained.
	//
	// Boolean (true only) or numeric (1).
	HouseTrained *bool

	// Return results that are declawed.
	//
	// Boolean (true only) or numeric (1).
	Declawed *bool

	// Return results that have special needs.
	//
	// Boolean (true only) or numeric (1).
	SpecialNeeds *bool

	// Return results by location.
	//
	// Format: "city, state"; "latitude,longitude"; or postal code.
	Location *string

	// Return results within distance (in miles) of the given location.
	//
	// Requires Location to be set. Default: 100. Max: 500.
	Distance *int

	// Return results published before this date/time.
	//
	// Must be a valid ISO-8601 date/time string, e.g. "2019-10-07T19:13:01+00:00".
	Before *string

	// Return results published after this date/time.
	//
	// Must be a valid ISO-8601 date/time string, e.g. "2019-10-07T19:13:01+00:00".
	After *string

	// Attribute to sort by; a leading dash requests reverse-order sort.
	//
	// Values: "recent", "-recent", "distance", "-distance".
	Sort *string

	// Specifies which page of results to return.
	//
	// Default: 1.
	Page *int

	// Maximum number of results to return per 'page'.
	//
	// Default: 20. Max: 100.
	Limit *int
}
