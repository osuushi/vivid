// +build !wasm

package render

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func DefaultBeam() StyleBeam {
	if !isTTY() {
		return &PlainBeam{}
	}
	ansiBeam := &ANSIBeam{}
	ansiBeam.UseColor = supports256Color()
	ansiBeam.TrueColor = supportsTrueColor()
	return ansiBeam
}

func TerminalWidth() (int, error) {
	tputAnswer, err := tput("cols")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(tputAnswer)
}

func tput(arg string) (string, error) {
	result, err := exec.Command("tput", arg).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), nil
}

func isTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func supports256Color() bool {
	tputAnswer, err := tput("colors")
	if err != nil {
		return false
	}
	return tputAnswer == "256"
}

func supportsTrueColor() bool {
	// Unfortunately, detection of true color is pretty spotty. Terminals will
	// report this env var, but not over SSH
	envVar := os.Getenv("COLORTERM")
	return envVar == "truecolor" || envVar == "24bit"
}
