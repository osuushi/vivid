package vivian

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/kr/pretty"
)

func TestParseString(t *testing.T) {
	check := func(expr string, wantErr bool, expectedChildren ...interface{}) {
		actualAst, err := ParseString(expr)
		if err != nil {
			if !wantErr {
				t.Errorf("Unexpected error for %q: %v", expr, err)
			}
		} else if wantErr {
			t.Errorf("Expected %q to error", expr)
		} else if diff := deep.Equal(actualAst.Content.Children, expectedChildren); diff != nil {
			t.Errorf("Wrong parse for %q: %v\n Got: %s", expr, diff, pretty.Sprint(actualAst.Content.Children))
		}
	}

	// Basic expressions
	check("foo", false, &TextNode{Text: "foo"})

	check("foo@@@][", false, &TextNode{Text: "foo@]["})

	check("@foo[]", false, &ContentNode{
		Tags:     []string{"foo"},
		Children: nil,
	})

	check("@foo@bar[]", false, &ContentNode{
		Tags:     []string{"foo", "bar"},
		Children: nil,
	})

	check("@-foo.bar.baz", false, &InputNode{
		Path: []string{"foo", "bar", "baz"},
	})

	// Swapped tokens
	check("@{}@foo{}", false, &ContentNode{
		Tags:     []string{"foo"},
		Children: nil,
	})

	check("@#><#foo><", false, &ContentNode{
		Tags:     []string{"foo"},
		Children: nil,
	})

	// Escape hell
	check("@!()foo@!!!)([]", false, &TextNode{Text: "foo@!)([]"})

	// Multi-expressions
	check("This is an @-adjective @-noun.", false,
		&TextNode{Text: "This is an "},
		&InputNode{
			Path: []string{"adjective"},
		},
		&TextNode{Text: " "},
		&InputNode{
			Path: []string{"noun"},
		},
	)

	// Chomp character
	check("This is the @-smokingPrefix~ smoking section.", false,
		&TextNode{Text: "This is the "},
		&InputNode{
			Path: []string{"smokingPrefix"},
		},
		&TextNode{Text: "smoking section."},
	)

	// Error cases
	check("@", true)
	check("@[", true)
	check("@[[", true)
	check("@ab", true)
	check("@$[", true)
	check("@$[[", true)
	check("@$ab", true)
	check("]", true)
	check("foo]", true)
	check("foo @[", true)
	check("This is a @test", true)
	check("This has @- but no input path", true)
}
