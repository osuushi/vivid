package rich

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestRGBFromHex(t *testing.T) {
	check := func(str string, valid bool, expectedComponents ...uint8) {
		actual, ok := RGBFromHex(str)
		if ok != valid {
			word := "invalid"
			if valid {
				word = "valid"
			}
			t.Errorf("Expected %q to be %s", str, word)
		}

		if valid {
			expected := &RGB{
				R: expectedComponents[0],
				G: expectedComponents[1],
				B: expectedComponents[2],
			}

			if diff := deep.Equal(actual, expected); diff != nil {
				t.Errorf("For input %q:\n%s", str, strings.Join(diff, "\n"))
			}
		}
	}

	check("foo", false)
	check("", false)
	check("fff", true, 0xff, 0xff, 0xff)
	check("aaa", true, 0xaa, 0xaa, 0xaa)
	check("ADA", true, 0xaa, 0xdd, 0xaa)
	check("019", true, 0x00, 0x11, 0x99)
	check("ffffff", true, 0xff, 0xff, 0xff)
	check("accede", true, 0xac, 0xce, 0xde)
	check("DEC0DE", true, 0xde, 0xc0, 0xde)
}
