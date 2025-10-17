package main

import (
	"fmt"
	"image/color"
	"net"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// App holds the core components of our application.
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	model   *Model
	ui      *UI
}

// Model holds the application's data state.
type Model struct {
	allRoutes        []Route
	displayRoutes    []Route
	storedRoutes     []Route
	selectedRowIndex int
	showStaticOnly   bool
}

// UI holds all the widgets for the application view.
type UI struct {
	table          *widget.Table
	destEntry      *ValidatingEntry
	gwEntry        *ValidatingEntry
	ifaceSelect    *widget.Select
	saveCheck      *widget.Check
	addBtn         *widget.Button
	storedSelect   *widget.Select
	applyStoredBtn *widget.Button
	filterCheck    *widget.Check
	deleteBtn      *widget.Button
}

func validateIP(s string) bool {
	return net.ParseIP(s) != nil
}

func validateCIDRorIP(s string) bool {
	if _, _, err := net.ParseCIDR(s); err == nil {
		return true
	}
	return net.ParseIP(s) != nil
}

func main() {
	myApp := NewApp()
	myApp.Run()
}

// NewApp is the constructor for our application.
func NewApp() *App {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(newZebraTheme())
	window := fyneApp.NewWindow("Linux Route Manager")
	window.Resize(fyne.NewSize(850, 600))

	stored, _ := LoadRoutes()
	model := &Model{
		storedRoutes:     stored,
		selectedRowIndex: -1,
	}

	a := &App{
		fyneApp: fyneApp,
		window:  window,
		model:   model,
		ui:      &UI{},
	}

	return a
}

// Run sets up the UI and starts the application loop.
func (a *App) Run() {
	a.window.SetContent(a.buildMainLayout())
	a.refreshAllData()
	a.window.ShowAndRun()
}

func (a *App) buildMainLayout() *fyne.Container {
	topBar := a.buildTopBar()
	storedBar := a.buildStoredBar()
	actionsBar := a.buildActionsBar()
	infoSection := a.buildInfoSection()
	a.ui.table = a.buildRouteTable()

	formContainer := container.NewVBox(topBar, storedBar, actionsBar, infoSection)
	return container.NewBorder(formContainer, nil, nil, nil, a.ui.table)
}

func (a *App) buildInfoSection() *widget.Accordion {
	helpText := `A route tells your computer how to send network traffic.

• Destination & Subnet Mask: This is the target address.
  - For a single computer, use its IP and /32 (e.g., 10.226.98.107/32).
  - For a whole network, use the network address and a mask (e.g., 192.168.1.0/24).
  - The mask (/32, /24) defines the size of the target "neighborhood".

• Gateway: The "exit door" for traffic leaving your local network.
  - This is usually the IP address of your router (e.g., 192.168.1.1).
  - Traffic to the internet or other networks must go through a gateway.

• Interface: The physical hardware used to send the data.
  - This is your network card, like Wi-Fi (e.g., wlp1s0) or Ethernet (e.g., enp2s0).`

	label := widget.NewLabel(helpText)
	label.Wrapping = fyne.TextWrapWord
	accordion := widget.NewAccordion(widget.NewAccordionItem("How to Add a Static Route?", label))
	accordion.Close(0)
	return accordion
}

func (a *App) buildTopBar() *fyne.Container {
	a.ui.destEntry = NewValidatingEntry(validateCIDRorIP)
	a.ui.destEntry.Entry.SetPlaceHolder("Destination (e.g., 8.8.8.8/32)")

	a.ui.gwEntry = NewValidatingEntry(validateIP)
	a.ui.gwEntry.Entry.SetPlaceHolder("Gateway (e.g., 192.168.1.1)")

	a.ui.ifaceSelect = widget.NewSelect([]string{}, nil)
	a.ui.ifaceSelect.PlaceHolder = "Interface"
	a.ui.saveCheck = widget.NewCheck("Save on restart", nil)
	a.ui.addBtn = widget.NewButton("Add Route", a.addRoute)

	ifaces, _ := GetActiveInterfaces()
	a.ui.ifaceSelect.Options = ifaces

	return container.NewGridWithColumns(5, a.ui.destEntry, a.ui.gwEntry, a.ui.ifaceSelect, a.ui.saveCheck, a.ui.addBtn)
}

func (a *App) buildStoredBar() *fyne.Container {
	a.ui.storedSelect = widget.NewSelect([]string{}, nil)
	a.ui.storedSelect.PlaceHolder = "(Select saved route to apply)"
	a.ui.applyStoredBtn = widget.NewButton("Apply Saved", a.applyStoredRoute)
	a.updateStoredRoutesDropdown()
	return container.NewBorder(nil, nil, nil, a.ui.applyStoredBtn, a.ui.storedSelect)
}

func (a *App) buildActionsBar() *fyne.Container {
	a.ui.filterCheck = widget.NewCheck("Show Static Only", func(checked bool) {
		a.model.showStaticOnly = checked
		a.refreshDisplay()
	})
	a.ui.deleteBtn = widget.NewButton("Delete Selected", a.deleteSelectedRoute)
	return container.NewHBox(a.ui.filterCheck, a.ui.deleteBtn)
}

func (a *App) buildRouteTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) { return len(a.model.displayRoutes) + 1, 5 },
		func() fyne.CanvasObject { return container.NewCenter() },
		func(id widget.TableCellID, cell fyne.CanvasObject) { a.updateTableCell(id, cell) },
	)
	t.SetColumnWidth(0, 40)
	t.SetColumnWidth(1, 220)
	t.SetColumnWidth(2, 220)
	t.SetColumnWidth(3, 150)
	t.SetColumnWidth(4, 80)
	return t
}

// --- App Logic & Event Handlers ---
func (a *App) addRoute() {
	if a.ui.gwEntry.validationState == stateInvalid || a.ui.destEntry.validationState == stateInvalid {
		dialog.ShowInformation("Invalid Input", "Please correct the fields with a red border.", a.window)
		return
	}

	destStr := a.ui.destEntry.Entry.Text
	gwStr := a.ui.gwEntry.Entry.Text

	if destStr == "" || gwStr == "" || a.ui.ifaceSelect.Selected == "" {
		dialog.ShowInformation("Error", "All fields must be filled.", a.window)
		return
	}

	r := Route{Destination: destStr, Gateway: gwStr, Interface: a.ui.ifaceSelect.Selected, Saved: a.ui.saveCheck.Checked}
	if err := AddRoute(r); err != nil {
		dialog.ShowInformation("Cancelled", "Root privileges not granted or command failed.", a.window)
		return
	}

	if r.Saved {
		a.model.storedRoutes = append(a.model.storedRoutes, r)
		if err := SaveRoutes(a.model.storedRoutes); err != nil {
			dialog.ShowError(err, a.window)
		}
		// THIS IS THE CRITICAL LINE THAT FIXES THE BUG
		a.updateStoredRoutesDropdown()
	}

	a.ui.destEntry.Entry.SetText("")
	a.ui.gwEntry.Entry.SetText("")
	a.refreshAllData()
}

func (a *App) applyStoredRoute() {
	if a.ui.storedSelect.Selected == "" {
		dialog.ShowInformation("Error", "No stored route selected", a.window)
		return
	}
	parts := strings.SplitN(a.ui.storedSelect.Selected, ":", 2)
	idx, err := strconv.Atoi(parts[0])
	if err != nil || idx < 1 || idx > len(a.model.storedRoutes) {
		dialog.ShowError(fmt.Errorf("invalid selection"), a.window)
		return
	}
	r := a.model.storedRoutes[idx-1]
	if err = AddRoute(r); err != nil {
		dialog.ShowInformation("Cancelled", "Root privileges not granted or command failed.", a.window)
	} else {
		a.refreshAllData()
	}
}

func (a *App) deleteSelectedRoute() {
	if a.model.selectedRowIndex < 0 || a.model.selectedRowIndex >= len(a.model.displayRoutes) {
		dialog.ShowInformation("No Selection", "Please check a route to delete.", a.window)
		return
	}
	routeToDelete := a.model.displayRoutes[a.model.selectedRowIndex]
	if routeToDelete.Destination == "default" {
		dialog.ShowInformation("Action Not Allowed", "Deleting a default route is prevented.", a.window)
		return
	}
	_ = DeleteRoute(routeToDelete)
	a.refreshAllData()
}

// --- Data & UI Refresh Methods ---
func (a *App) refreshAllData() {
	a.model.allRoutes, _ = ListRoutes()
	a.refreshDisplay()
}

func (a *App) refreshDisplay() {
	a.model.displayRoutes = nil
	if !a.model.showStaticOnly {
		a.model.displayRoutes = a.model.allRoutes
	} else {
		for _, r := range a.model.allRoutes {
			if r.Saved {
				a.model.displayRoutes = append(a.model.displayRoutes, r)
			}
		}
	}
	a.model.selectedRowIndex = -1
	a.ui.table.Refresh()
}

func (a *App) updateStoredRoutesDropdown() {
	var opts []string
	for i, r := range a.model.storedRoutes {
		opts = append(opts, fmt.Sprintf("%d: %s via %s", i+1, r.Destination, r.Gateway))
	}
	a.ui.storedSelect.Options = opts
	a.ui.storedSelect.Refresh()
}

func (a *App) updateTableCell(id widget.TableCellID, cell fyne.CanvasObject) {
	a.ui.table.SetRowHeight(id.Row, theme.TextSize()+theme.Padding()*2)
	container := cell.(*fyne.Container)
	rowIndex := id.Row - 1
	if id.Row == 0 {
		var label *widget.Label
		if len(container.Objects) > 0 {
			label, _ = container.Objects[0].(*widget.Label)
		}
		if label == nil {
			label = widget.NewLabel("")
			label.Alignment = fyne.TextAlignCenter
			container.Objects = []fyne.CanvasObject{label}
		}
		headers := []string{"", "Destination", "Gateway", "Interface", "Saved"}
		label.SetText(headers[id.Col])
		return
	}
	if rowIndex >= len(a.model.displayRoutes) {
		return
	}
	if id.Col == 0 {
		var check *widget.Check
		if len(container.Objects) > 0 {
			check, _ = container.Objects[0].(*widget.Check)
		}
		if check == nil {
			check = widget.NewCheck("", nil)
			container.Objects = []fyne.CanvasObject{check}
		}
		check.SetChecked(a.model.selectedRowIndex == rowIndex)
		check.OnChanged = func(checked bool) {
			if checked {
				a.model.selectedRowIndex = rowIndex
			} else {
				a.model.selectedRowIndex = -1
			}
			a.ui.table.Refresh()
		}
	} else {
		route := a.model.displayRoutes[rowIndex]
		var label *widget.Label
		if len(container.Objects) > 0 {
			label, _ = container.Objects[0].(*widget.Label)
		}
		if label == nil {
			label = widget.NewLabel("")
			container.Objects = []fyne.CanvasObject{label}
		}
		var text string
		switch id.Col {
		case 1:
			text = route.Destination
		case 2:
			text = route.Gateway
		case 3:
			text = route.Interface
		case 4:
			if route.Saved {
				text = "Yes"
			} else {
				text = "No"
			}
		}
		label.SetText(text)
	}
}

// --- Custom Theme ---
type zebraTheme struct{ fyne.Theme }

func (z *zebraTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.NRGBA{R: 0x2a, G: 0x2a, B: 0x2e, A: 0xff}
	}
	return z.Theme.Color(name, variant)
}
func newZebraTheme() fyne.Theme { return &zebraTheme{Theme: theme.DefaultTheme()} }
