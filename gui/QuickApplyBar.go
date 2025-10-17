package gui

import (
	"fmt"
	"log"
	"route-manager/gui/components"
	"route-manager/routemanager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// QuickApplyBar is a component for quickly re-applying a saved route.
type QuickApplyBar struct {
	View fyne.CanvasObject

	OnApply  func(route routemanager.StaticRoute)
	OnDelete func(route routemanager.StaticRoute)

	// Internal references
	routes       []routemanager.StaticRoute
	dropdown     *components.ChoiceList
	applyButton  *components.CustomButton
	deleteButton *components.CustomButton
}

// NewQuickApplyBar creates a new instance of the component.
func NewQuickApplyBar() *QuickApplyBar {
	bar := &QuickApplyBar{}

	bar.applyButton = components.NewCustomButton("Apply Selected Route", func() {
		if bar.OnApply != nil {
			selectedText := bar.dropdown.Selected()
			for _, r := range bar.routes {
				routeStr := formatRoute(r)
				if routeStr == selectedText {
					bar.OnApply(r)
					break
				}
			}
		}
	})
	bar.applyButton.SetMinWidth(180.0)

	bar.deleteButton = components.NewCustomButton("", func() {
		if bar.OnDelete != nil {
			selectedText := bar.dropdown.Selected()
			for _, r := range bar.routes {
				if formatRoute(r) == selectedText {
					bar.OnDelete(r)
					break
				}
			}
		}
	})
	bar.deleteButton.SetIcon(theme.DeleteIcon())

	bar.dropdown = components.NewChoiceList([]string{})

	buttonGroup := container.NewHBox(bar.applyButton, bar.deleteButton)

	bar.View = container.New(NewProportionalLayout(1, 5),
		bar.dropdown.View,
		buttonGroup,
	)

	bar.Refresh() // Load initial data
	return bar
}

// Refresh reloads the route data and handles the empty state.
func (b *QuickApplyBar) Refresh() {
	routes, err := routemanager.LoadRoutes()
	if err != nil {
		log.Printf("ERROR: Failed to load routes: %v", err)
		return
	}

	if len(routes) > 10 {
		routes = routes[len(routes)-10:]
	}
	b.routes = routes

	var options []string
	if len(b.routes) == 0 {
		options = []string{"[No saved routes]"}
		b.dropdown.View.Disable()
		b.applyButton.Disable()
		b.deleteButton.Disable() // Disable delete button when list is empty
	} else {
		for _, r := range b.routes {
			options = append(options, formatRoute(r))
		}
		b.dropdown.View.Enable()
		b.applyButton.Enable()
		b.deleteButton.Enable() // Enable delete button when list has items
	}

	b.dropdown.View.Options = options
	if len(options) > 0 {
		b.dropdown.View.SetSelected(options[0])
	} else {
		b.dropdown.View.ClearSelected()
	}
	b.dropdown.View.Refresh()
}

// formatRoute is a helper to create a consistent display string for a route.
func formatRoute(r routemanager.StaticRoute) string {
	return fmt.Sprintf("%s via %s (dev %s)", r.Destination, r.Gateway, r.Interface)
}
