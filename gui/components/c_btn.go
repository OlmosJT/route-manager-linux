package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// CustomButton is now a custom widget that embeds widget.BaseWidget.
type CustomButton struct {
	widget.BaseWidget
	button   *widget.Button
	minWidth float32
}

// customButtonRenderer is the private renderer for our custom widget.
type customButtonRenderer struct {
	btn *CustomButton
}

// Layout defines how the child objects are positioned.
func (r *customButtonRenderer) Layout(size fyne.Size) {
	r.btn.button.Resize(size)
}

// MinSize calculates the minimum size required for the widget.
// This is where we enforce our custom minimum width.
func (r *customButtonRenderer) MinSize() fyne.Size {
	min := r.btn.button.MinSize()
	if min.Width < r.btn.minWidth {
		min.Width = r.btn.minWidth
	}
	return min
}

// Refresh is called when the widget needs to be redrawn.
func (r *customButtonRenderer) Refresh() {
	r.btn.button.Refresh()
}

// Objects returns all the child objects that should be rendered.
func (r *customButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.btn.button}
}

// Destroy is called when the widget is removed.
func (r *customButtonRenderer) Destroy() {}

// CreateRenderer is the constructor for our custom renderer.
func (b *CustomButton) CreateRenderer() fyne.WidgetRenderer {
	return &customButtonRenderer{btn: b}
}

// NewCustomButton is the constructor for our custom CustomButton widget.
func NewCustomButton(label string, onTapped func()) *CustomButton {
	btn := &CustomButton{}
	// This is essential for all custom widgets.
	btn.ExtendBaseWidget(btn)

	btn.button = widget.NewButton(label, onTapped)
	return btn
}

func (b *CustomButton) SetIcon(icon fyne.Resource) {
	b.button.SetIcon(icon)
}

// SetMinWidth allows setting the minimum width for the button.
func (b *CustomButton) SetMinWidth(width float32) {
	b.minWidth = width
}

// Enable makes the button clickable.
func (b *CustomButton) Enable() {
	b.button.Enable()
}

// Disable makes the button unclickable.
func (b *CustomButton) Disable() {
	b.button.Disable()
}
