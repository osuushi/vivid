package vivian

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/kr/pretty"
)

func TestParseString(t *testing.T) {
	check := func(expr string, wantErr bool, expectedChildren ...Node) {
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
		Tag:      "foo",
		Children: []Node{},
	})

	check("@foo@bar[]", false, &ContentNode{
		Tag: "foo",
		Children: []Node{
			&ContentNode{
				Tag:      "bar",
				Children: []Node{},
			},
		},
	})

	check("@-foo.bar.baz", false, &InputNode{
		Path: []string{"foo", "bar", "baz"},
	})

	// Swapped tokens
	check("@{}@foo{}", false, &ContentNode{
		Tag:      "foo",
		Children: []Node{},
	})

	check("@#><#foo><", false, &ContentNode{
		Tag:      "foo",
		Children: []Node{},
	})

	// Escape hell
	check("@!()foo@!!!)([]", false, &TextNode{Text: "foo@!)([]"})

	// Multi-expressions
	check("This is an @-adjective @-noun~ .", false,
		&TextNode{Text: "This is an "},
		&InputNode{
			Path: []string{"adjective"},
		},
		&TextNode{Text: " "},
		&InputNode{
			Path: []string{"noun"},
		},
		&TextNode{Text: "."},
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
	check("My email is ada@example.com", true)
	check("This is a @test", true)
	check("This has @- but no input path", true)
	check("This @-path.ends.in.a.dot.", true)
	check("This @-path..has..two.dots", true)
}
