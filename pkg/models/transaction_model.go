package models

type Status string

const (
	Success Status = "success"
	Revert  Status = "revert"
	Pending Status = "pending"
)

type Transaction struct {
	Hash    string `json:"hash"`
	Chainid uint   `json:"chainid" validate:"required"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Status  Status `json:"status"`
}
