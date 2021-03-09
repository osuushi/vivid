package render

import (
	"strings"
	"unicode"

	"github.com/osuushi/vivid/vivian"
)

// Tag names that cause a cell to be created, and therefore must be moved to the
// top level.
var cellCreatingNamePrefixes = []string{
	"min",
	"max",
	"wrap",
	"fixed",
	"left",
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
