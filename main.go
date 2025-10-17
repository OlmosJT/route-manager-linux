package main

import (
	"fmt"
	"log"
	"route-manager/gui"
	"route-manager/routemanager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog" // ⭐️ IMPORT the dialog package
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Route Manager")

	// 1. Create the UI components
	header := gui.NewAppHeader()
	quickApply := gui.NewQuickApplyBar()
	routeListPlaceholder := widget.NewLabel("List of active/saved routes will appear here...")

	// 2. Define the application's core logic

	// Logic for the TOP header (adding a NEW route)
	header.OnAdd = func(route routemanager.StaticRoute, save bool) {
		log.Printf("Adding route: %+v (Save: %v)", route, save)

		// a. Apply the route to the OS (requires sudo)
		if err := routemanager.Add(route); err != nil {
			log.Printf("ERROR: Failed to apply route: %v", err)
			// ⭐️ Show an error dialog
			dialog.ShowError(err, myWindow)
			return
		}

		// b. If the user checked "Save", append it to the JSON file
		if save {
			if err := routemanager.AppendRoute(route); err != nil {
				log.Printf("ERROR: Failed to save route: %v", err)
				// ⭐️ Show an error dialog
				dialog.ShowError(err, myWindow)
			} else {
				// c. If save was successful, refresh the quick apply bar
				quickApply.Refresh()
			}
		}

		// ⭐️ Show a success dialog
		successMsg := fmt.Sprintf("Successfully applied route to %s", route.Destination)
		dialog.ShowInformation("Success", successMsg, myWindow)

		header.ClearFields()
	}

	// Logic for the QUICK APPLY bar (applying an EXISTING saved route)
	quickApply.OnApply = func(route routemanager.StaticRoute) {
		log.Printf("Quick applying route: %+v", route)

		if err := routemanager.Add(route); err != nil {
			log.Printf("ERROR: Failed to quick-apply route: %v", err)
			// ⭐️ Show an error dialog
			dialog.ShowError(err, myWindow)
			return // Stop execution if there was an error
		}

		// ⭐️ Show a success dialog
		successMsg := fmt.Sprintf("Successfully re-applied route to %s", route.Destination)
		dialog.ShowInformation("Success", successMsg, myWindow)
	}

	// 3. Assemble the main layout
	topPanel := container.NewVBox(
		header.View,
		quickApply.View,
	)

	content := container.NewBorder(
		topPanel, // Place the grouped panel at the top
		nil,
		nil,
		nil,
		routeListPlaceholder,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}
