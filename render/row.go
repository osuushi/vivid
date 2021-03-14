package render

import (
	"strings"

	"github.com/osuushi/vivid/rich"
	"github.com/osuushi/vivid/vivian"
	"github.com/thomaso-mirodin/intmath/intgr"
)

// Public API to the package. You create a Row from a string, and that row can
// be fed a context to render it.

type Row struct {
	Cells []*Cell
}

func MakeRow(input string) (*Row, error) {
	ast, err := vivian.ParseString(input)
	if err != nil {
		return nil, err
	}
	return makeRowFromAst(ast)
}

func makeRowFromAst(ast *vivian.Ast) (*Row, error) {
	// We must have a "true" vivian tree with hoisted cell creators before we can
	// produce cells.
	err := hoistCells(ast)
	if err != nil {
		return nil, err
	}
	row := &Row{}
	row.Cells, err = cellsFromAst(ast)
	return row, err
}

// Context must be a generic structure, either []interface{} or
// map[interface{}]interface{}, and all values must be either primitives, or more
// of the same.
//
// Returns an array of lines, which are ANSI styled strings
func (row *Row) Render(width int, beam StyleBeam, context interface{}) ([]string, error) {
	sizedCells := AllocateCellSizes(row.Cells, width)
	cellLines := make([][]rich.RichString, len(sizedCells))

	for i, sc := range sizedCells {
		err := layoutAndStyleCell(sc, i > 0, context, &cellLines[i])
		if err != nil {
			return nil, err
		}
	}

	// Pad all cells to tallest
	maxHeight := 0
	for _, lines := range cellLines {
		maxHeight = intgr.Max(maxHeight, len(lines))
	}

	for i, lines := range cellLines {
		spacerCount := maxHeight - len(lines)
		if spacerCount == 0 {
			continue
		}

		sc := sizedCells[i]
		spacerWidth := sc.Width
		if i > 0 {
			spacerWidth++ // account for left pad
		}
		spacer := rich.MakeSpacer(spacerWidth, &rich.Style{
			Background: sc.Cell.Background,
		})

		for j := 0; j < spacerCount; j++ {
			cellLines[i] = append(cellLines[i], spacer)
		}
	}

	result := make([]string, maxHeight)

	// Scan each line
	for i := 0; i < maxHeight; i++ {
		var builder strings.Builder

		// Get the current line from each cell line group
		for _, lines := range cellLines {
			line := lines[i]
			for _, r := range line {
				beam.ScanRune(r, &builder)
			}
		}
		beam.Terminate(&builder)
		result[i] = builder.String()
	}

	return result, nil
}

func layoutAndStyleCell(
	sc *SizedCell,
	leftPad bool,
	context interface{},
	dest *[]rich.RichString,
) error {
	cellStyle := &rich.Style{
		Background: sc.Cell.Background,
	}
	content, err := stylizeNodes(sc.Cell.Content, context, cellStyle)
	if err != nil {
		return err
	}

	lines := renderCell(content, sc)

	if leftPad {
		for j, line := range lines {
			lines[j] = rich.Concat(
				rich.NewRichString(" ", cellStyle),
				line,
			)
		}
	}
	*dest = lines
	return nil
}
