package main

import (
	"fmt"
	"log"
	"route-manager/gui"
	"route-manager/routemanager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Route Manager")

	// 1. Create the UI components
	header := gui.NewAppHeader()
	quickApply := gui.NewQuickApplyBar()
	routeTable := gui.NewRouteTable() // Create the route table

	// 2. Define the application's core logic

	// Logic for adding a NEW route
	header.OnAdd = func(route routemanager.StaticRoute, save bool) {
		if err := routemanager.Add(route); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if save {
			if err := routemanager.AppendRoute(route); err != nil {
				dialog.ShowError(err, myWindow)
			}
		}
		dialog.ShowInformation("Success", "Successfully applied route", myWindow)
		header.ClearFields()
		// Refresh other components
		quickApply.Refresh()
		routeTable.Refresh()
	}

	// Logic for applying an EXISTING saved route
	quickApply.OnApply = func(route routemanager.StaticRoute) {
		if err := routemanager.Add(route); err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		dialog.ShowInformation("Success", "Successfully re-applied route", myWindow)
		routeTable.Refresh() // Refresh the table to show the new active route
	}

	// Logic for DELETING a route from the table
	routeTable.OnDelete = func(route routemanager.StaticRoute) {
		// Ask for confirmation before deleting
		confirmCallback := func(confirm bool) {
			if !confirm {
				return
			}
			// User confirmed, proceed with deletion
			if err := routemanager.Delete(route); err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			// Also delete it from our saved JSON file, if it exists there
			if err := routemanager.DeleteRoute(route); err != nil {
				log.Printf("INFO: Could not delete route from JSON (it may not have been saved): %v", err)
			}

			successMsg := fmt.Sprintf("Successfully deleted route to %s", route.Destination)
			dialog.ShowInformation("Success", successMsg, myWindow)

			// Refresh components to reflect the change
			quickApply.Refresh()
			routeTable.Refresh()
		}
		dialog.ShowConfirm("Confirm Deletion", "Are you sure you want to delete this route?", confirmCallback, myWindow)
	}

	// 3. Assemble the main layout
	topPanel := container.NewVBox(
		header.View,
		quickApply.View,
	)

	content := container.NewBorder(
		topPanel,
		nil,
		nil,
		nil,
		routeTable, // Use the new RouteTable component
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}
