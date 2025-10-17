package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// HelpSection is a component that provides instructional text in a collapsible accordion.
type HelpSection struct {
	View fyne.CanvasObject
}

// NewHelpSection creates a new help section component.
func NewHelpSection() *HelpSection {
	// 1. Explanation of Subnet Masks
	subnetIntro := widget.NewLabel("A subnet mask defines the size of a network. The CIDR notation (like /24) is a shorthand for the full address.")
	subnetIntro.Wrapping = fyne.TextWrapWord

	subnetTable := widget.NewForm(
		widget.NewFormItem("Single Host (/32):", widget.NewLabel("255.255.255.255")),
		widget.NewFormItem("Common LAN (/24):", widget.NewLabel("255.255.255.0")),
		widget.NewFormItem("Larger Network (/16):", widget.NewLabel("255.255.0.0")),
	)

	// 2. Explanation of Destination Address
	destIntro := widget.NewLabel("The Destination is usually a network address (ending in .0), not a specific device. Use /32 to route to a single host.")
	destIntro.Wrapping = fyne.TextWrapWord

	destExamples := widget.NewLabel(
		"Good: 10.226.98.0/24 (Routes traffic for the entire 10.226.98.x network)\n" +
			"Specific: 10.226.98.107/32 (Routes traffic for only the 10.226.98.107 host)",
	)

	// 3. Explanation of how to find an IP
	findIpIntro := widget.NewLabel("If you only know a domain name (e.g., google.com), you can find its IP address using a terminal command:")
	findIpIntro.Wrapping = fyne.TextWrapWord

	findIpCommands := widget.NewLabel("ping google.com\nnslookup google.com")

	// 4. Assemble all the help content in a vertical box
	helpContent := container.NewVBox(
		subnetIntro,
		subnetTable,
		widget.NewSeparator(),
		destIntro,
		destExamples,
		widget.NewSeparator(),
		findIpIntro,
		findIpCommands,
	)

	// 5. Create the Accordion item
	accordionItem := widget.NewAccordionItem(
		"Need Help? Click to Expand Instructions",
		helpContent,
	)

	return &HelpSection{
		View: widget.NewAccordion(accordionItem),
	}
}
