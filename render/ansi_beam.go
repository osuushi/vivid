package render

import (
	"strconv"
	"strings"

	"github.com/osuushi/vivid/rich"
)

// StyleBeam which converts to ANSI formatting (24-bit color)

type ANSIBeam struct {
	// Track the style we're currently rendering so we know which new escape codes
	// need to be emitted
	currentStyle *rich.Style
}

func (beam *ANSIBeam) ScanRune(r rich.RichRune, b *strings.Builder) {

}

func (beam *ANSIBeam) Terminate(b *strings.Builder) {

}

func writeSGR(content string, b *strings.Builder) {
	b.WriteString(CSI)
	b.WriteString(content)
	b.WriteRune(SGRSuffix)
}

// True-color ansi code from RGB
func WriteSGRColor(background bool, color *rich.RGB, b *strings.Builder) {
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
	b.WriteRune(SGRSuffix)
}
