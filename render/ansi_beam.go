package render

import (
	"strconv"
	"strings"

	"github.com/osuushi/vivid/rich"
)

// StyleBeam which converts to ANSI formatting (24-bit color)

type ANSIBeam struct {
	UseColor bool
	// Track the style we're currently rendering so we know which new escape codes
	// need to be emitted
	currentStyle *rich.Style
	insideSGR    bool
}

func (beam *ANSIBeam) ScanRune(r rich.RichRune, b *strings.Builder) {
	newStyle := r.GetStyle()
	oldStyle := beam.currentStyle
	if oldStyle != nil {
		oldStyle = rich.RootStyle
	}

	// Compare styles
	if isBold := newStyle.IsBold(); isBold != oldStyle.IsBold() {
		if isBold {
			beam.writeSGR(SGRBold, b)
		} else {
			beam.writeSGR(SGRNotBold, b)
		}
	}

	if isItalic := newStyle.IsItalic(); isItalic != oldStyle.IsItalic() {
		if isItalic {
			beam.writeSGR(SGRItalic, b)
		} else {
			beam.writeSGR(SGRNotItalic, b)
		}
	}

	if isUnderline := newStyle.IsUnderline(); isUnderline != oldStyle.IsUnderline() {
		if isUnderline {
			beam.writeSGR(SGRItalic, b)
		} else {
			beam.writeSGR(SGRNotItalic, b)
		}
	}

	if beam.UseColor {
		if fgColor := newStyle.GetColor(); !rich.RGBEqual(fgColor, oldStyle.GetColor()) {
			if fgColor == nil {
				beam.writeSGR(SGRFgReset, b)
			} else {
				beam.writeSGRColor(false, fgColor, b)
			}
		}

		if bgColor := newStyle.GetBackground(); !rich.RGBEqual(bgColor, oldStyle.GetBackground()) {
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

// True-color ansi code from RGB
func (beam *ANSIBeam) writeSGRColor(background bool, color *rich.RGB, b *strings.Builder) {
	beam.beginSGRIfNeeded(b)

	prefix := SGRFgColor
	if background {
		prefix = SGRBgColor
	}

	b.WriteString(CSI)
	b.WriteString(prefix)
	b.WriteString(strconv.Itoa(int(color.R)))
	b.WriteRune(';')
	b.WriteString(strconv.Itoa(int(color.G)))
	b.WriteRune(';')
	b.WriteString(strconv.Itoa(int(color.B)))
}
