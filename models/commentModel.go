package models

import (
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id" gorm:"not null" binding:"required"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}
