package model

import "time"

type Link struct {
	ID        int       `json:"id"`
	Alias     string    `json:"alias"`
	URL       string    `json:"url"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}
