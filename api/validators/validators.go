package validators

import "github.com/getkin/kin-openapi/openapi3"

func DefineCustomOpenapiFormatValidators() {
	openapi3.DefineStringFormatValidator("uri", NewURIValidator())
	openapi3.DefineStringFormatValidator("email", openapi3.NewRegexpFormatValidator(openapi3.FormatOfStringForEmail))
}
