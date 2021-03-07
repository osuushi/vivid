package vivian

import (
	"fmt"
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

	for i, child := range ast.Content.Children {
		innerStrings[i] = child.String(ast)
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
	return ""
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
