package validators

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func DefineCustomOpenapiFormatValidators() {
	defineBodyDecoders()
	openapi3.DefineStringFormatValidator("uri", NewURIValidator())
	openapi3.DefineStringFormatValidator("email", openapi3.NewRegexpFormatValidator(openapi3.FormatOfStringForEmail))
}

func defineBodyDecoders() {
	openapi3filter.RegisterBodyDecoder("image/jpeg", openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder("image/png", openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder("video/mp4", openapi3filter.FileBodyDecoder)
}
