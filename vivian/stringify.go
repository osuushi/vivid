package vivian

import (
	"fmt"
	"regexp"
	"strings"
)

func (ast *Ast) String() string {
	header := fmt.Sprintf("%c%c%c", ast.TagMarker, ast.OpenBrace, ast.CloseBrace)
	if header == "@[]" {
		header = ""
	} else if header[0] != '@' {
		header = "@" + header
	}

	innerStrings := make([]string, len(ast.Content.Children))

	children := ast.Content.Children
	for i, child := range children {
		currentString := child.String(ast)
		innerStrings[i] = currentString
		// Check if we need to add a chomp
		if i > 0 {
			currentString := innerStrings[i]
			_, ok := children[i-1].(*InputNode)
			// If input precedes a string that doesn't start with a space, use a chomp.
			if ok {
				matched, err := regexp.MatchString("^\\s", currentString)
				if err != nil {
					panic(err)
				}
				if !matched {
					innerStrings[i-1] += "~ "
				}
			}
		}
	}
	return header + strings.Join(innerStrings, "")
}

func (node *ContentNode) String(root *Ast) string {
	innerStrings := make([]string, len(node.Children))
	for i, child := range node.Children {
		innerStrings[i] = child.String(root)
	}

	return string(root.TagMarker) +
		node.Tag +
		string(root.OpenBrace) +
		strings.Join(innerStrings, "") +
		string(root.CloseBrace)
}

func (node *InputNode) String(root *Ast) string {
	pathString := strings.Join(node.Path, ".")
	return string(root.TagMarker) + "-" + pathString
}

func (node *TextNode) String(root *Ast) string {
	text := node.Text
	text = strings.ReplaceAll(
		text,
		string(root.TagMarker),
		string(root.TagMarker)+string(root.TagMarker),
	)
	text = strings.ReplaceAll(
		text,
		string(root.CloseBrace),
		string(root.TagMarker)+string(root.CloseBrace),
	)
	return text
}
