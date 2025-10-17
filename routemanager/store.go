package routemanager

import (
	"encoding/json"
	"os"
)

const routesFile = "routes.json"

// SaveRoutes writes a slice of StaticRoute structs to the JSON file.
// This is the low-level function that overwrites the file.
func SaveRoutes(routes []StaticRoute) error {
	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return err
	}
	// 0644 are standard file permissions (read/write for owner, read-only for others).
	return os.WriteFile(routesFile, data, 0644)
}

// LoadRoutes reads all routes from the JSON file.
func LoadRoutes() ([]StaticRoute, error) {
	data, err := os.ReadFile(routesFile)
	if err != nil {
		// If the file simply doesn't exist, that's not a critical error.
		// It just means no routes have been saved yet. We return an empty list.
		if os.IsNotExist(err) {
			return []StaticRoute{}, nil
		}
		return nil, err // For any other error (e.g., permissions), return it.
	}

	var routes []StaticRoute
	if err = json.Unmarshal(data, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}

// AppendRoute adds a single new route to the routes.json file.
// It uses the safe "Read-Modify-Write" pattern.
func AppendRoute(newRoute StaticRoute) error {
	// 1. Read
	routes, err := LoadRoutes()
	if err != nil {
		return err
	}

	// 2. Modify
	routes = append(routes, newRoute)

	// 3. Write
	return SaveRoutes(routes)
}

// DeleteRoute removes a specific route from the routes.json file.
// It also uses the "Read-Modify-Write" pattern.
func DeleteRoute(routeToDelete StaticRoute) error {
	// 1. Read
	routes, err := LoadRoutes()
	if err != nil {
		return err
	}

	// 2. Modify: Create a new slice containing only the routes we want to keep.
	var updatedRoutes []StaticRoute
	for _, route := range routes {
		// This checks if the current route is the one we want to delete.
		if route.Interface == routeToDelete.Interface &&
			route.Destination == routeToDelete.Destination &&
			route.Gateway == routeToDelete.Gateway {
			continue // Skip adding it to the new slice
		}
		updatedRoutes = append(updatedRoutes, route)
	}

	// 3. Write
	return SaveRoutes(updatedRoutes)
}
