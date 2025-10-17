package main

import (
	"log"
	"os/exec"
	"strings"
)

// GetActiveInterfaces returns all non-loopback, non-bridge interfaces
func GetActiveInterfaces() ([]string, error) {
	cmd := exec.Command("ip", "-o", "link", "show")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ifaces []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) > 1 {
			name := strings.Fields(parts[1])[0]
			if name != "lo" && !strings.HasPrefix(name, "br-") && name != "docker0" {
				ifaces = append(ifaces, name)
			}
		}
	}
	return ifaces, nil
}

// ListRoutes parses the output of "ip route show"
func ListRoutes() ([]Route, error) {
	cmd := exec.Command("ip", "route", "show")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var routes []Route
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		r := Route{}
		r.Destination = parts[0]
		for i, p := range parts {
			if p == "via" && i+1 < len(parts) {
				r.Gateway = parts[i+1]
			}
			if p == "dev" && i+1 < len(parts) {
				r.Interface = parts[i+1]
			}
		}
		if lineContainsStatic(line) {
			r.Saved = true
		}
		routes = append(routes, r)
	}
	return routes, nil
}

func lineContainsStatic(line string) bool {
	return strings.Contains(line, "proto static") || strings.Contains(line, "via")
}

// AddRoute uses pkexec to add a route with root privileges
func AddRoute(r Route) error {
	cmd := exec.Command("pkexec", "ip", "route", "add", r.Destination, "via", r.Gateway, "dev", r.Interface)
	return cmd.Run()
}

// DeleteRoute uses pkexec to delete a route with root privileges
func DeleteRoute(r Route) error {
	cmd := exec.Command("pkexec", "ip", "route", "del", r.Destination, "via", r.Gateway, "dev", r.Interface)
	return cmd.Run()
}
