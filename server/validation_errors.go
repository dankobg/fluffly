package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func collectFieldErrors(err error, out *[]api.ValidationDetail) {
	fmt.Println("KIN API ERROR: ", err.Error())

	var me openapi3.MultiError
	if errors.As(err, &me) {
		for _, e := range me {
			collectFieldErrors(e, out)
		}

		return
	}

	var se *openapi3.SchemaError
	if errors.As(err, &se) {
		p := "/" + strings.Join(se.JSONPointer(), "/")

		*out = append(*out, api.ValidationDetail{
			In:      "body",
			Pointer: p,
			Reason:  se.Reason,
		})
		if se.Origin != nil {
			collectFieldErrors(se.Origin, out)
		}

		return
	}

	var re *openapi3filter.RequestError
	if errors.As(err, &re) {
		var in, pointer string
		if re.Parameter != nil {
			in = re.Parameter.In
			pointer = re.Parameter.Name
		}

		reason := re.Reason
		if reason == "" {
			reason = re.Err.Error()
		}

		*out = append(*out, api.ValidationDetail{
			In:      in,
			Pointer: pointer,
			Reason:  reason,
		})
		// if re.Err != nil {
		// 	collectFieldErrors(re.Err, out)
		// }
	}

	var pe *openapi3filter.ParseError
	if errors.As(err, &pe) {
		var in, pointer string

		in = "body"

		path := pe.Path()
		if path != nil && len(path) >= 0 {
			if _, ok := path[0].(string); ok {
				pointer = "/" + path[0].(string)
			}
		}

		reason := pe.Cause.Error()

		*out = append(*out, api.ValidationDetail{
			In:      in,
			Pointer: pointer,
			Reason:  reason,
		})
		if pe.RootCause() != nil {
			collectFieldErrors(pe.RootCause(), out)
		}
	}
}

func handleMultiError(me openapi3.MultiError) (int, error) {
	var details []api.ValidationDetail
	collectFieldErrors(me, &details)

	valErr := api.GenericErrorResponse{
		Code:       "validation",
		Message:    "validation failed",
		Reason:     new("Request failed schema validation"),
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}

	bb, err := json.MarshalIndent(&valErr, "", "  ")
	if err != nil {
		return 500, err
	}

	errResp := fmt.Errorf("%s", string(bb))

	return int(valErr.StatusCode), errResp
}
