package models

import (
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	id          int
	description string
	isCompleted bool
	createdAt   time.Time
}

func (t *Task) ToCSVRecord() []string {
	return []string{
		strconv.Itoa(t.id),
		t.description,
		strconv.FormatBool(t.isCompleted),
		t.createdAt.Format(time.RFC3339),
	}
}

func FromCSVRecord(record []string) (*Task, error) {
	if len(record) != 4 {
		return nil, fmt.Errorf("invalid record length")
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	completed, err := strconv.ParseBool(record[2])
	if err != nil {
		return nil, fmt.Errorf("invalid completed status: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, record[3])
	if err != nil {
		return nil, fmt.Errorf("invalid created_at: %w", err)
	}

	return &Task{
		id:          id,
		description: record[1],
		isCompleted: completed,
		createdAt:   createdAt,
	}, nil

}
