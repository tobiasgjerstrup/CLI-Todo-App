package main

import (
	"cli-todo-app/internal/todo"
	"flag"
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

		updateCmd := flag.NewFlagSet("update", flag.ContinueOnError)
		name := updateCmd.String("name", "", "New task name")
		description := updateCmd.String("description", "", "New task description")
		state := updateCmd.String("state", "", "New task state")

		if err := updateCmd.Parse(os.Args[3:]); err != nil {
			slog.Error("invalid update flags", "error", err)
			return
		}

		todo.UpdateTask(id, name, description, state)
		return
	}

	if os.Args[1] == "get" {
		slog.Debug("Listing tasks")
		todo.GetTasks()
		return
	}

	firstArg := os.Args[1]

	println("First argument:", firstArg)
}
