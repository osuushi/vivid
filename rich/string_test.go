package rich

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestHasStringPrefix(t *testing.T) {
	checkPrefix := func(str string, prefix string, expected bool) {
		richStr := NewRichString(str, nil)
		if richStr.HasStringPrefix(prefix) != expected {
			t.Errorf("%q.HasStringPrefix(%q) expected %v", str, prefix, expected)
		}
	}

	checkPrefix("blah blah", "blah", true)
	checkPrefix("blah blah", "blahx", false)
	checkPrefix("foo bar", "oo b", false)
	checkPrefix("hello", "", true)
}

func TestSplit(t *testing.T) {
	checkSplit := func(str string, delim string, expected ...string) {
		richString := NewRichString(str, nil)
		split := richString.Split(delim)
		actual := []string{}
		for _, s := range split {
			actual = append(actual, s.String())
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf(
				"%q.Split(%q) expected %v but got %v",
				str,
				delim,
				expected,
				actual,
			)
		}
	}

	checkSplit("foo bar baz", " ", "foo", "bar", "baz")
	checkSplit("foo", "", "f", "o", "o")
	checkSplit("this, is, a, test", ", ", "this", "is", "a", "test")
	checkSplit("::foo::bar::", "::", "", "foo", "bar", "")
}

func TestConcat(t *testing.T) {
	actual := Concat(
		NewRichString("foo", &Style{Bold: On}),
		NewRichString("bar", nil),
		NewRichString("baz", &Style{Bold: On}),
	)

	actualString := actual.String()
	actualBolds := make([]bool, len(actual))
	for i, richRune := range actual {
		actualBolds[i] = richRune.IsBold()
	}

	expectedString := "foobarbaz"
	expectedBolds := []bool{
		true, true, true, // foo
		false, false, false, //bar
		true, true, true, // baz
	}

	if actualString != expectedString {
		t.Errorf("Expected %q but got %q", expectedString, actualString)
	}

	if diff := deep.Equal(expectedBolds, actualBolds); diff != nil {
		t.Error(diff)
	}
}
