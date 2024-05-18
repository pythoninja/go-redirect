package model

type Links struct {
	Id           int    `json:"id"`
	ShortLink    string `json:"short_link"`
	LongLink     string `json:"long_link"`
	ClickCounter int    `json:"clicks"`
}
