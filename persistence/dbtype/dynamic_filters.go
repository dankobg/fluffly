package dbtype

import (
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type FilterType int

const (
	FilterTypeString FilterType = iota
	FilterTypeBool
	FilterTypeInt
)

type DynamicFilter struct {
	Type     FilterType
	Multiple bool
	Name     string
	Value    any
}

func CreateDynamicPropertiesFilter(schema *jsonschema.Schema, input map[string][]string) ([]DynamicFilter, error) {
	out := make([]DynamicFilter, 0)

	for k, v := range input {
		prop, exists := schema.Properties[k]
		if !exists {
			continue
		}

		switch prop.Types[0] {
		case "boolean":
			arr := make([]bool, len(v))
			for i, x := range v {
				b, err := parseBool(x)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", k, err)
				}

				arr[i] = b
			}

			filter := DynamicFilter{
				Type:     FilterTypeBool,
				Multiple: true,
				Name:     k,
				Value:    arr,
			}
			out = append(out, filter)
		case "string":
			filter := DynamicFilter{
				Type:     FilterTypeString,
				Multiple: true,
				Name:     k,
				Value:    v,
			}
			out = append(out, filter)
		default:
			filter := DynamicFilter{
				Type:     FilterTypeString,
				Multiple: true,
				Name:     k,
				Value:    v,
			}
			out = append(out, filter)
		}
	}

	return out, nil
}

func parseBool(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean: %q", s)
	}
}
