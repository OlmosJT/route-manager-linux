package main

import (
	"encoding/json"
	"os"
)

const routeFile = "route.json"

func SaveRoutes(routes []Route) error {
	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(routeFile, data, 0644)
}

func LoadRoutes() ([]Route, error) {
	if _, err := os.Stat(routeFile); os.IsNotExist(err) {
		return []Route{}, nil
	}
	data, err := os.ReadFile(routeFile)
	if err != nil {
		return nil, err
	}
	var routes []Route
	err = json.Unmarshal(data, &routes)
	return routes, err
}
