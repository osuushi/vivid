package vivian

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Control characters. We normalize delimiters/braces to these before parsing,
// and restore them after parsing. This makes the parser itself delimiter
// agnostic.
const (
	tagMarker  = '\x01'
	openBrace  = '\x02'
	closeBrace = '\x03'
)

const allowedDelimiters = "!#$%^&*"
const allowedBraces = "[]()<>{}"

func ParseString(str string) (*Ast, error) {
	ast := &Ast{
		TagMarker:  '@',
		OpenBrace:  '[',
		CloseBrace: ']',
	}

	// Handle delimiter mode
	if strings.HasPrefix(str, "@") {

	}

	stringReader := bufio.NewReader(strings.NewReader(str))
	reader, writer := io.Pipe()
	go transformInput(ast, stringReader, writer)

	result, err := ParseReader("input", reader)
	if err != nil {
		return nil, err
	}

	content, ok := result.(*ContentNode)
	if !ok {
		return nil, fmt.Errorf(
			"Unexpected type for parsed content: %s",
			reflect.TypeOf(result).String(),
		)
	}

	ast.Content = content
	transformOutput(ast, ast)

	return ast, nil
}

// Only panics in here because the only ways for these to fail are catastrophic
// internal errors.
func transformInput(ast *Ast, stringReader *bufio.Reader, writer *io.PipeWriter) {
	bufWriter := bufio.NewWriter(writer)
	// Wrap with the root tag
	bufWriter.WriteRune(tagMarker)
	bufWriter.WriteString("root")
	bufWriter.WriteRune(openBrace)

	for {
		char, _, err := stringReader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		switch char {
		case ast.TagMarker:
			char = tagMarker
		case ast.OpenBrace:
			char = openBrace
		case ast.CloseBrace:
			char = closeBrace
		}

		_, err = bufWriter.WriteRune(char)
		if err != nil {
			panic(err)
		}
	}

	// Close the root tag
	bufWriter.WriteRune(closeBrace)

	// Flush and close the writer
	err := bufWriter.Flush()
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}
}

// Undo the delimiter transformation to restore original characters
func transformOutput(ast *Ast, currentNode interface{}) {
	switch currentNode.(type) {
	case *Ast:
		transformOutput(ast, currentNode.(*Ast).Content)
	case *ContentNode:
		for _, child := range currentNode.(*ContentNode).Children {
			transformOutput(ast, child)
		}
	case *TextNode:
		textNode := currentNode.(*TextNode)
		text := textNode.Text
		newText := make([]rune, len(text))
		for i, r := range text {
			switch r {
			case tagMarker:
				r = ast.TagMarker
			case openBrace:
				r = ast.OpenBrace
			case closeBrace:
				r = ast.CloseBrace
			}
			newText[i] = r
		}
		textNode.Text = string(newText)
	}
}
