package render

import (
	"testing"

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
