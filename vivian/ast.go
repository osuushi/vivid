package vivian

// Root node
type Ast struct {
	// What character is used to denote the start of a tag; default is @
	TagMarker rune
	// What character opens a tag (default <)
	OpenBrace rune
	// What character closes a tag (default >)
	CloseBrace rune

	// The actual content
	Content *ContentNode
}

// Tagged content
type ContentNode struct {
	// Tag names, which are alphanumeric strings. The parser is agnostic to the
	// actual value of these tags. Validation and parsing of parameterized tags is
	// left to the consumer.
	Tag string

	// Children can be ContentNodes, TextNodes, or InterpolationNodes
	Children []interface{}
}

// A node representing an input from the template context
type InputNode struct {
	// Path to the property, which will be dot separated identifiers of the form [A-Za-z1-9_]
	Path []string
}

type TextNode struct {
	Text string
}
