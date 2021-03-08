package render

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/osuushi/vivid/vivian"
)

// This takes all cell-creating nodes and moves them to the top level, splitting
// up existing nodes as needed.
func hoistCells(ast *vivian.Ast) error {
	_, err := hoistContentNode(ast.Content)
	return err
}

// Tag names that cause a cell to be created, and therefore must be moved to the
// top level.
var cellCreatingNamePrefixes = []string{
	"min",
	"max",
	"auto",
	"wrap",
	"fixed",
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

func isRootNode(node vivian.Node) bool {
	contentNode, ok := node.(*vivian.ContentNode)
	if !ok { // Non-content node never a cell creator
		return false
	}
	return contentNode.Tag == "_root"
}

func isCellCreatorNode(node vivian.Node) bool {
	contentNode, ok := node.(*vivian.ContentNode)
	if !ok { // Non-content node never a cell creator
		return false
	}
	return isCellCreator(contentNode.Tag)
}

// Copy a content node, but not its children
func copyContentNode(node *vivian.ContentNode) *vivian.ContentNode {
	newStruct := *node
	newStruct.Children = []vivian.Node{}
	return &newStruct
}

// Recursively hoist all content nodes, replacing this node if necessary with
// new children for the parent.
func hoistContentNode(node *vivian.ContentNode) ([]vivian.Node, error) {
	// If the node is childless, hoisting is a noop
	if node.Children == nil {
		return []vivian.Node{node}, nil
	}

	err := hoistChildren(node)
	if err != nil {
		return nil, err
	}

	if isRootNode(node) {
		return nil, nil
	} else if isCellCreator(node.Tag) {
		// Cell creators require validation, but no manipulation once teir children
		// have been hoisted above.
		err = validateCellCreatorChildren(node)
		if err != nil {
			return nil, err
		}
		// We never split a cell creator
		return []vivian.Node{node}, nil
	} else {
		return hoistStyle(node)
	}
}

// Hoist all children first. This process may turn any child node into
// multiple child nodes, so it replaces the child list
func hoistChildren(node *vivian.ContentNode) error {
	if node.Children == nil {
		return nil
	}

	newChildren := []vivian.Node{}
	for _, child := range node.Children {
		// Recursively hoist child nodes
		childContentNode, ok := child.(*vivian.ContentNode)
		if ok {
			// Recurse for content node
			replacement, err := hoistContentNode(childContentNode)
			if err != nil {
				return err
			}
			newChildren = append(newChildren, replacement...)
		} else {
			// Not a content node, so it can't have children
			newChildren = append(newChildren, child)
		}
	}

	node.Children = newChildren
	return nil
}

func hoistStyle(node *vivian.ContentNode) ([]vivian.Node, error) {
	if node.Children == nil {
		return []vivian.Node{node}, nil
	}

	anyCellCreators := false
	for _, child := range node.Children {
		if isCellCreatorNode(child) {
			anyCellCreators = true
			break
		}
	}
	// If the noded doesn't contain a cell creator, hoisting is a no-op. Note that
	// it is guaranteed that any cell creator descendents at this point will be a
	// consecutive lineage. So not finding a cell creator in the children means
	// that this entire subtree is styling.
	if !anyCellCreators {
		return []vivian.Node{node}, nil
	}

	// The existence of a cell creator child guarantees that this node will be
	// replaced
	replacement := []vivian.Node{}

	// We use this to coalesce non-creator nodes into one split, so that we don't
	// split to every child individually.
	currentRun := []vivian.Node{}
	appendCoalescedRun := func() {
		if len(currentRun) == 0 {
			return
		}
		clone := copyContentNode(node)
		clone.Children = currentRun
		replacement = append(replacement, clone)
		currentRun = []vivian.Node{}
	}

	for _, child := range node.Children {
		// Style node has to be injected to the end of the cell creator chain for
		// cell creators
		if isCellCreatorNode(child) {
			appendCoalescedRun()
			injectStyle(node, child.(*vivian.ContentNode))
			replacement = append(replacement, child)
		} else { // Not a cell creator
			// Gather consecutive non-cell-creators up so they can be added to a clone
			// later
			currentRun = append(currentRun, child)
		}
	}
	// Coalesce the last run if there is one
	appendCoalescedRun()

	return replacement, nil
}

func injectStyle(styleNode, cellCreator *vivian.ContentNode) {
	// Injection must happen at the end of the cell creator lineage, so recurse if
	// we're not there.
	endOfLineage :=
		cellCreator.Children == nil ||
			len(cellCreator.Children) != 1 || // guaranteed by validateCellCreatorChildren
			!isCellCreatorNode(cellCreator.Children[0])

	if !endOfLineage {
		injectStyle(
			styleNode,
			cellCreator.Children[0].(*vivian.ContentNode),
		)
		return
	}

	// Insert style node clone
	clone := copyContentNode(styleNode)
	clone.Children = cellCreator.Children
	cellCreator.Children = []vivian.Node{clone}
}

// Once descendents are hoisted, a cell creator can have:
// 1. Zero children
// 2. One cell creator child
// 3. Multiple children that are not cell creators
//
// A cell creator cannot have a _mix_ of styles and cell creators as children,
// as this would imply a cell subdivision, which is not allowed.
func validateCellCreatorChildren(node *vivian.ContentNode) error {
	// One or fewer children is already allowed
	if node.Children == nil || len(node.Children) <= 1 {
		return nil
	}

	// For multiple children, no child may be a cell creator
	for _, child := range node.Children {
		if isCellCreatorNode(child) {
			return fmt.Errorf(
				"Sizing element %q cannot be subdivided by %q.\n"+
					"  May contain exactly one sizing element, OR zero or more style elements.",
				node.Tag,
				child.(*vivian.ContentNode).Tag,
			)
		}
	}
	return nil
}
