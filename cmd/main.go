package main

import (
	"cli-todo-app/internal/store"
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

		flags := flag.NewFlagSet("update", flag.ContinueOnError)
		name := flags.String("name", "", "New task name")
		description := flags.String("description", "", "New task description")
		state := flags.String("state", "", "New task state")

		if err := flags.Parse(os.Args[3:]); err != nil {
			slog.Error("invalid update flags", "error", err)
			return
		}

		todo.UpdateTask(id, name, description, state)
		return
	}

	if os.Args[1] == "get" {
		slog.Debug("Listing tasks")
		flags := flag.NewFlagSet("get", flag.ContinueOnError)
		search := flags.String("search", "", "Search")
		state := flags.String("state", "active", "Search State")
		id := flags.Int("id", -1, "Search id")

		if err := flags.Parse(os.Args[2:]); err != nil {
			slog.Error("invalid get flags", "error", err)
			return
		}

		filter := store.TaskFilter{
			Search: *search,
			State:  *state,
		}
		if *id >= 0 {
			filter.ID = id
		}

		todo.GetTasks(filter)
		return
	}

	firstArg := os.Args[1]

	println("First argument:", firstArg)
}
