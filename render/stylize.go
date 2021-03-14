package render

import (
	"fmt"
	"strings"

	"github.com/osuushi/vivid/rich"
	"github.com/osuushi/vivid/vivian"
	"github.com/spf13/cast"
)

/*
A cell's content consists of nothing but style, text, and input nodes. Once our
template context is provided, we can produce a RichString to put in the cell.
This is completely isolated from layout.
*/

func stylizeNodes(nodes []vivian.Node, context interface{}, style *rich.Style) (rich.RichString, error) {
	parts := []rich.RichString{}
	for _, node := range nodes {
		part, err := stylizeNode(node, context, style)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return rich.Concat(parts...), nil
}

func stylizeNode(node vivian.Node, context interface{}, style *rich.Style) (rich.RichString, error) {
	switch node := node.(type) {
	case *vivian.ContentNode:
		newStyle, err := styleFromContentNode(node.Tag, style)
		if err != nil {
			return nil, err
		}
		return stylizeNodes(node.Children, context, newStyle)
	case *vivian.InputNode:
		val := cast.ToString(accessContext(context, node.Path))
		return rich.NewRichString(val, style), nil
	case *vivian.TextNode:
		return rich.NewRichString(node.Text, style), nil
	}

	return nil, fmt.Errorf("Unknown node type for %v", node)
}

func styleFromContentNode(tag string, parentStyle *rich.Style) (*rich.Style, error) {
	// Tags are case insensitive
	tag = strings.ToLower(tag)
	style := &rich.Style{
		Parent: parentStyle,
	}
	switch {
	case tag == "bold" || tag == "b":
		style.Bold = rich.On
	case tag == "italic" || tag == "i":
		style.Italic = rich.On
	case tag == "underline" || tag == "u":
		style.Underline = rich.On
	default:
		if color, ok := parseColor(tag); ok {
			style.Color = color
		} else {
			return nil, fmt.Errorf("Unrecognized tag: %q", tag)
		}
	}
	return style, nil
}
