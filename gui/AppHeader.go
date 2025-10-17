// file: gui/AppHeader.go
package gui

import (
	"route-manager/gui/components"
	"route-manager/routemanager"
	"route-manager/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// AppHeader struct is unchanged.
type AppHeader struct {
	View fyne.CanvasObject

	OnAdd func(route routemanager.StaticRoute, save bool)

	destInput    *components.InputField
	gatewayInput *components.InputField
	addButton    *components.CustomButton
}

// NewAppHeader creates a new header component.
func NewAppHeader() *AppHeader {
	header := &AppHeader{}

	header.destInput = components.NewInputField("Destination (e.g. 10.226.98.107/32)", validators.ValidateCIDR)
	header.gatewayInput = components.NewInputField("Gateway (e.g. 10.226.35.1)", validators.ValidateIP)

	header.destInput.SetMinWidth(160.0)
	header.gatewayInput.SetMinWidth(160.0)

	interfaceNames := routemanager.GetInterfaceNames()
	interfaceChoice := components.NewChoiceList(interfaceNames)
	saveCheckbox := components.NewCustomCheckbox("Save")

	header.addButton = components.NewCustomButton("Add Route", func() {
		if header.OnAdd != nil {
			route := routemanager.StaticRoute{
				Destination: header.destInput.Text(),
				Gateway:     header.gatewayInput.Text(),
				Interface:   interfaceChoice.Selected(),
			}
			header.OnAdd(route, saveCheckbox.IsChecked())
		}
	})
	header.addButton.SetMinWidth(120.0)
	header.addButton.Disable()

	var isDestValid, isGatewayValid bool
	checkOverallValidation := func() {
		if isDestValid && isGatewayValid {
			header.addButton.Enable()
		} else {
			header.addButton.Disable()
		}
	}
	header.destInput.OnValidationChanged = func(isValid bool) {
		isDestValid = isValid
		checkOverallValidation()
	}
	header.gatewayInput.OnValidationChanged = func(isValid bool) {
		isGatewayValid = isValid
		checkOverallValidation()
	}

	header.View = container.New(NewProportionalLayout(2, 5),
		header.destInput,
		header.gatewayInput,
		interfaceChoice.View,
		saveCheckbox.View,
		header.addButton,
	)

	return header
}

// ClearFields method is unchanged.
func (h *AppHeader) ClearFields() {
	h.destInput.SetText("")
	h.gatewayInput.SetText("")
}
