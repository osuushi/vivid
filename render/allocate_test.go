package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
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
func getCellText(sizedCell *SizedCell) string {
	cell := sizedCell.Cell
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
		texts = append(texts, getCellText(cell))
	}
	return strings.Join(texts, " ")
}

func TestShyestNode(t *testing.T) {
	check := func(input, expected string) {
		cells := makeCells(input)
		list := makeSizedCellList(cells)
		actual := getCellText(list.shyestNode().val)
		if actual != expected {
			t.Errorf(
				"shyestNode() for %q\nExpected: %q\nGot:%q",
				input, expected, actual,
			)
		}
	}

	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", "baz")
	check("@fixed30@shy[foo] @fixed30[bar] @fixed30[baz]", "foo")
	check("@fixed30@shy[foo] @fixed30@shy[bar] @fixed30[baz]", "bar")
	check("@fixed30@shy[foo] @fixed30@shy[bar] @fixed30@shy[baz]", "baz")
	check("@fixed30@shy2[foo] @fixed30@shy[bar] @fixed30@shy[baz]", "foo")
	check("@fixed30@shy2[foo] @fixed30@shy[bar] @fixed30@shy2[baz]", "baz")
	check("@fixed30@shy4[foo] @fixed30@shy2[bar] @fixed30@shy3[baz]", "foo")
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
		list = list.trimShyCells(width)
		sized := list.toSlice()

		actual := getAllCellText(sized)
		if actual != expected {
			t.Errorf(
				"For %q at width %d,\nexpected %q but got %q",
				input, width, expected, actual,
			)
		}
	}

	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", 92, "foo bar baz")
	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", 70, "foo bar")
	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", 61, "foo bar")
	// not enough room with padding
	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", 60, "foo")
	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", 30, "foo")
	check("@fixed30[foo] @fixed30[bar] @fixed30[baz]", 20, "")
	check("@fixed30@shy[foo] @fixed30[bar] @fixed30[baz]", 30, "bar")
	check("@fixed30[foo] @fixed30@shy[bar] @fixed30@shy[baz]", 92, "foo bar baz")
	check("@fixed30[foo] @fixed30@shy[bar] @fixed30@shy[baz]", 70, "foo bar")
	check("@fixed30[foo] @fixed30@shy2[bar] @fixed30@shy[baz]", 70, "foo baz")
	check("@fixed30[foo] @fixed30@shy[bar] @fixed30[baz]", 70, "foo baz")
	check("@fixed30[foo] @fixed30@shy[bar] @fixed30@glue[baz]", 70, "foo")
}

func TestExpandCells(t *testing.T) {
	check := func(input string, twoPass bool, expectedWidths ...int) {
		width := 150
		cells := makeCells(input)
		list := makeSizedCellList(cells)
		list.applyMinimumSizes(width)
		if twoPass {
			list.expandCells(width, true)
		}
		list.expandCells(width, !twoPass)

		actualWidths := []int{}
		list.each(func(node *scNode) {
			actualWidths = append(actualWidths, node.val.Width)
		})
		if diff := deep.Equal(actualWidths, expectedWidths); diff != nil {
			t.Errorf("For input: %q\n%s", input, strings.Join(diff, "\n"))
		}
	}

	// One cell by its lonesome
	check(
		"@min30[]",
		true, 150,
	)

	// Remember that there is one space for padding between cells
	check(
		"@fixed30[] @min40@strut[]",
		false, 30, 119,
	)
	// After greedy cell has its turn, there's nothing left to allocate
	check(
		"@min30[] @min40@strut[]",
		true, 30, 119,
	)

	// No greed, so both cells take the free space evenly. Note that they don't
	// end up at even sizes because their minimums are different
	check(
		"@min30[] @min40@max90[]",
		true, 70, 79,
	)

	// The second cell is greedy, so it's going to eat its fill of free space in
	// this first pass.
	check(
		"@min30[] @min40@max90@strut[]",
		false, 30, 90,
	)

	// On the non-greedy pass, the first cell takes the rest
	check(
		"@min30[] @min40@max90@strut[]",
		true, 59, 90,
	)

	// In this case, the nongreedy cells divide free space amongst themselves
	// before the nongreedy cell in the middle has a chance
	check(
		"@max60@strut[] @min40[] @max70@strut[]",
		false, 54, 40, 54,
	)

	// Sorry buddy, but the greedy cells took all the free space.
	check(
		"@max70@strut[] @min40[] @max70@strut[]",
		true, 54, 40, 54,
	)

	// Similar to the above case, but the first greedy cell has a head start now
	// because of its min
	check(
		"@min20@max70@strut[] @min40[] @max70@strut[]",
		true, 64, 40, 44,
	)

	// A similar variation, but this time the first cell's ravenous appetite
	// causes it to reach its capacity, leaving more for the other greedy cell.
	check(
		"@min20@max60@strut[] @min40[] @max70@strut[]",
		true, 60, 40, 48,
	)

	// We didn't use up all the space because both cells were fixed. That's fine.
	check(
		"@fixed30[] @fixed30[]",
		true, 30, 30,
	)
}
