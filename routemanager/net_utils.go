package routemanager

import (
	"log"
	"net"
)

func GetInterfaceNames() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
		return []string{}
	}

	var names []string
	for _, i := range interfaces {
		// Filter out loopback interfaces (like 'lo') and interfaces that are down.
		isUp := i.Flags&net.FlagUp != 0
		isLoopback := i.Flags&net.FlagLoopback != 0

		if isUp && !isLoopback {
			names = append(names, i.Name)
		}
	}
	return names
}
