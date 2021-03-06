package rich

type RichRune struct {
	// Note that this field is private. This is important because it is completely
	// unsafe to modify a rune's style in place. The style is a pointer as a way
	// of building in some compression in the normal ways of constructing rich
	// strings. When you construct a rich string, you will typically use a single
	// style for every rune, in which case every rune will point at the same
	// style.
	style *Style
	Rune  rune
}

func (r *RichRune) GetColor() *RGB {
	return r.style.GetColor()
}

func (r *RichRune) GetBackground() *RGB {
	return r.style.GetBackground()
}

func (r *RichRune) IsBold() bool {
	return r.style.IsBold()
}

func (r *RichRune) IsItalic() bool {
	return r.style.IsItalic()
}

func (r *RichRune) IsUnderline() bool {
	return r.style.IsUnderline()
}
