package routemanager

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// Add applies a static route to the system's routing table.
// It uses RouteReplace which acts as an "upsert" (update or insert),
// making it safer than RouteAdd as it won't fail if the route already exists.
func Add(route StaticRoute) error {
	link, err := netlink.LinkByName(route.Interface)
	if err != nil {
		return fmt.Errorf("interface %s not found: %w", route.Interface, err)
	}

	_, dst, err := net.ParseCIDR(route.Destination)
	if err != nil {
		return fmt.Errorf("invalid destination CIDR %s: %w", route.Destination, err)
	}

	gw := net.ParseIP(route.Gateway)
	if gw == nil {
		return fmt.Errorf("invalid gateway IP %s", route.Gateway)
	}

	routeObj := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Gw:        gw,
	}

	return netlink.RouteReplace(routeObj)
}

// Delete removes a static route from the system's routing table.
func Delete(route StaticRoute) error {
	link, err := netlink.LinkByName(route.Interface)
	if err != nil {
		return fmt.Errorf("interface %s not found: %w", route.Interface, err)
	}

	_, dst, err := net.ParseCIDR(route.Destination)
	if err != nil {
		return fmt.Errorf("invalid destination CIDR %s: %w", route.Destination, err)
	}

	// ⭐️ Safety check to prevent deleting the default route.
	// The IsUnspecified method checks for 0.0.0.0 (IPv4) or :: (IPv6).
	if dst.IP.IsUnspecified() {
		return errors.New("deleting the default route is not allowed")
	}

	gw := net.ParseIP(route.Gateway)
	if gw == nil {
		return fmt.Errorf("invalid gateway IP %s", route.Gateway)
	}

	routeObj := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Gw:        gw,
	}

	return netlink.RouteDel(routeObj)
}
