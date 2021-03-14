package render

// Control sequence introducer
const CSI = "\033["
const SGRSuffix = 'm'

// SGR codes
const (
	SGRReset       = "0"
	SGRBold        = "1"
	SGRItalic      = "3"
	SGRUnderline   = "4"
	SGRFgTrueColor = "38;2;"
	SGRBgTrueColor = "48;2;"
	SGRFg256Color  = "38;5;"
	SGRBg256Color  = "48;5;"

	SGRNotBold      = "22"
	SGRNotItalic    = "23"
	SGRNotUnderline = "24"

	SGRFgReset = "39"
	SGRBgReset = "49"
)
