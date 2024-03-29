package render

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/kr/pretty"
	"github.com/osuushi/vivid/rich"
)

func TestJustifyLine(t *testing.T) {
	var input rich.RichString
	input = rich.NewRichString("Whan that aprill with his shoures soote", nil)
	check := func(width int, expected string) {
		actual := justifyLine(input, width).String()
		if actual != expected {
			t.Errorf("Whan that…: for width %d, Expected \n%q\nGot:\n%q", width, expected, actual)
		}
	}

	check(50, "Whan   that   aprill   with  his   shoures   soote")
	check(49, "Whan   that  aprill   with   his  shoures   soote")
	check(40, "Whan that aprill  with his shoures soote")
	check(39, "Whan that aprill with his shoures soote")
}

func TestAlignRow(t *testing.T) {
	var input rich.RichString
	input = rich.NewRichString("Whan that aprill with his shoures soote", nil)
	check := func(width int, alignment Alignment, alignName string, expected string) {
		actual := alignRow(input, alignment, width).String()
		if actual != expected {
			t.Errorf("Whan that… %s: for width %d, Expected \n%q\nGot:\n%q", alignName, width, expected, actual)
		}
	}

	check(50, Left, "left", "Whan that aprill with his shoures soote           ")
	check(50, Right, "right", "           Whan that aprill with his shoures soote")
	check(50, Center, "center", "     Whan that aprill with his shoures soote      ")
	check(49, Center, "center", "     Whan that aprill with his shoures soote     ")
	check(50, Justify, "justify", "Whan   that   aprill   with  his   shoures   soote")

	// Noops
	check(39, Left, "left", "Whan that aprill with his shoures soote")
	check(39, Right, "right", "Whan that aprill with his shoures soote")
	check(39, Center, "center", "Whan that aprill with his shoures soote")
	check(39, Justify, "justify", "Whan that aprill with his shoures soote")
}

var testParagraph = "This   is\ta test of paragraph slicing. Here's a long word: \"antidisestablishmentarianism\". Wow. What a mouthful."

func TestSliceParagraph(t *testing.T) {
	check := func(width int, expected []string) {
		actual := sliceParagraph(rich.NewRichString(testParagraph, nil), width)

		actualStrings := []string{}
		for _, row := range actual {
			actualStrings = append(actualStrings, row.String())
		}

		if diff := deep.Equal(actualStrings, expected); diff != nil {
			pretty.Println(actualStrings)
			t.Errorf("For width %d\n%s", width, strings.Join(diff, "\n"))
		}
	}

	check(27, []string{
		"This is a test of paragraph",
		"slicing. Here's a long",
		"word:",
		"\"antidisestablishmentariani",
		"sm\". Wow. What a mouthful.",
	})

	check(30, []string{"This is a test of paragraph",
		"slicing. Here's a long word:",
		"\"antidisestablishmentarianism\"",
		". Wow. What a mouthful.",
	})

	check(40, []string{
		"This is a test of paragraph slicing.",
		"Here's a long word:",
		"\"antidisestablishmentarianism\". Wow.",
		"What a mouthful.",
	})

	check(70, []string{
		"This is a test of paragraph slicing. Here's a long word:",
		"\"antidisestablishmentarianism\". Wow. What a mouthful.",
	})
}

func TestScanNextLine(t *testing.T) {
	check := func(actual rich.RichString, expected string) {
		if expected == "" {
			if actual != nil {
				t.Errorf("Expected nil for last remainder, not %s", actual.String())
			}
		}

		actualString := actual.String()
		if actualString != expected {
			// These tests are backward-dependent, so we can't keep going with this
			// test after a failure.
			t.Fatalf("\nExpected: %q\nGot:%q", expected, actual)
		}
	}
	input := rich.NewRichString("It has multiple paragraphs and wraps repeatedly", nil)
	var line rich.RichString

	line, input = scanNextLine(input, 8)
	check(line, "It has")
	check(input, " multiple paragraphs and wraps repeatedly")

	line, input = scanNextLine(input, 8)
	check(line, "multiple")
	check(input, " paragraphs and wraps repeatedly")

	line, input = scanNextLine(input, 8)
	check(line, "paragrap")
	check(input, "hs and wraps repeatedly")

	line, input = scanNextLine(input, 8)
	check(line, "hs and")
	check(input, " wraps repeatedly")

	line, input = scanNextLine(input, 8)
	check(line, "wraps")
	check(input, " repeatedly")

	line, input = scanNextLine(input, 8)
	check(line, "repeated")
	check(input, "ly")

	line, input = scanNextLine(input, 8)
	check(line, "ly")
	check(input, "")

	// regression test
	input = rich.NewRichString("Fighting, Psychic", nil)
	line, input = scanNextLine(input, 16)
	check(line, "Fighting,")
	check(input, " Psychic")

	line, input = scanNextLine(input, 16)
	check(line, "Psychic")
	check(input, "")
}

func TestNormalizeWhitespace(t *testing.T) {
	check := func(input string, expected string) {
		actual := normalizeWhitespace(rich.NewRichString(input, nil)).String()
		if actual != expected {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	}

	check("Spaces are fine here", "Spaces are fine here")
	check("This   has multi  spaces", "This has multi spaces")
	check(" this has one leading space", "this has one leading space")
	check("   this has leading spaces", "this has leading spaces")
	check("this has one trailing space ", "this has one trailing space")
	check("this has trailing spaces   ", "this has trailing spaces")
	check("\t  this is \t   \t lousy \t \t with tabs", "this is lousy with tabs")
}

func TestTruncateContentToWidth(t *testing.T) {
	check := func(input string, width int, expected string) {
		actual := truncateContentToWidth(rich.NewRichString(input, nil), width).String()
		if actual != expected {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	}

	check("This is a test", 1, "…")
	check("This is a test", 2, "T…")
	check("This is a test", 7, "This i…")
	check("This is a test", 13, "This is a te…")
	check("This is a test", 14, "This is a test")
}

// Smoke test, putting it all together
func TestRenderContent(t *testing.T) {
	input := strings.TrimSpace(`
		This is a   test that   has multiple paragraphs.

		See? This is technically  the third  paragraph because of the empty line.
	`)

	actualRichStrings := renderContent(
		rich.NewRichString(input, nil),
		44,
		true,
		Left,
	)
	actual := []string{}
	for _, rs := range actualRichStrings {
		actual = append(actual, rs.String())
	}
	expected := []string{
		"This is a test that has multiple paragraphs.",
		"                                            ",
		"See? This is technically the third paragraph",
		"because of the empty line.                  ",
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Errorf("\n" + strings.Join(diff, "\n"))
	}
}
