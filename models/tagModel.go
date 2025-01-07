package models

import (
	"time"
)

type Tag struct {
	ID        int       `json:"id"`
	User_id   int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Task_Tags struct {
	User_id int `json:"user_id"`
	Tag_id  int `json:"tag_id"`
	Task_id int `json:"task_id"`
}
