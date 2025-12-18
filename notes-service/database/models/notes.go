package models

import "time"

type Notes struct {
	Id           int        `json:"id"`
	UserPublicID string     `json:"user_public_id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Category     string     `json:"category"`
	Color        string     `json:"color"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
