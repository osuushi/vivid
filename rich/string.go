package rich

// A RichString is represented as a slice of RichRunes, with core string methods
// reimplemented. RichStrings are mutable and are safe to modify in place
// (although not thread safe).

type RichString []RichRune

// Create a new rich string with a uniform style. Pass nil for an unstyled
// string.
func NewRichString(str string, style *Style) RichString {
	if style == nil {
		style = rootStyle
	}

	runes := []rune(str)
	richRunes := make(RichString, len(runes))
	for i, r := range runes {
		richRunes[i].Rune = r
		richRunes[i].style = style
	}
	return richRunes
}

// Linked list type used internally
type stringList struct {
	value RichString
	next  *stringList
}

func (l *stringList) length() int {
	if l.next == nil {
		return 1
	}
	return l.next.length() + 1
}

func (l *stringList) toSlice() []RichString {
	result := make([]RichString, l.length())
	for i := range result {
		result[i] = l.value
		l = l.next
	}
	return result
}

// Check if s is a prefix of r ignoring formatting
func (r RichString) HasStringPrefix(s string) bool {
	// Prefix can't be longer than the string itself
	if len(s) > len(r) {
		return false
	}
	for i, rune := range s {
		if r[i].Rune != rune {
			return false
		}
	}
	return true
}

// This is useful for debugging and testing etc. but it doesn't render
// formatting.
func (r RichString) String() string {
	runes := make([]rune, len(r))
	for i, r := range r {
		runes[i] = r.Rune
	}
	return string(runes)
}

// Check if a RichString is equal to a bare string, ignoring formatting
func (r RichString) EqualsString(s string) bool {
	return len(s) == len(r) && r.HasStringPrefix(s)
}

// Helper that splits the string recursively by creating a linked list
func (r RichString) listSplit(delim string) *stringList {
	for i := range r {
		if r[i:].HasStringPrefix(delim) {
			return &stringList{
				// Value is everything up until the delimiter
				value: r[:i],
				next:  r[i+len(delim):].listSplit(delim),
			}
		}
	}
	// If we made it out of the loop, there was no delimiter, so we just return
	// the entire string (this will always happen for the last substring)
	return &stringList{
		value: r,
		next:  nil,
	}
}

func (r RichString) splitRunes() []RichString {
	result := make([]RichString, len(r))
	for i := range r {
		result[i] = r[i : i+1]
	}
	return result
}

// Similar to strings.Split. Note that the result will be a reslicing of the
// original string, not a copy.
func (r RichString) Split(delim string) []RichString {
	if delim == "" {
		// Special case; listSplit would just stack overflow on a list of empty
		// strings
		return r.splitRunes()
	}
	return r.listSplit(delim).toSlice()
}

// Concatenate strings into a new string.
func Concat(richStrings ...RichString) RichString {
	totalLength := 0
	for _, r := range richStrings {
		totalLength += len(r)
	}

	// Preallocate so that appends don't have to resize the slice.
	result := make(RichString, 0, totalLength)
	for _, r := range richStrings {
		result = append(result, r...)
	}

	return result
}
