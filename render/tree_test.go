package render

import (
	"strings"
	"testing"

	"github.com/osuushi/vivid/vivian"
)

func TestHoistCells(t *testing.T) {
	check := func(input string, expectedOptional ...string) {
		expected := input
		if len(expectedOptional) == 1 {
			expected = expectedOptional[0]
		}

		ast, err := vivian.ParseString(input)
		if err != nil {
			t.Errorf("Unexpected parse error for %q: %v", input, err)
			return
		}
		err = hoistCells(ast)
		if err != nil {
			t.Errorf("Unexpected hoist error for %q: %v", input, err)
			return
		}

		actual := ast.String()
		if actual != expected {
			t.Errorf("For input: %q\n"+
				"Expected: %q\n"+
				"Got: %q",
				input, expected, actual,
			)
		}
	}

	checkErr := func(input string, expected string) {
		ast, err := vivian.ParseString(input)
		if err != nil {
			t.Errorf("Unexpected parse error for %q: %v", input, err)
			return
		}
		err = hoistCells(ast)
		if err != nil {
			if !strings.Contains(err.Error(), expected) {
				t.Errorf(
					"Expected hoist error for %q: to contain %q\nGot: %v",
					input, expected, err,
				)
				return
			}
		} else {
			actual := ast.String()
			t.Errorf("For input: %q\n"+
				"Expected error with: %q\n"+
				"Got successful hoist: %q",
				input, expected, actual,
			)
		}
	}

	check("foo")
	check("@red[foo]")
	check("@wrap[foo] @wrap[bar]")
	check("@red[@wrap[foo]]", "@wrap[@red[foo]]")
	check(
		"@red[@wrap[foo] @wrap[bar]]",
		"@wrap[@red[foo]]@red[ ]@wrap[@red[bar]]",
	)

	check(
		"@red[@wrap[@bold[foo]] @wrap[@blue[bar]]]",
		"@wrap[@red[@bold[foo]]]@red[ ]@wrap[@red[@blue[bar]]]",
	)
	check(
		"@red@wrap[@-hello] @bold@max30[@black[darkness] my old friend]",
		"@wrap[@red[@-hello]] @max30[@bold[@black[darkness] my old friend]]",
	)

	// Errors
	checkErr("@wrap[@max30[hello] @auto[world]]", "subdivided by \"max30\"")
	checkErr("@wrap[@max30[hello] @red[world]]", "subdivided by \"max30\"")
	checkErr("@wrap[@max30[hello] world]", "subdivided by \"max30\"")
	checkErr("@wrap[@max30[@auto[hello]] world]", "subdivided by \"max30\"")
	checkErr("@wrap[@max30[@auto[hello] world]]", "subdivided by \"auto\"")
}
