package render

import (
	"strings"

	"github.com/osuushi/vivid/rich"
	"golang.org/x/image/colornames"
)

func parseColor(tag string) (*rich.RGB, bool) {
	// First check html color names
	rgba, ok := colornames.Map[tag]
	if ok {
		return rich.RGBA2RGB(rgba), true
	}

	return rich.RGBFromHex(tag)
}

func parseBgColor(tag string) (*rich.RGB, bool) {
	if !strings.HasPrefix(tag, "bg") {
		return nil, false
	}
	return parseColor(tag[2:])
}

func isColor(tag string) bool {
	_, ok := parseColor(tag)
	return ok
}

func isBgColor(tag string) bool {
	_, ok := parseBgColor(tag)
	return ok
}
