package todo

import (
	"cli-todo-app/internal/store"
	"fmt"
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
	tasks, err := store.ReadTasks(StorageType)
	if err != nil {
		slog.Error("failed to read tasks", "storage", StorageType, "error", err)
		return
	}

	if len(tasks) == 0 {
		slog.Info("No tasks found.")
		return
	}

	fmt.Printf("%-4s  %-20s  %-10s  %s\n", "ID", "Name", "State", "Description")
	fmt.Println("----  --------------------  ----------  -----------")
	for _, task := range tasks {
		fmt.Printf("%-4d  %-20s  %-10s  %s\n", task.ID, task.Name, task.State, task.Description)
	}
}

func UpdateTask(id int, name *string, description *string, state *string) {
	task, err := store.ReadTask(StorageType, id)
	if err != nil {
		slog.Error("failed to read task", "storage", StorageType, "error", err)
		return
	}

	changed := false
	if *name != "" {
		task.Name = *name
		changed = true
	}
	if *description != "" {
		task.Description = *description
		changed = true
	}
	if *state != "" {
		task.State = *state
		changed = true
	}

	if !changed {
		slog.Error("no update fields provided; use --name/--description/--state")
		return
	}

	if err := store.UpdateTask(StorageType, *task); err != nil {
		slog.Error("failed to update task", "storage", StorageType, "error", err)
		return
	}
	slog.Info("Updated task", "id", id)
}
