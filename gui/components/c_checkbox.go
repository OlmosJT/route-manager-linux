package components

import (
	"fyne.io/fyne/v2/widget"
)

type CustomCheckbox struct {
	View *widget.Check
}

func NewCustomCheckbox(label string) *CustomCheckbox {
	return &CustomCheckbox{
		View: widget.NewCheck(label, nil),
	}
}

func (c *CustomCheckbox) IsChecked() bool {
	return c.View.Checked
}
