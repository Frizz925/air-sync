package models

type Session struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt int64     `json:"created_at"`
}
