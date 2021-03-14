package render

import (
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

// Access the data in an arbitrary map/slice structure by a path.
// Arrays are treated exactly as if they were objects with numeric strings for
// keys and a `length` property.
//
// Output is always cast to a string. If a value is missing or not
// stringifyable, an empty string is returned.

func accessContext(context interface{}, components []string) string {
	result := accessContextWithComponents(context, components)
	if array, ok := result.([]interface{}); ok {
		stringified := make([]string, len(array))
		for i, val := range array {
			stringified[i] = cast.ToString(val)
		}
		return strings.Join(stringified, ", ")
	}
	return cast.ToString(result)
}

func accessContextWithComponents(context interface{}, components []string) interface{} {
	// No components left
	if len(components) == 0 {
		return context
	}

	if context == nil {
		return nil
	}

	firstComponent := components[0]
	components = components[1:]

	switch context := context.(type) {
	case map[string]interface{}:
		return accessContextWithComponents(
			context[firstComponent],
			components,
		)
	case map[interface{}]interface{}:
		return accessContextWithComponents(
			context[firstComponent],
			components,
		)
	case []interface{}:
		if firstComponent == "length" {
			return accessContextWithComponents(len(context), components)
		} else {
			i, err := strconv.Atoi(firstComponent)
			// Can't access array or string with non-numeric index
			if err != nil {
				return nil
			}
			return accessContextWithComponents(context[i], components)
		}
	default: // Attempt to access anything off of a non-object non-array yields nil
		return nil
	}
}
