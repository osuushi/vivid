package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/osuushi/vivid/vivian"
)

func makeCells(s string) []*Cell {
	ast, err := vivian.ParseString(s)
	if err != nil {
		panic(err)
	}

	hoistCells(ast)
	cells, err := cellsFromAst(ast)
	if err != nil {
		panic(err)
	}
	return cells
}

// Helper to get cell text. Assumes the cell ultimately just contains a text
// node
func getInnerCellText(cell *Cell) string {
	node := cell.Content[0]
	for {
		textNode, ok := node.(*vivian.TextNode)
		if ok {
			return textNode.Text
		}
		contentNode, ok := node.(*vivian.ContentNode)
		if !ok {
			panic(fmt.Sprintf("Unexpected node type: %#v", node))
		}
		node = contentNode.Children[0]
	}
}

func getAllCellText(cells []*SizedCell) string {
	texts := []string{}
	for _, cell := range cells {
		texts = append(texts, getInnerCellText(cell.Cell))
	}
	return strings.Join(texts, " ")
}

func TestTrimShyCells(t *testing.T) {
	check := func(
		input string,
		width int,
		expected string,
	) {
		cells := makeCells(input)

		list := makeSizedCellList(cells)
		list.applyMinimumSizes(width)
		list.trimShyCells(width)
		sized := list.toSlice()

		actual := getAllCellText(sized)
		if actual != expected {
			t.Errorf(
				"For %q at width %d,\nexpected %q but got %q",
				input, width, expected, actual,
			)
		}
	}

	check("@fixed30[foo]@fixed30@shy[bar]@fixed30@shy[baz]", 90, "foo bar baz")
	check("@fixed30[foo]@fixed30@shy[bar]@fixed30@shy[baz]", 70, "foo bar")
	check("@fixed30[foo]@fixed30@shy2[bar]@fixed30@shy[baz]", 70, "foo baz")
	check("@fixed30[foo]@fixed30@shy[bar]@fixed30[baz]", 70, "foo baz")
	check("@fixed30[foo]@fixed30@shy[bar]@fixed30@glue[baz]", 70, "foo")
	check("@fixed30[foo] @fixed30@shy[bar] @fixed30@glue[baz]", 70, "foo")
}
