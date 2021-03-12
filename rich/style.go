package rich

// Style for a RichRune.
//
// Styles should be constructed via struct literals, but typically you should
// access their fields through the provided getters, which will proved proper
// inheritance
//
// Note that it is generally unsafe to modify a style once it has been sent to a
// method in the rich package.

type Style struct {
	Color, Background       *RGB
	Bold, Italic, Underline Tristate
	Parent                  *Style // if nil, the root Style is implied
}

var rootStyle = &Style{
	// Note that "nil" has a special meaning for colors in the root style. A nil
	// color in any non-root style will mean "inherit", but a nil style in the
	// root will use the "normal" color provided by the reset escape sequence.
	// This is the only supported method of accessing the "normal" style - by not
	// applying a color in the first place, and is offered so that lightweight
	// usecases don't have to think about color.
	Bold:      Off,
	Italic:    Off,
	Underline: Off,
}

func (s *Style) getParent() *Style {
	if s == rootStyle {
		panic("Tried to access parent of root style")
	}
	if s.Parent == nil {
		return rootStyle
	} else {
		return s.Parent
	}
}

// Inheritance helper for color access.
func (s *Style) inheritColor(get func(s *Style) *RGB) *RGB {
	if val := get(s); val != nil {
		return val
	} else if s == rootStyle {
		// nil is actually allowed in root style, albeit with a different meaning.
		return nil
	} else {
		// Inherit
		return s.getParent().inheritColor(get)
	}
}

func (s *Style) GetColor() *RGB {
	return s.inheritColor(func(s *Style) *RGB { return s.Color })
}

func (s *Style) GetBackground() *RGB {
	return s.inheritColor(func(s *Style) *RGB { return s.Background })
}

// Inheritance helper for Tristate access. Since the value must eventually be
// set, always returns a boolean.
//
// Note the implication that no root style Tristate may be Unset, or this will
// panic when it tries to get the parent from the root style.
func (s *Style) inheritTristate(get func(*Style) Tristate) bool {
	if val := get(s); val == Unset {
		// If unset, inherit
		return s.getParent().inheritTristate(get)
	} else {
		// If set, convert to boolean
		return val == On
	}
}

func (s *Style) IsBold() bool {
	return s.inheritTristate(func(s *Style) Tristate { return s.Bold })
}

func (s *Style) IsItalic() bool {
	return s.inheritTristate(func(s *Style) Tristate { return s.Italic })
}

func (s *Style) IsUnderline() bool {
	return s.inheritTristate(func(s *Style) Tristate { return s.Underline })
}

// Inject newRoot at the root, creating a new style
func (s *Style) Rebase(newRoot *Style) *Style {
	if s == rootStyle || s == nil {
		return newRoot
	}

	// Shallow clone by dereference
	styleStruct := *s
	styleStruct.Parent = s.Parent.Rebase(newRoot)
	return &styleStruct
}
