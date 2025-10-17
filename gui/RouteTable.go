package gui

import (
	"image/color"
	"route-manager/routemanager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type RouteTable struct {
	widget.BaseWidget
	OnDelete func(route routemanager.StaticRoute)

	table          *widget.Table
	deleteButton   *widget.Button
	filterCheck    *widget.Check
	allRoutes      []routemanager.SystemRoute
	filteredRoutes []routemanager.SystemRoute
	selectedID     int // Index in filteredRoutes, adjusted for header
}

func NewRouteTable() *RouteTable {
	t := &RouteTable{selectedID: -1}
	t.ExtendBaseWidget(t)
	return t
}

func (t *RouteTable) CreateRenderer() fyne.WidgetRenderer {
	// 1. CREATE CONTROLS
	t.filterCheck = widget.NewCheck("Show only static routes", func(checked bool) {
		t.applyFilter(checked)
	})

	t.deleteButton = widget.NewButtonWithIcon("Delete Selected Route", theme.DeleteIcon(), func() {
		if t.selectedID >= 0 && t.OnDelete != nil {
			routeToDelete := t.filteredRoutes[t.selectedID]
			staticRoute := routemanager.StaticRoute{
				Interface:   routeToDelete.Interface,
				Destination: routeToDelete.Destination,
				Gateway:     routeToDelete.Gateway,
			}
			t.OnDelete(staticRoute)
		}
	})
	t.deleteButton.Disable()

	// 2. BUILD THE TABLE WITH AN INTEGRATED HEADER
	headers := []string{"Destination", "Gateway", "Interface", "Protocol"}
	t.table = &widget.Table{
		Length: func() (int, int) {
			// Add 1 to the row count for our header row
			return len(t.filteredRoutes) + 1, len(headers)
		},
		CreateCell: func() fyne.CanvasObject {
			return container.NewStack(
				canvas.NewRectangle(color.Transparent),
				widget.NewLabel(""),
			)
		},
		UpdateCell: func(id widget.TableCellID, cell fyne.CanvasObject) {
			stack := cell.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)

			if id.Row == 0 { // This is the HEADER row
				label.SetText(headers[id.Col])
				label.TextStyle.Bold = true
				bg.FillColor = color.Transparent
			} else { // These are DATA rows
				route := t.filteredRoutes[id.Row-1] // Adjust index for header
				label.TextStyle.Bold = false
				var text string
				switch id.Col {
				case 0:
					text = route.Destination
				case 1:
					text = route.Gateway
				case 2:
					text = route.Interface
				case 3:
					text = route.Protocol
				}
				label.SetText(text)

				// Visual selection logic
				if (id.Row - 1) == t.selectedID {
					bg.FillColor = theme.FocusColor()
				} else {
					bg.FillColor = color.Transparent
				}
			}
			bg.Refresh()
			label.Refresh()
		},
		OnSelected: func(id widget.TableCellID) {
			if id.Row == 0 { // Don't allow selecting the header
				t.table.UnselectAll()
				t.selectedID = -1
				t.deleteButton.Disable()
				return
			}
			t.selectedID = id.Row - 1 // Adjust index for header
			if t.filteredRoutes[t.selectedID].IsStatic {
				t.deleteButton.Enable()
			} else {
				t.deleteButton.Disable()
			}
			t.table.Refresh()
		},
	}
	// These widths now apply to both the "header" and data cells perfectly
	t.table.SetColumnWidth(0, 250)
	t.table.SetColumnWidth(1, 200)
	t.table.SetColumnWidth(2, 150)
	t.table.SetColumnWidth(3, 100)

	// 3. ASSEMBLE THE FINAL LAYOUT
	controlBar := container.NewHBox(t.deleteButton, t.filterCheck)

	// Use a VBox to stack the controls above the table
	content := container.NewBorder(controlBar, nil, nil, nil, t.table)

	t.Refresh()
	return widget.NewSimpleRenderer(content)
}

// Refresh and applyFilter methods are adjusted to handle the new selection logic
func (t *RouteTable) Refresh() {
	t.allRoutes = routemanager.ListSystemRoutes()
	t.table.UnselectAll()
	t.selectedID = -1
	t.deleteButton.Disable()
	t.applyFilter(t.filterCheck.Checked)
}

func (t *RouteTable) applyFilter(onlyStatic bool) {
	t.table.UnselectAll() // Clear selection when filtering
	t.selectedID = -1
	t.deleteButton.Disable()

	if !onlyStatic {
		t.filteredRoutes = t.allRoutes
	} else {
		var filtered []routemanager.SystemRoute
		for _, r := range t.allRoutes {
			if r.IsStatic {
				filtered = append(filtered, r)
			}
		}
		t.filteredRoutes = filtered
	}
	t.table.Refresh()
}
