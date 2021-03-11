package render

import (
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/osuushi/vivid/rich"
	"github.com/osuushi/vivid/vivian"
)

// Assert that predicate matches wherever there's an X and not where there isn't
func assertStyleMatch(
	t *testing.T,
	input string,
	context interface{},
	expectedText string,
	mask string,
	predicateName string,
	predicate func(rich.RichRune) bool,
) {
	ast, err := vivian.ParseString(input)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := stylizeNodes(
		ast.Content.Children,
		context,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	for i, r := range rs {
		matches := predicate(r)
		shouldMatch := mask[i] == 'X'

		if rs.String() != expectedText {
			t.Errorf("%q:\nExpected: %q, got %q", input, expectedText, rs.String())
			return
		}

		if matches != shouldMatch {
			pad := strings.Repeat(" ", i+1) // +1 for quote
			infinitive := "to"
			if !shouldMatch {
				infinitive = "not to"
			}
			t.Errorf(
				"%q\nExpected %q predicate %s match at:\n%q\n%s\nBad rune style:\n%s",
				input,
				predicateName,
				infinitive,
				rs.String(),
				pad+"^",
				pretty.Sprint(r.GetStyle()),
			)

			return
		}
	}
}

func TestStylizeNodes(t *testing.T) {
	assertStyleMatch(
		t,
		"this is @red[a] @bold@green[@-noun]",
		map[interface{}]interface{}{"noun": "test"},
		"this is a test",
		"        X     ",
		"red",
		func(r rich.RichRune) bool {
			color := r.GetColor()
			if color == nil {
				return false
			}
			return *color == rich.RGB{R: 0xff, G: 0x00, B: 0x00}
		},
	)

	assertStyleMatch(
		t,
		"this is @red[a] @bold@green[@-noun]",
		map[interface{}]interface{}{"noun": "test"},
		"this is a test",
		"          XXXX",
		"green",
		func(r rich.RichRune) bool {
			color := r.GetColor()
			if color == nil {
				return false
			}
			return *color == rich.RGB{R: 0x00, G: 0x80, B: 0x00}
		},
	)

	assertStyleMatch(
		t,
		"this is @red[a] @bold@green[@-noun]",
		map[interface{}]interface{}{"noun": "test"},
		"this is a test",
		"          XXXX",
		"bold",
		func(r rich.RichRune) bool {
			return r.IsBold()
		},
	)

	assertStyleMatch(
		t,
		"this is @fff@bg345678[hex colors]",
		nil,
		"this is hex colors",
		"        XXXXXXXXXX",
		"xfff foreground x345678 background",
		func(r rich.RichRune) bool {
			color := r.GetColor()
			if color == nil {
				return false
			}
			if *color != (rich.RGB{R: 0xff, G: 0xff, B: 0xff}) {
				return false
			}

			color = r.GetBackground()
			if color == nil {
				return false
			}
			if *color != (rich.RGB{R: 0x34, G: 0x56, B: 0x78}) {
				return false
			}
			return true
		},
	)
}
