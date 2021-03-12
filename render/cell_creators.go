package render

import (
	"strings"
	"unicode"

	"github.com/osuushi/vivid/vivian"
)

// Tag names that cause a cell to be created, and therefore must be moved to the
// top level.
var cellCreatingNamePrefixes = []string{
	// ** Types with parameter
	"min",   // minimum width
	"max",   // maximum width
	"fixed", // same min and max
	"wrap",  // wrap text; parameter is same as @wrap@fixedNN

	// Override priority for hiding this cell if there isn't enough space for
	// everyone's minimum. If no parameter is given, it defaults to 1 (normal
	// shyness is 0). Ties are broken from right to left.
	"shy",

	// Do not display the previous cell if this cell is hidden, and vice versa.
	// These can be chained across three or more cells as well, and the group will
	// effectively take maximum shyness out of the group.
	"glue",

	// ** Types without parameter

	// greedy cell; greedy cells eat up all free space up to their max before
	// nongreedy cells have a chance to take on space above their min
	"strut",

	// Alignments; these are cell creators because it is meaningless to be aligned
	// outside the context of a boundary.
	"left", // this is a noop, but kept for consistency
	"right",
	"center",
	"justify",
}

func isNumeric(str string) bool {
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isCellCreator(name string) bool {
	if isBgColor(name) {
		return true
	}

	for _, prefix := range cellCreatingNamePrefixes {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		suffix := strings.TrimPrefix(name, prefix)
		if isNumeric(suffix) {
			return true
		}
	}
	return false
}

func isCellCreatorNode(node vivian.Node) bool {
	contentNode, ok := node.(*vivian.ContentNode)
	if !ok { // Non-content node never a cell creator
		return false
	}
	return isCellCreator(contentNode.Tag)
}
