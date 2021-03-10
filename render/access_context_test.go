package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-yaml/yaml"
)

const example = `
- foo:
		bar: 2
	baz:
		- 1
		- 3
- qux: hello
	gribble:
		- gralt: 42
`

func parseYaml(data string) interface{} {
	var obj interface{}
	// Fraught to put space indented strings in a go file, so just convert the
	// tabs here.
	data = strings.ReplaceAll(data, "\t", "  ")
	err := yaml.Unmarshal([]byte(data), &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func TestAccessContext(t *testing.T) {
	check := func(data interface{}, path string, expected string) {
		actual := accessContext(data, path)
		if actual != expected {
			t.Errorf("For %q, expected %q but got %q", path, expected, actual)
		}
	}

	input := parseYaml(example)
	fmt.Println("Using example:", example)
	check(input, "0.foo.bar", "2")
	check(input, "0.foo.bar.bad", "")
	check(input, "0.baz.1", "3")
	check(input, "0.baz.length", "2")
	check(input, "length", "2")
	check(input, "1.gribble.0.gralt", "42")
}
