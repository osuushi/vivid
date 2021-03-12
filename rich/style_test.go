package rich

import (
	"testing"
)

func TestSimple(t *testing.T) {
	style := &Style{
		Color: &RGB{0xff, 0x00, 0x00},
		// No background
		Bold:   On,
		Italic: Off,
	}

	expectedColor := RGB{0xff, 0x00, 0x00}
	if actual := *style.GetColor(); actual != expectedColor {
		t.Errorf("style.GetColor(): Expected %v but got %v", expectedColor, actual)
	}

	if actual := style.GetBackground(); actual != nil {
		t.Errorf("style.GetBackground(): Expected %v but got %v", nil, actual)
	}

	if actual := style.IsBold(); !actual {
		t.Errorf("style.IsBold(): Expected true")
	}

	if actual := style.IsItalic(); actual {
		t.Errorf("style.IsItalic(): Expected false")
	}

	if actual := style.IsUnderline(); actual {
		t.Errorf("style.IsUnderline(): Expected false")
	}
}

func TestInherited(t *testing.T) {
	parent := &Style{
		Color:      &RGB{0xff, 0xff, 0x00},
		Background: &RGB{0xC0, 0x01, 0xD0},
		// No background
		Bold:      On,
		Underline: On,
	}

	style := &Style{
		Color:  &RGB{0xff, 0x00, 0x00},
		Bold:   Off,
		Parent: parent,
	}

	expectedColor := RGB{0xff, 0x00, 0x00}
	if actual := *style.GetColor(); actual != expectedColor {
		t.Errorf("style.GetColor(): Expected %v but got %v", expectedColor, actual)
	}

	expectedColor = RGB{0xC0, 0x01, 0xD0}
	if actual := *style.GetBackground(); actual != expectedColor {
		t.Errorf("style.GetBackground(): Expected %v but got %v", expectedColor, actual)
	}

	if actual := style.IsBold(); actual {
		t.Errorf("style.IsBold(): Expected false")
	}

	if actual := style.IsItalic(); actual {
		t.Errorf("style.IsItalic(): Expected false")
	}

	if actual := style.IsUnderline(); !actual {
		t.Errorf("style.IsUnderline(): Expected true")
	}
}

func TestRebase(t *testing.T) {
	style := &Style{
		Parent: &Style{
			Parent: &Style{},
		},
	}
	newRoot := &Style{}
	rebasedStyle := style.Rebase(newRoot)
	if rebasedStyle.Parent.Parent.Parent != newRoot {
		t.Errorf("Expected newRoot to be injected after root")
	}
}
