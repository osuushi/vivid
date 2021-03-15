package render

import (
	"fmt"
	"unicode"

	"github.com/osuushi/vivid/rich"
)

// Rendering of an individual cell, given a RichString.
// Output is an array of RichString rows.
//
// NOTE: This will destroy the input string. Most of the layout work is done
// in-place, including space compaction and trimming. That reduces allocations
// and copies, but it will garble the input argument.
func renderCell(content rich.RichString, sizedCell *SizedCell) []rich.RichString {
	return renderContent(
		content,
		sizedCell.Width,
		sizedCell.Cell.Wrap,
		sizedCell.Cell.Alignment,
	)
}

func renderContent(content rich.RichString, width int, wrap bool, alignment Alignment) []rich.RichString {
	paragraphs := content.Split("\n")
	paragraphs = trimEmptyParagraphs(paragraphs)
	if !wrap {
		// Use only the first pargraph and truncate it
		line := normalizeWhitespace(rich.Concat(paragraphs...))
		paragraphs = []rich.RichString{
			truncateContentToWidth(line, width),
		}
	}
	rows := []rich.RichString{}
	for _, p := range paragraphs {
		paragraphRows := sliceParagraph(p, width)
		alignParagraphRows(paragraphRows, alignment, width)
		rows = append(rows, paragraphRows...)
	}

	return rows
}

func trimEmptyParagraphs(paragraphs []rich.RichString) []rich.RichString {
	for len(paragraphs) > 1 && isAllWhiteSpace(paragraphs[0]) {
		paragraphs = paragraphs[1:]
	}

	return paragraphs
}

func isAllWhiteSpace(line rich.RichString) bool {
	for _, r := range line {
		if !unicode.IsSpace(r.Rune) {
			return false
		}
	}
	return true
}

func truncateContentToWidth(content rich.RichString, width int) rich.RichString {
	if len(content) <= width {
		return content
	}

	if width == 1 {
		// Special case. Give the elipsis the style of the first character
		style := content[0].GetStyle()
		return rich.NewRichString("…", style)
	} else {
		// Slice to truncate, leaving room for ellipsis
		return content[:width-1].Append("…")
	}
}

func normalizeWhitespace(content rich.RichString) rich.RichString {
	if len(content) == 0 {
		// Since the other case returns a copy, so should this. It's very cheap, and
		// this defends against bugs later.
		return rich.RichString{}
	}

	// Make in place by writing back to the original array. Since len(result) <=
	// len(content), this will never have to resize
	result := content[:0]
	// Consider the string to be left-extended with spaces; we trim leading and
	// trailing space.
	wasSpace := true
	for _, c := range content {
		r := c.Rune
		// Replace all whitespace with spaces
		if unicode.IsSpace(r) && r != ' ' {
			r = ' '
		}

		// Collapse spaces by skipping if adjacent
		if r == ' ' && wasSpace {
			continue
		}
		c.Rune = r
		result = append(result, c)
		wasSpace = r == ' '
	}

	// Check if last character is a space. If so, omit it
	if len(result) > 0 && result[len(result)-1].Rune == ' ' {
		result = result[:len(result)-1]
	}

	return result
}

// Convert a single block of text (with no newlines) and split it into rows of a
// given maximum width (without alignment)
func sliceParagraph(text rich.RichString, width int) []rich.RichString {
	text = normalizeWhitespace(text)
	lines := []rich.RichString{}

	for text != nil {
		var currentLine rich.RichString
		currentLine, text = scanNextLine(text, width)

		lines = append(lines, currentLine)
	}
	return lines
}

// Main workhorse of sliceParagraph. Returns a single line and a remainder, both
// slices of the original string (saving one copy). Returns nil for second argument if there is
// nothing remaining.
func scanNextLine(input rich.RichString, width int) (rich.RichString, rich.RichString) {
	if len(input) == 0 {
		return input, nil
	}

	// Ignore leading space
	if input[0].Rune == ' ' {
		input = input[1:]
	}

	lastSpaceIndex := -1
	for i, r := range input {
		if r.Rune == ' ' {
			lastSpaceIndex = i
		}

		if i >= width { // We hit the max line length
			if lastSpaceIndex > 0 { // Normal case where we have a space to split at
				return input[:lastSpaceIndex], input[lastSpaceIndex:]
			} else { // A single word has occupied the entire line
				return input[:i], input[i:]
			}
		}
	}

	// If we make it out of the loop, the string ended before we hit the width.
	return input, nil
}

func alignParagraphRows(rows []rich.RichString, alignment Alignment, width int) {
	// All but last row
	for i, row := range rows[:len(rows)-1] {
		rows[i] = alignRow(row, alignment, width)
	}
	// Alignment switches to left for Justify on the last row of a paragraph
	if alignment == Justify {
		alignment = Left
	}
	rows[len(rows)-1] = alignRow(rows[len(rows)-1], alignment, width)
}

func alignRow(row rich.RichString, alignment Alignment, width int) rich.RichString {
	// Copy the row, since up to this point, we will have used slices to avoid
	// copies
	row = append(rich.RichString{}, row...)
	switch alignment {
	case Left:
		return append(row, rich.MakeSpacer(width-len(row), nil)...)
	case Right:
		return append(rich.MakeSpacer(width-len(row), nil), row...)
	case Center:
		return centerRow(row, width)
	case Justify:
		return justifyLine(row, width)
	}
	panic(fmt.Sprintf("Unknown alignment: %v", alignment))
}

func centerRow(row rich.RichString, width int) rich.RichString {
	freeSpace := width - len(row)
	leftSide := freeSpace / 2
	row = append(rich.MakeSpacer(leftSide, nil), row...)
	rightSide := freeSpace - leftSide
	row = append(row, rich.MakeSpacer(rightSide, nil)...)
	return row
}

// Justify a line of text. This assumes that it is not the last line in a
// paragraph, which is always left-aligned.
func justifyLine(line rich.RichString, width int) rich.RichString {
	totalSpace := width - len(line)

	// With no space to allocate, we can just return the original string
	if totalSpace == 0 {
		return line
	}

	originalSpaces := 0

	// First scan the number of spaces
	for _, r := range line {
		if r.Rune == ' ' {
			originalSpaces += 1
		}
	}

	// Allocate a new amount of space to each original space. We add the original
	// spaces to totalSpace, since they were counted in the length
	allocator := dither(totalSpace+originalSpaces, originalSpaces)

	result := make(rich.RichString, 0, width)
	for _, r := range line {
		if r.Rune == ' ' {
			// Place a spacer
			result = append(result, rich.MakeSpacer(allocator(), nil)...)
		} else {
			result = append(result, r)
		}
	}
	return result
}
