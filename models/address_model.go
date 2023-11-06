package models

type Address struct {
	Address string   `json:"address" validate:"required"`
	Labels  []string `json:"labels,omitempty"`
}
