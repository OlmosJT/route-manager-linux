package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// ValidatorFunc and InputField struct are unchanged
type ValidatorFunc func(string) bool
type InputField struct {
	widget.BaseWidget
	entry               *widget.Entry
	border              *canvas.Rectangle
	validator           ValidatorFunc
	minWidth            float32
	OnValidationChanged func(bool)
}

// Renderer struct is unchanged
type inputFieldRenderer struct {
	field *InputField
}

func (r *inputFieldRenderer) Layout(size fyne.Size) {
	// The border takes the full size of the widget
	r.field.border.Resize(size)

	// Inset the entry by 2 pixels on each side
	padding := float32(2)
	entryPos := fyne.NewPos(padding, padding)
	entrySize := fyne.NewSize(size.Width-padding*2, size.Height-padding*2)
	r.field.entry.Move(entryPos)
	r.field.entry.Resize(entrySize)
}

func (r *inputFieldRenderer) MinSize() fyne.Size {
	min := r.field.entry.MinSize()
	// Add the padding to the total minimum size
	padding := float32(2)
	min = min.Add(fyne.NewSize(padding*2, padding*2))

	if min.Width < r.field.minWidth {
		min.Width = r.field.minWidth
	}
	return min
}

func (r *inputFieldRenderer) Refresh() {
	r.field.border.Refresh()
	r.field.entry.Refresh()
}
func (r *inputFieldRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.field.border, r.field.entry}
}
func (r *inputFieldRenderer) Destroy() {}

// CreateRenderer is unchanged
func (f *InputField) CreateRenderer() fyne.WidgetRenderer {
	return &inputFieldRenderer{field: f}
}

// NewInputField and other public methods are unchanged
func NewInputField(placeholder string, validator ValidatorFunc) *InputField {
	field := &InputField{
		validator: validator,
	}
	field.ExtendBaseWidget(field)
	field.entry = widget.NewEntry()
	field.entry.SetPlaceHolder(placeholder)
	field.border = canvas.NewRectangle(color.Transparent)
	field.border.CornerRadius = 5
	field.entry.OnChanged = func(text string) {
		isValid := field.validator(text)
		if text == "" {
			field.border.FillColor = color.Transparent
		} else if isValid {
			field.border.FillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 40}
		} else {
			field.border.FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 40}
		}
		field.border.Refresh()
		if field.OnValidationChanged != nil {
			field.OnValidationChanged(isValid)
		}
	}
	return field
}
func (f *InputField) SetMinWidth(width float32) {
	f.minWidth = width
}
func (f *InputField) Text() string {
	return f.entry.Text
}
func (f *InputField) SetText(text string) {
	f.entry.SetText(text)
}
