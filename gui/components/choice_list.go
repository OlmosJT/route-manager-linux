package components

import (
	"fyne.io/fyne/v2/widget"
)

type ChoiceList struct {
	View *widget.Select
}

func NewChoiceList(options []string) *ChoiceList {
	selectWidget := widget.NewSelect(options, nil)
	if len(options) > 0 {
		selectWidget.SetSelected(options[0])
	}

	return &ChoiceList{View: selectWidget}
}

func (c *ChoiceList) Selected() string {
	return c.View.Selected
}
