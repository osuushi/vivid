package render

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"unicode"

	"github.com/thomaso-mirodin/intmath/intgr"

	"github.com/osuushi/vivid/rich"
	"github.com/osuushi/vivid/vivian"
)

type Cell struct {
	MinWidth int
	MaxWidth int
	Wrap     bool

	Greedy bool

	Shyness int
	Glue    bool

	Background *rich.RGB
	Alignment  Alignment
	Content    []vivian.Node
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

		// If a cell was created implicitly and contains no non-space or input, it
		// is omitted.
		if isEmptyImplicitCell(currentCellNodes) {
			currentCellNodes = []vivian.Node{}
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
		MinWidth: 1,
		// Out of laziness, using MaxInt16 to avoid having to worry about overflows.
		// This is far wider than any cell would reasonably be.
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
		if len(contentNode.Children) == 0 {
			node = nil
		} else {
			node = contentNode.Children[0]
		}
	}
	cell.Content = contentNode.Children
	return cell, nil
}

func isEmptyImplicitCell(nodes []vivian.Node) bool {
	for _, node := range nodes {
		if !isEmptyNode(node) {
			return false
		}
	}
	return true
}

// Is the node "empty", meaning it contains nothing but whitespace, regardless
// of style?
func isEmptyNode(node vivian.Node) bool {
	switch node := node.(type) {
	case *vivian.ContentNode: // empty if all children are empty
		for _, child := range node.Children {
			if !isEmptyNode(child) {
				return false
			}
			return true
		}
	case *vivian.InputNode: // An input node is never empty
		return false
	case *vivian.TextNode: // Empty if nothing but whitespace
		return isStringWhitespace(node.Text)
	default:
		panic(fmt.Sprintf("Unknown node type: %v", node))
	}
	return false
}

func isStringWhitespace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

var tagParsePattern *regexp.Regexp

func init() {
	tagParsePattern = regexp.MustCompile("(?i)([a-z]+)(\\d*)")
}

func applyTag(tag string, cell *Cell) error {
	// Special case is bg color, which does not follow the typical <name><number>
	// pattern
	color, isColor := parseBgColor(tag)
	if isColor {
		cell.Background = color
		return nil
	}

	match := tagParsePattern.FindStringSubmatch(tag)
	name := match[1]
	// Only invalid parse is empty string, and it's fine for that to be zero
	param64, _ := strconv.ParseInt(match[2], 10, 64)
	param := int(param64)

	switch name {
	case "min":
		cell.MinWidth = intgr.Max(param, 1)
	case "max":
		cell.MaxWidth = param
	case "fixed":
		cell.MinWidth = param
		cell.MaxWidth = param
	case "wrap":
		cell.Wrap = true
		// Number parameter is shorthand for @wrap@fixedNN
		if param != 0 {
			cell.MinWidth = param
			cell.MaxWidth = param
		}
	case "strut":
		cell.Greedy = true
	case "shy":
		cell.Shyness = intgr.Max(param, 1)
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
