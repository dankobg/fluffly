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

var supportedFileTypes = []string{
	"image/jpg",
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
	"image/avif",
	"video/mp4",
	"video/ogg",
	"video/mpeg",
	"video/webm",
}

func defineBodyDecoders() {
	for _, ct := range supportedFileTypes {
		openapi3filter.RegisterBodyDecoder(ct, openapi3filter.FileBodyDecoder)
	}
}
