package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// Define states for our validator
const (
	stateNeutral = iota // An empty field
	stateValid
	stateInvalid
)

// ValidatingEntry is a custom widget that wraps an Entry and validates its content.
type ValidatingEntry struct {
	widget.BaseWidget
	Entry           *widget.Entry
	Validator       func(string) bool
	validationState int // Use an int for three states: neutral, valid, invalid
}

// NewValidatingEntry creates a new instance of our custom entry.
func NewValidatingEntry(validator func(string) bool) *ValidatingEntry {
	e := &ValidatingEntry{
		Validator:       validator,
		validationState: stateNeutral, // Start as neutral
	}
	e.Entry = widget.NewEntry()
	e.Entry.OnChanged = func(s string) {
		e.validate(s)
	}
	e.ExtendBaseWidget(e)
	return e
}

func (e *ValidatingEntry) validate(s string) {
	if e.Validator == nil {
		e.validationState = stateValid
		return
	}

	var newState int
	if s == "" {
		newState = stateNeutral
	} else if e.Validator(s) {
		newState = stateValid
	} else {
		newState = stateInvalid
	}

	if e.validationState != newState {
		e.validationState = newState
		e.Refresh()
	}
}

// CreateRenderer is a part of the Fyne widget API.
func (e *ValidatingEntry) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 2
	border.CornerRadius = theme.Size(theme.SizeNameInputRadius)
	border.Hide() // Start with the border hidden

	return &validatingEntryRenderer{
		entry:  e,
		border: border,
		objects: []fyne.CanvasObject{
			e.Entry,
			border,
		},
	}
}

// --- Renderer for our custom widget ---
type validatingEntryRenderer struct {
	entry   *ValidatingEntry
	border  *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *validatingEntryRenderer) Layout(size fyne.Size) {
	r.entry.Entry.Resize(size)
	r.border.Resize(size)
}

func (r *validatingEntryRenderer) MinSize() fyne.Size {
	return r.entry.Entry.MinSize()
}

func (r *validatingEntryRenderer) Refresh() {
	switch r.entry.validationState {
	case stateNeutral:
		r.border.Hide()
	case stateValid:
		// A nice green color for valid input
		r.border.StrokeColor = &color.NRGBA{R: 0x2e, G: 0x8b, B: 0x57, A: 0xff} // SeaGreen
		r.border.Show()
	case stateInvalid:
		r.border.StrokeColor = theme.Color(theme.ColorRed)
		r.border.Show()
	}
	r.border.Refresh()
	r.entry.Entry.Refresh()
}

func (r *validatingEntryRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *validatingEntryRenderer) Destroy() {}
