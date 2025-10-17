package gui

import "fyne.io/fyne/v2"

// ProportionalLayout now includes a Gap field.
type ProportionalLayout struct {
	ExpandingElements int
	Gap               float32
}

// NewProportionalLayout now accepts a gap size.
func NewProportionalLayout(expanding int, gap float32) fyne.Layout {
	return &ProportionalLayout{ExpandingElements: expanding, Gap: gap}
}

// Layout function updated to include gaps.
func (p *ProportionalLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	totalMinWidth := float32(0)
	for _, obj := range objects {
		totalMinWidth += obj.MinSize().Width
	}

	// Account for the total width of all gaps
	totalGap := p.Gap * float32(len(objects)-1)
	if totalGap < 0 {
		totalGap = 0
	}

	extraSpace := size.Width - totalMinWidth - totalGap
	if extraSpace < 0 {
		extraSpace = 0
	}

	var spacePerExpandingElement float32
	if p.ExpandingElements > 0 {
		spacePerExpandingElement = extraSpace / float32(p.ExpandingElements)
	}

	x := float32(0)
	for i, obj := range objects {
		minWidth := obj.MinSize().Width
		var newWidth float32

		if i < p.ExpandingElements {
			newWidth = minWidth + spacePerExpandingElement
		} else {
			newWidth = minWidth
		}

		obj.Move(fyne.NewPos(x, 0))
		obj.Resize(fyne.NewSize(newWidth, size.Height))
		// Move x past the current object AND the gap
		x += newWidth + p.Gap
	}
}

// MinSize function updated to include gaps.
func (p *ProportionalLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minWidth := float32(0)
	maxHeight := float32(0)
	for _, obj := range objects {
		minWidth += obj.MinSize().Width
		if obj.MinSize().Height > maxHeight {
			maxHeight = obj.MinSize().Height
		}
	}

	// Add the total width of all gaps to the container's minimum size
	if len(objects) > 1 {
		totalGap := p.Gap * float32(len(objects)-1)
		minWidth += totalGap
	}

	return fyne.NewSize(minWidth, maxHeight)
}
