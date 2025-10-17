package routemanager

type StaticRoute struct {
	Destination string `json:"destination"`
	Interface   string `json:"interface"`
	Gateway     string `json:"gateway"`
}

type SystemRoute struct {
	Interface   string
	Destination string
	Gateway     string
	Protocol    string // e.g., "static", "kernel", "dhcp"
	IsStatic    bool   // A flag to easily identify deletable routes
}
