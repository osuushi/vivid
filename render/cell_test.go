package render

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/osuushi/vivid/vivian"
)

func TestApplyTagAllCellCreatorsHandled(t *testing.T) {
	// Ensure that no case panics in applyTag
	for _, prefix := range cellCreatingNamePrefixes {
		cell := makeDefaultCell()
		err := applyTag(prefix, cell)
		if err != nil && strings.HasPrefix(err.Error(), "Unhandled cell creator type") {
			t.Error(err.Error())
		}
	}
}

func TestCellsFromAst(t *testing.T) {
	ast, err := vivian.ParseString("hello @-val @max30[this is] @fixed50@green@right[a test]")
	hoistCells(ast)

	cells, err := cellsFromAst(ast)
	if err != nil {
		t.Error(err)
	}

	// Clear content nodes, since we don't want to test those after this point
	for _, cell := range cells {
		cell.Content = nil
	}

	expected := []*Cell{
		makeDefaultCell(),
		&Cell{
			MaxWidth: 30,
		},
		// Note that space was culled because it was in its own implicit cell
		&Cell{
			MinWidth:  50,
			MaxWidth:  50,
			Alignment: Right,
		},
	}

	if diff := deep.Equal(cells, expected); diff != nil {
		t.Error("\n" + strings.Join(diff, "\n"))
	}
}
