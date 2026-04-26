package todo

import (
	"cli-todo-app/internal/store"
	"log/slog"
)

const StorageType = "sqlite"

func CreateTask(name string, description string, state string) {
	tasks, err := store.ReadTasks(StorageType)
	if err != nil {
		slog.Error("failed to read tasks", "storage", StorageType, "error", err)
		return
	}

	nextID := len(tasks) + 1

	err = store.WriteTask(StorageType, store.Task{
		ID:          nextID,
		Name:        name,
		Description: description,
		State:       state,
	})
	if err != nil {
		slog.Error("failed to write task", "storage", StorageType, "error", err)
		return
	}

	slog.Info("Created task", "name", name, "id", nextID)
}

func GetTasks() {

}

func UpdateTask() {

}
