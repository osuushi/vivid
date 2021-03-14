package rich

type RichRune struct {
	Style *Style
	Rune  rune
}

func (r *RichRune) GetColor() *RGB {
	return r.GetStyle().GetColor()
}

func (r *RichRune) GetBackground() *RGB {
	return r.GetStyle().GetBackground()
}

func (r *RichRune) IsBold() bool {
	return r.GetStyle().IsBold()
}

func (r *RichRune) IsItalic() bool {
	return r.GetStyle().IsItalic()
}

func (r *RichRune) IsUnderline() bool {
	return r.GetStyle().IsUnderline()
}

// Return a copy of the rune's style
func (r *RichRune) GetStyle() *Style {
	styleStruct := *RootStyle
	if r.Style != nil {
		styleStruct = *r.Style
	}
	return &styleStruct
}
