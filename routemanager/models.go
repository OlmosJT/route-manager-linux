package routemanager

type StaticRoute struct {
	Destination string `json:"destination"`
	Interface   string `json:"interface"`
	Gateway     string `json:"gateway"`
}
