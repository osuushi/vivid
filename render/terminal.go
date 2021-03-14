// +build !wasm

package render

import (
	"os"

	"github.com/olekukonko/ts"
)

func DefaultBeam() StyleBeam {
	if !isTTY() {
		return &PlainBeam{}
	}
	ansiBeam := &ANSIBeam{}
	ansiBeam.UseColor = supportsTrueColor()
	return ansiBeam
}

func TerminalWidth() (int, error) {
	if !isTTY() {
		// TODO: Need a better choice for this
		return 250, nil
	}
	size, err := ts.GetSize()
	if err != nil {
		return 0, err
	}

	return size.Col(), nil
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
