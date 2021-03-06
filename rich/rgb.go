package rich

import "image/color"

type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	// Easiest way to match color.RGBA's behavior is just to convert and call it.
	return color.RGBA{c.R, c.G, c.B, 0xff}.RGBA()
}
