package main

import (
	"cli-todo-app/internal/store"
	"fmt"
	"log/slog"
	"os"
)

const storageType = "sqlite"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if len(os.Args) < 2 {
		slog.Info("Usage: <program> <args>")
		slog.Info("Use help for more info")
		return
	}

	if os.Args[1] == "help" {
		slog.Info("Usage: <program> <args>")
		slog.Info("Available commands: help, create, get, delete, update")
		return
	}

	if os.Args[1] == "create" {
		slog.Debug("Creating a new task...")

		name := "Unnamed Task"
		description := ""
		state := "pending"

		if len(os.Args) > 2 {
			name = os.Args[2]
		}

		if len(os.Args) > 3 {
			description = os.Args[3]
		}

		if len(os.Args) > 4 {
			state = os.Args[4]
		}

		tasks, err := store.ReadTasks(storageType)
		if err != nil {
			slog.Error("failed to read tasks", "storage", storageType, "error", err)
			return
		}

		nextID := len(tasks) + 1

		err = store.WriteTask(storageType, store.Task{
			ID:          nextID,
			Name:        name,
			Description: description,
			State:       state,
		})
		if err != nil {
			slog.Error("failed to write task", "storage", storageType, "error", err)
			return
		}

		slog.Info("Created task", "name", name, "id", nextID)
		// fmt.Printf("Loaded %d tasks from store.json\n", len(tasks.Tasks))
		return
	}

	if os.Args[1] == "get" || os.Args[1] == "delete" || os.Args[1] == "update" {
		tasks, err := store.ReadTasks(storageType)
		if err != nil {
			slog.Error("failed to read tasks", "storage", storageType, "error", err)
			return
		}

		if len(tasks) == 0 {
			slog.Info("No tasks found.")
			return
		}

		slog.Debug("Listing tasks")
		fmt.Printf("%-4s  %-20s  %-10s  %s\n", "ID", "Name", "State", "Description")
		fmt.Println("----  --------------------  ----------  -----------")
		for _, task := range tasks {
			fmt.Printf("%-4d  %-20s  %-10s  %s\n", task.ID, task.Name, task.State, task.Description)
		}
		return
	}

	firstArg := os.Args[1]

	println("First argument:", firstArg)
	// TODO: setup data storage (SQLite)

	// TODO: structure the project properly
}
