package routemanager

import (
	"log"

	"github.com/vishvananda/netlink"
)

func ListSystemRoutes() []SystemRoute {
	var systemRoutes []SystemRoute

	// Passing nil for the link gets all routes on the system.
	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		log.Printf("ERROR: Could not list system routes: %v", err)
		return systemRoutes
	}

	for _, r := range routes {
		// Skip routes that don't have a destination
		if r.Dst == nil {
			continue
		}

		gateway := ""
		if r.Gw != nil {
			gateway = r.Gw.String()
		}

		link, err := netlink.LinkByIndex(r.LinkIndex)
		if err != nil {
			log.Printf("WARN: Could not find link for index %d: %v", r.LinkIndex, err)
			continue
		}

		// The type of r.Protocol is netlink.RouteProtocol, so we pass it directly
		protocol, isStatic := interpretProtocol(r.Protocol)

		systemRoutes = append(systemRoutes, SystemRoute{
			Interface:   link.Attrs().Name,
			Destination: r.Dst.String(),
			Gateway:     gateway,
			Protocol:    protocol,
			IsStatic:    isStatic,
		})
	}

	return systemRoutes
}

func interpretProtocol(p netlink.RouteProtocol) (string, bool) {
	// These values correspond to the constants defined in Linux's networking headers.
	const (
		protoKernel = 2
		protoBoot   = 3
		protoStatic = 4
		protoDhcp   = 16
	)

	switch p {
	case protoKernel:
		return "kernel", false // Not deletable
	case protoBoot:
		return "static", true // Deletable
	case protoStatic:
		return "static", true // Deletable
	case protoDhcp:
		return "dhcp", false // Not deletable
	default:
		// Any other protocol is likely user-added (e.g., from a routing daemon)
		// and we'll consider it "static" for our purposes.
		return "static", true
	}
}
