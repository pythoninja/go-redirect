package model

type Health struct {
	Status         string `json:"status"`
	Version        string `json:"version"`
	Environment    string `json:"env"`
	DatabaseStatus string `json:"db"`
}
