package models

type Label struct {
	Label     string   `json:"label"`
	Addresses []string `json:"address"`
}
