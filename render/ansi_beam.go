package render

import (
	"strconv"
	"strings"

	gookitColor "github.com/gookit/color"
	"github.com/osuushi/vivid/rich"
)

// StyleBeam which converts to ANSI formatting (24-bit color)

type ANSIBeam struct {
	UseColor  bool
	TrueColor bool
	// Track the style we're currently rendering so we know which new escape codes
	// need to be emitted
	currentStyle realizedStyle
	insideSGR    bool
}

type realizedStyle struct {
	Color, Background       *rich.RGB
	Bold, Italic, Underline bool
}

func (beam *ANSIBeam) ScanRune(r rich.RichRune, b *strings.Builder) {
	newStyle := r.GetStyle()
	oldStyle := &beam.currentStyle

	// Compare styles
	if isBold := newStyle.IsBold(); isBold != oldStyle.Bold {
		oldStyle.Bold = isBold
		if isBold {
			beam.writeSGR(SGRBold, b)
		} else {
			beam.writeSGR(SGRNotBold, b)
		}
	}

	if isItalic := newStyle.IsItalic(); isItalic != oldStyle.Italic {
		oldStyle.Italic = isItalic
		if isItalic {
			beam.writeSGR(SGRItalic, b)
		} else {
			beam.writeSGR(SGRNotItalic, b)
		}
	}

	if isUnderline := newStyle.IsUnderline(); isUnderline != oldStyle.Underline {
		oldStyle.Underline = isUnderline
		if isUnderline {
			beam.writeSGR(SGRUnderline, b)
		} else {
			beam.writeSGR(SGRNotUnderline, b)
		}
	}

	if beam.UseColor {
		if fgColor := newStyle.GetColor(); !rich.RGBEqual(fgColor, oldStyle.Color) {
			oldStyle.Color = fgColor
			if fgColor == nil {
				beam.writeSGR(SGRFgReset, b)
			} else {
				beam.writeSGRColor(false, fgColor, b)
			}
		}

		if bgColor := newStyle.GetBackground(); !rich.RGBEqual(bgColor, oldStyle.Background) {
			oldStyle.Background = bgColor
			if bgColor == nil {
				beam.writeSGR(SGRBgReset, b)
			} else {
				beam.writeSGRColor(true, bgColor, b)
			}
		}
	}

	beam.endSGRIfNeeded(b)
	b.WriteRune(r.Rune)
}

func (beam *ANSIBeam) Terminate(b *strings.Builder) {
	beam.currentStyle = realizedStyle{}
	beam.writeSGR(SGRReset, b)
	beam.endSGRIfNeeded(b)
}

// Start SGR sequence. If we're inside a sequence, emit a semicolon to start the
// next sequence.
func (beam *ANSIBeam) beginSGRIfNeeded(b *strings.Builder) {
	if beam.insideSGR {
		b.WriteRune(';')
	} else {
		b.WriteString(CSI)
	}
	beam.insideSGR = true
}

func (beam *ANSIBeam) endSGRIfNeeded(b *strings.Builder) {
	if !beam.insideSGR {
		return
	}

	b.WriteRune(SGRSuffix)
	beam.insideSGR = false
}

func (beam *ANSIBeam) writeSGR(sequence string, b *strings.Builder) {
	beam.beginSGRIfNeeded(b)
	b.WriteString(sequence)
}

func (beam *ANSIBeam) writeSGRColor(background bool, color *rich.RGB, b *strings.Builder) {
	if beam.TrueColor {
		beam.writeSGRTrueColor(background, color, b)
	} else {
		beam.writeSGR256Color(background, color, b)
	}
}

// True-color ansi code from RGB
func (beam *ANSIBeam) writeSGRTrueColor(background bool, color *rich.RGB, b *strings.Builder) {
	beam.beginSGRIfNeeded(b)

	prefix := SGRFgTrueColor
	if background {
		prefix = SGRBgTrueColor
	}

	b.WriteString(prefix)
	b.WriteString(strconv.Itoa(int(color.R)))
	b.WriteRune(';')
	b.WriteString(strconv.Itoa(int(color.G)))
	b.WriteRune(';')
	b.WriteString(strconv.Itoa(int(color.B)))
}

func (beam *ANSIBeam) writeSGR256Color(background bool, color *rich.RGB, b *strings.Builder) {
	beam.beginSGRIfNeeded(b)
	code := gookitColor.Rgb2short(color.R, color.G, color.B)

	prefix := SGRFg256Color
	if background {
		prefix = SGRBg256Color
	}
	b.WriteString(prefix)
	b.WriteString(strconv.Itoa(int(code)))
}
