package rich

type RichRune struct {
	Style *Style
	Rune  rune
}

func (r *RichRune) GetColor() *RGB {
	return r.Style.GetColor()
}

func (r *RichRune) GetBackground() *RGB {
	return r.Style.GetBackground()
}

func (r *RichRune) IsBold() bool {
	return r.Style.IsBold()
}

func (r *RichRune) IsItalic() bool {
	return r.Style.IsItalic()
}

func (r *RichRune) IsUnderline() bool {
	return r.Style.IsUnderline()
}

// Return a copy of the rune's style
func (r *RichRune) GetStyle() *Style {
	styleStruct := *r.Style
	return &styleStruct
}
