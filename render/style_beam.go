package render

import (
	"strings"

	"github.com/osuushi/vivid/rich"
)

// Interface for rendering a rich string by scanning it from left to right.
type StyleBeam interface {
	ScanRune(rich.RichRune, *strings.Builder)
	// Terminate a line. The beam should be restored to its initial state, and
	// must be reusable for multiple lines
	Terminate(*strings.Builder)
}

// Simplest beam which discards all style
type PlainBeam struct{}

func (beam *PlainBeam) ScanRune(r rich.RichRune, b *strings.Builder) {
	b.WriteRune(r.Rune)
}

// Nothing has to be done to terminate with the plain beam
func (beam *PlainBeam) Terminate(*strings.Builder) {}
