package models

import (
	"time"
)

// Comment represents a comment on a video
type Comment struct {
	ID       uint64    `json:"id"`
	UserID   uint64    `json:"user_id"`
	VideoID  uint64    `json:"video_id"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// TableName sets the name of the table in the database for the Comment model
func (Comment) TableName() string {
	return "comments"
}
