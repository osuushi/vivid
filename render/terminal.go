// +build !wasm

package render

import "os"

func DefaultBeam() StyleBeam {
	if !isTTY() {
		return &PlainBeam{}
	}
	ansiBeam := &ANSIBeam{}
	ansiBeam.UseColor = supportsTrueColor()
	return ansiBeam
}

func isTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func supportsTrueColor() bool {
	// Unfortunately, detection of true color is pretty spotty. Terminals will
	// report this env var, but not over SSH
	envVar := os.Getenv("COLORTERM")
	return envVar == "truecolor" || envVar == "24bit"
}
