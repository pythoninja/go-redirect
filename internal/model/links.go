package model

import "time"

type Link struct {
	Id        int       `json:"id"`
	ShortLink string    `json:"short_link"`
	LongLink  string    `json:"long_link"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}
