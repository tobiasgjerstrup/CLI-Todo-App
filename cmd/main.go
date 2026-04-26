package main

import (
	"cli-todo-app/internal/store"
	"cli-todo-app/internal/todo"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

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
		slog.Info("Available commands: help, create, get, update")
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

		todo.CreateTask(name, description, state)
		return
	}

	if os.Args[1] == "update" {
		slog.Debug("Updating a task...")

		if len(os.Args) < 3 {
			slog.Error("Task ID is required for update")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			slog.Error("invalid task ID", "error", err)
			return
		}

		task, err := store.ReadTask(todo.StorageType, id)
		if err != nil {
			slog.Error("failed to read task", "storage", todo.StorageType, "error", err)
			return
		}

		updateCmd := flag.NewFlagSet("update", flag.ContinueOnError)
		name := updateCmd.String("name", "", "New task name")
		description := updateCmd.String("description", "", "New task description")
		state := updateCmd.String("state", "", "New task state")

		if err := updateCmd.Parse(os.Args[3:]); err != nil {
			slog.Error("invalid update flags", "error", err)
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

		if err := store.UpdateTask(todo.StorageType, *task); err != nil {
			slog.Error("failed to update task", "storage", todo.StorageType, "error", err)
			return
		}
		slog.Info("Updated task", "id", id)
		return
	}

	if os.Args[1] == "get" {
		slog.Debug("Listing tasks")
		tasks, err := store.ReadTasks(todo.StorageType)
		if err != nil {
			slog.Error("failed to read tasks", "storage", todo.StorageType, "error", err)
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
		return
	}

	firstArg := os.Args[1]

	println("First argument:", firstArg)
}
