package models

import (
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Description string
	Completed bool
	CreatedAt   time.Time
}

func (t *Task) ToCSVRecord() []string {
	return []string{
		strconv.Itoa(t.ID),
		t.Description,
		strconv.FormatBool(t.Completed),
		t.CreatedAt.Format(time.RFC3339),
	}
}

func FromCSVRecord(record []string) (*Task, error) {
	if len(record) != 4 {
		return nil, fmt.Errorf("invalid record length")
	}

	ID, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	completed, err := strconv.ParseBool(record[2])
	if err != nil {
		return nil, fmt.Errorf("invalid completed status: %w", err)
	}

	CreatedAt, err := time.Parse(time.RFC3339, record[3])
	if err != nil {
		return nil, fmt.Errorf("invalid created_at: %w", err)
	}

	return &Task{
		ID:          ID,
		Description: record[1],
		Completed: completed,
		CreatedAt:   CreatedAt,
	}, nil

}
