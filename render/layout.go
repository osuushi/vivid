package render

import (
	"fmt"
	"unicode"

	"github.com/osuushi/vivid/rich"
)

// Rendering of an individual cell, given a RichString.
// Output is an array of RichString rows.
func renderCell(sizedCell *SizedCell, content rich.RichString) []rich.RichString {
	width := sizedCell.Width
	cell := sizedCell.Cell
	paragraphs := content.Split("\n")
	if !cell.Wrap {
		// Use only the first pargraph and truncate it
		paragraphs = []rich.RichString{
			truncateContentToWidth(paragraphs[0], width),
		}
	}
	rows := []rich.RichString{}
	for _, p := range paragraphs {
		paragraphRows := sliceParagraph(p, width)
		alignParagraphRows(paragraphRows, cell.Alignment, width)
		rows = append(rows, paragraphRows...)
	}

	return rows
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
	result := make(rich.RichString, 0, len(content))
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
		if i > width { // We hit the max line length
			if lastSpaceIndex > 0 { // Normal case where we have a space to split at
				return input[:lastSpaceIndex], input[lastSpaceIndex:]
			} else { // A single word has occupied the entire line
				return input[:i], input[i:]
			}
		}

		if r.Rune == ' ' {
			lastSpaceIndex = i
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
		return append(row, makeSpacer(width-len(row))...)
	case Right:
		return append(makeSpacer(width-len(row)), row...)
	case Center:
		freeSpace := width - len(row)
		leftSide := freeSpace / 2
		row = append(makeSpacer(leftSide), row...)
		rightSide := freeSpace - leftSide
		return append(row, makeSpacer(rightSide)...)
	case Justify:
		return justifyLine(row, width)
	}
	panic(fmt.Sprintf("Unknown alignment: %v", alignment))
}

func makeSpacer(width int) rich.RichString {
	result := make(rich.RichString, width)
	for i := range result {
		result[i] = rich.RichRune{Rune: ' '}
	}
	return result
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
			result = append(result, makeSpacer(allocator())...)
		} else {
			result = append(result, r)
		}
	}
	return result
}
