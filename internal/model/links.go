package model

import "time"

type Link struct {
	Id        int       `json:"id"`
	Alias     string    `json:"alias"`
	Url       string    `json:"url"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}
