package render

// Control sequence introducer
const CSI = "\033["
const SGRSuffix = 'm'

// SGR codes
const (
	SGRReset     = "0"
	SGRBold      = "1"
	SGRItalic    = "3"
	SGRUnderline = "4"
	SGRFgColor   = "38;2;"
	SGRBgColor   = "48;2;"

	SGRNotBold      = "22"
	SGRNotItalic    = "23"
	SGRNotUnderline = "24"

	SGRFgReset = "39"
	SGRBgReset = "49"
)
