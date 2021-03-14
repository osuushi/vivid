package render

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/kr/pretty"
	"github.com/osuushi/vivid/rich"
)

// Test beam just handles basic text style by putting symbols before each
// styled character (which doesn't count against the width)

type TestBeam struct{}

func (t *TestBeam) ScanRune(r rich.RichRune, b *strings.Builder) {
	if r.IsItalic() {
		b.WriteRune('𝐼')
	}
	if r.IsBold() {
		b.WriteRune('𝗕')
	}
	if r.IsUnderline() {
		b.WriteRune('⎁')
	}
	b.WriteRune(r.Rune)
}

func (t *TestBeam) Terminate(*strings.Builder) {}

func TestRender(t *testing.T) {
	context := map[interface{}]interface{}{
		"name": "Ada",
		"columns": []interface{}{
			"This is a test column.\n\nIt has multiple paragraphs and wraps repeatedly",
			"This column only wraps but has no paragraphs.",
			"This line is short",
		},
	}
	check := func(input string, expected ...string) {
		row, err := MakeRow(input)
		if err != nil {
			t.Fatal(err)
		}

		result, err := row.Render(
			50,
			&TestBeam{},
			context,
		)

		if err != nil {
			t.Fatal(err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			pretty.Println(result)
			t.Errorf("On input: %q\n%s", input, strings.Join(diff, "\n"))
		}
	}

	check(
		"My name is @-name",
		"My name is Ada                                    ",
	)

	check(
		"@right[My name is @-name]",
		"                                    My name is Ada",
	)

	check(
		"@right@fixed10[Name:] @bold[@-name]",
		"     Name: 𝗕A𝗕d𝗕a                                    ",
	)

	check(
		"@right@fixed10[Name:] @min5@center[@b@i[@-name]]",
		"     Name:                   𝐼𝗕A𝐼𝗕d𝐼𝗕a                  ",
	)

	check(
		"@wrap@justify@fixed20[@-columns.0] @wrap@fixed15[@-columns.1] @bold[@-columns.2]",
		"This   is   a   test This column     𝗕T𝗕h𝗕i𝗕s𝗕 𝗕l𝗕i𝗕n𝗕e𝗕 𝗕i𝗕… ",
		"column.              only wraps but               ",
		"                     has no                       ",
		"It    has   multiple paragraphs.                  ",
		"paragraphs and wraps                              ",
		"repeatedly                                        ",
	)

	// Using a strut to give paragraphs room to breathe

	check(
		"@wrap@justify@fixed26[@-columns.0] @strut[] @wrap@fixed20[@-columns.1]",
		"This is a test column.        This column only    ",
		"                              wraps but has no    ",
		"It has multiple paragraphs    paragraphs.         ",
		"and wraps repeatedly                              ",
	)

	// Too small for second column

	check(
		"@wrap@justify@fixed40[@-columns.0] @strut[] @wrap@fixed20[@-columns.1]",
		"This is a test column.                            ",
		"                                                  ",
		"It  has  multiple paragraphs  and  wraps          ",
		"repeatedly                                        ",
	)

	// Too small for second column, first column is shy
	check(
		"@wrap@shy@justify@fixed40[@-columns.0] @strut[] @wrap@fixed20[@-columns.1]",
		"                              This column only    ",
		"                              wraps but has no    ",
		"                              paragraphs.         ",
	)

	// Now the strut is glued, so it gets deleted
	check(
		"@wrap@shy@justify@fixed40[@-columns.0] @strut@glue[] @wrap@fixed20[@-columns.1]",
		"This column only    ",
		"wraps but has no    ",
		"paragraphs.         ",
	)

	// Now the strut is glued, so it gets deleted
	check(
		"@wrap@shy@justify@fixed40[@-columns.0] @strut@glue[] @wrap@min20[@-columns.1]",
		"This column only wraps but has no paragraphs.     ",
	)
}
