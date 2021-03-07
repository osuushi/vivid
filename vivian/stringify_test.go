package vivian

import "testing"

func TestStringify(t *testing.T) {
	check := func(input string, expectedOptional ...string) {
		expected := input
		if len(expectedOptional) == 1 {
			expected = expectedOptional[0]
		}

		ast, err := ParseString(input)
		if err != nil {
			t.Errorf("Unexpected error for %q: %v", input, err)
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

	check("foo")
	check("@red[foo]")
	check("@${} ")
	check("My email is ada@@example.com")
	check("@$()This apple costs $$5.00 (if you can believe it$)")
	check("@red[@-foo] is the value of @green[foo]")
	check("Let's see how @-this.chomp~ worked")
	check("Let's see how @-this.chomp@red[worked]", "Let's see how @-this.chomp~ @red[worked]")
	check("Let's see how @-this.chomp~ @red[worked]")
	check("Let's see how @-this.chomp[worked", "Let's see how @-this.chomp~ [worked")
}
