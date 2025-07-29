package models

import "time"

type Task struct {
	id int 
	description string
	isCompleted bool
	createdAt time.Time
}

