package model

type Health struct {
	Status         string     `json:"status"`
	Environment    string     `json:"env"`
	DatabaseStatus string     `json:"db"`
	Build          *BuildInfo `json:"build"`
}
