package main

type Route struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Interface   string `json:"interface"`
	Saved       bool   `json:"saved"` // true = reapply on restart
}
