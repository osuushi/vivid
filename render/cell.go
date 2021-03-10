package render

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/thomaso-mirodin/intmath/intgr"

	"github.com/osuushi/vivid/vivian"
)

type Cell struct {
	MinWidth int
	MaxWidth int
	Wrap     bool

	Greedy bool

	Shyness int
	Glue    bool

	Alignment Alignment
	Content   []vivian.Node
}

type Alignment byte

const (
	Left Alignment = iota
	Right
	Center
	Justify
)

// Split an AST into cells. The AST must already be hoisted so that cell
// creators live at the top level
func cellsFromAst(hoistedAst *vivian.Ast) ([]*Cell, error) {
	cells := []*Cell{}
	currentCellNodes := []vivian.Node{}

	// When we hit the end of a run of non-cell creators, they get bundled up into
	// a default-configured cell.
	aggregateTopLevelCellIfNeeded := func() {
		if len(currentCellNodes) == 0 {
			return
		}
		cell := makeDefaultCell()
		cell.Content = currentCellNodes
		cells = append(cells, cell)
		currentCellNodes = []vivian.Node{}
	}

	for _, node := range hoistedAst.Content.Children {
		if isCellCreatorNode(node) {
			aggregateTopLevelCellIfNeeded()
			cell, err := cellFromCellCreator(node)
			if err != nil {
				return nil, err
			}
			cells = append(cells, cell)
		} else {
			currentCellNodes = append(currentCellNodes, node)
		}
	}
	aggregateTopLevelCellIfNeeded()

	return cells, nil
}

func makeDefaultCell() *Cell {
	return &Cell{
		MaxWidth: math.MaxInt16,
	}
}

func cellFromCellCreator(node vivian.Node) (*Cell, error) {
	cell := makeDefaultCell()
	var contentNode *vivian.ContentNode
	for isCellCreatorNode(node) {
		contentNode = node.(*vivian.ContentNode)
		err := applyTag(contentNode.Tag, cell)
		if err != nil {
			return nil, err
		}
		node = contentNode.Children[0]
	}
	cell.Content = contentNode.Children
	return cell, nil
}

var tagParsePattern *regexp.Regexp

func init() {
	tagParsePattern = regexp.MustCompile("(?i)([a-z]+)(\\d*)")
}

func applyTag(tag string, cell *Cell) error {
	match := tagParsePattern.FindStringSubmatch(tag)
	name := match[1]
	// Only invalid parse is empty string, and it's fine for that to be zero
	param64, _ := strconv.ParseInt(match[2], 10, 64)
	param := int(param64)

	switch name {
	case "min":
		cell.MinWidth = param
	case "max":
		cell.MaxWidth = param
	case "fixed":
		cell.MinWidth = param
		cell.MaxWidth = param
	case "wrap":
		cell.Wrap = true
		// Number parameter is shorthand for @wrap@fixedNN
		cell.MinWidth = param
		cell.MaxWidth = param
	case "strut":
		cell.Greedy = true
	case "shy":
		cell.Shyness = intgr.Min(param, 1)
	case "glue":
		cell.Glue = true
	case "left":
		// Pointless assignment, but for consistency
		cell.Alignment = Left
	case "right":
		cell.Alignment = Right
	case "center":
		cell.Alignment = Center
	case "justify":
		cell.Alignment = Justify
	default:
		// Only possible to escape other cases via a bug. No other tag would have
		// matched as a cell creator.
		return fmt.Errorf("Unhandled cell creator type %s", name)
	}
	return nil
}
