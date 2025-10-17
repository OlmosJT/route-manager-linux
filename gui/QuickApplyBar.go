package gui

import (
	"fmt"
	"log"
	"route-manager/gui/components"
	"route-manager/routemanager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// QuickApplyBar is a component for quickly re-applying a saved route.
type QuickApplyBar struct {
	View fyne.CanvasObject

	OnApply func(route routemanager.StaticRoute)

	// Internal references
	routes      []routemanager.StaticRoute
	dropdown    *components.ChoiceList
	applyButton *components.CustomButton // Add a reference to the button
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

	bar.dropdown = components.NewChoiceList([]string{})

	// ⭐️ FIX: Use our ProportionalLayout to get a consistent 5-pixel gap.
	// We make the dropdown the single expanding element.
	bar.View = container.New(NewProportionalLayout(1, 5),
		bar.dropdown.View,
		bar.applyButton,
	)

	bar.Refresh() // Load initial data
	return bar
}

// Refresh reloads the route data and handles the empty state.
func (b *QuickApplyBar) Refresh() {
	routes, err := routemanager.LoadRoutes()
	if err != nil {
		log.Printf("ERROR: Failed to load routes for quick apply bar: %v", err)
		return
	}

	if len(routes) > 10 {
		routes = routes[len(routes)-10:]
	}
	b.routes = routes

	var options []string
	// ⭐️ FIX: Handle the case where there are no saved routes.
	if len(b.routes) == 0 {
		options = []string{"[No saved routes]"}
		b.dropdown.View.Disable()
		b.applyButton.Disable()
	} else {
		for _, r := range b.routes {
			options = append(options, formatRoute(r))
		}
		b.dropdown.View.Enable()
		b.applyButton.Enable()
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
