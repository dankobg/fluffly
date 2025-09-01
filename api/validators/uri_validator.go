package validators

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

// NewURIValidator creates a new FormatValidator that validates the value is an URI
func NewURIValidator() openapi3.FormatValidator[string] {
	return openapi3.NewCallbackValidator(func(uri string) error {
		_, err := url.ParseRequestURI(uri)
		return err
	})
}
