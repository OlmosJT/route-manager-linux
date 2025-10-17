package validators

import "net"

// ValidateCIDR checks if a string is a valid CIDR notation (e.g., "192.168.1.0/24").
// It is exported because it starts with a capital 'V'.
func ValidateCIDR(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

// ValidateIP checks if a string is a valid IP address (e.g., "192.168.1.1").
// It is also exported.
func ValidateIP(s string) bool {
	return net.ParseIP(s) != nil
}
