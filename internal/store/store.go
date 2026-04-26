package store

import (
	"fmt"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
}

type TaskFilter struct {
	Search string
	State  string
	ID     *int
}

func WriteTask(storageType string, task Task) error {
	if storageType == "sqlite" {
		return writeTaskToSQLite(task)
	}

	return fmt.Errorf("unsupported storage type: %s", storageType)
}

func ReadTasks(storageType string, filter TaskFilter) ([]Task, error) {
	if storageType == "sqlite" {
		return readTasksFromSQLite(filter)
	}

	return nil, fmt.Errorf("unsupported storage type: %s", storageType)
}

func ReadTask(storageType string, id int) (*Task, error) {
	if storageType == "sqlite" {
		return readTaskFromSQLite(id)
	}
	return nil, fmt.Errorf("unsupported storage type: %s", storageType)
}

func UpdateTask(storageType string, task Task) error {
	if storageType == "sqlite" {
		return updateTaskInSQLite(task)
	}
	return fmt.Errorf("unsupported storage type: %s", storageType)
}
