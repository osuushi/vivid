package rich

import (
	"image/color"
	"strconv"
)

type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	// Easiest way to match color.RGBA's behavior is just to convert and call it.
	return color.RGBA{c.R, c.G, c.B, 0xff}.RGBA()
}

// Note that this is not correct if A < 1, since color.RGBA represents premul
// alpha. However, we're never converting from a transparent source, so this is
// moot.
func RGBA2RGB(c color.RGBA) *RGB {
	return &RGB{
		R: c.R,
		G: c.G,
		B: c.B,
	}
}

func RGBFromHex(hex string) (*RGB, bool) {
	// Validate hex
	val, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return nil, false
	}

	// If shorthand, double the digits
	if len(hex) == 3 {
		sixDigit := make([]rune, 6)
		for i, r := range hex {
			sixDigit[i*2] = r
			sixDigit[i*2+1] = r
		}
		hex = string(sixDigit)
		val, _ = strconv.ParseUint(hex, 16, 64)
	}

	// Now normalized that we're sure it's six digit hex
	result := &RGB{}
	result.B = uint8(val & 0xff)
	val = val >> 8
	result.G = uint8(val & 0xff)
	val = val >> 8
	result.R = uint8(val & 0xff)
	return result, true
}
