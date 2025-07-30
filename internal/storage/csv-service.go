package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"todo/models"
)

type CSVStorage struct {
	filename string
}

func NewCSVStorage(filename string) *CSVStorage {
	return &CSVStorage{
		filename: filename,
	}
}

func (s *CSVStorage) AddTask(description string) (*models.Task, error) {
	nextID, err := s.getNextID()
	if err != nil {
		return nil, fmt.Errorf("failed to get next id: %w", err)
	}

	task := &models.Task{
		ID:          nextID,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	err = s.writeTaskToFile(task)
	if err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	return task, nil
}

func (s *CSVStorage) UpdateTask(id int, newDescription *string, completed *bool) (*models.Task, error) {
	tasks, err := s.ListTask()
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	// Find and update the task
	var updatedTask *models.Task
	found := false
	for _, task := range tasks {
		if task.ID == id {
			if newDescription != nil {
				task.Description = *newDescription
			}
			if completed != nil {
				task.Completed = *completed
			}
			updatedTask = task
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	// Rewrite the entire file with updated tasks
	err = s.rewriteAllTasks(tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to save updated task: %w", err)
	}

	return updatedTask, nil
}

func (s *CSVStorage) DeleteTask(id int) error {
	tasks, err := s.ListTask()
	if err != nil {
		return fmt.Errorf("failed to read tasks: %w", err)
	}

	var filteredTasks []*models.Task
	found := false
	for _, task := range tasks {
		if task.ID != id {
			filteredTasks = append(filteredTasks, task)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("task with ID  %d not found", id)
	}
	return s.rewriteAllTasks(filteredTasks)
}

func (s *CSVStorage) ListTask() ([]*models.Task, error) {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return []*models.Task{}, nil
	}
	file, err := os.Open(s.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	tasks := make([]*models.Task, 0, len(records))
	for _, record := range records {
		task, err := models.FromCSVRecord(record)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *CSVStorage) rewriteAllTasks(tasks []*models.Task) error {
	// Create/truncate the file
	file, err := os.Create(s.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write all tasks
	for _, task := range tasks {
		record := task.ToCSVRecord()
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func (s *CSVStorage) getNextID() (int, error) {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return 1, nil
	}

	file, err := os.Open(s.filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV %w", err)
	}

	if len(records) == 0 {
		return 1, nil
	}

	maxID := 0
	for _, record := range records {
		if len(record) > 0 {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			if id > maxID {
				maxID = id
			}
		}

	}
	return maxID + 1, nil
}

func (s *CSVStorage) writeTaskToFile(task *models.Task) error {
	file, err := os.OpenFile(s.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := task.ToCSVRecord()
	err = writer.Write(record)
	if err != nil {
		return fmt.Errorf("faild to write record %w", err)
	}

	return nil

}
