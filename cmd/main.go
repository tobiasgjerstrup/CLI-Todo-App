package main

import (
	InputParser "cli-todo-app/internal/parser"
	"cli-todo-app/internal/store"
	"cli-todo-app/internal/todo"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

func printHelp(program string) {
	program = filepath.Base(program)

	fmt.Printf("Usage:\n")
	fmt.Printf("  %s <command> [arguments]\n\n", program)
	fmt.Println("Commands:")
	fmt.Println("  help                               Show this help message")
	fmt.Println("  create [name] [description] [state]")
	fmt.Println("                                     Create a new task")
	fmt.Println("  get [--search text] [--state s] [--id n]")
	fmt.Println("                                     List tasks using optional filters")
	fmt.Println("  update <id> [--name n] [--description d] [--state s]")
	fmt.Println("                                     Update fields for a task")
	fmt.Println()
	fmt.Println("Defaults:")
	fmt.Println("  create state: pending")
	fmt.Println("  get state: '' (anything but removed and completed)")
	fmt.Println("  get id: -1 (disabled)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s create \"Buy milk\" \"2%% from store\" pending\n", program)
	fmt.Printf("  %s get --search milk\n", program)
	fmt.Printf("  %s get --state active\n", program)
	fmt.Printf("  %s get --id 3\n", program)
	fmt.Printf("  %s update 3 --state done\n", program)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if len(os.Args) < 2 {
		printHelp(os.Args[0])
		return
	}

	if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp(os.Args[0])
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
		state := flags.String("state", "", "Search State")
		id := flags.Int("id", -1, "Search id")
		where := flags.String("where", "", "where clause")

		if err := flags.Parse(os.Args[2:]); err != nil {
			slog.Error("invalid get flags", "error", err)
			return
		}

		filter := store.TaskFilter{
			Search: *search,
		}
		if *state != "" {
			filter.State = *state
		}
		if *id >= 0 {
			filter.ID = id
		}
		if *where != "" {
			where, err := InputParser.Parse(where)
			if err != nil {
				slog.Error("something went wrong parsing where clause", "error", err)
			}
			filter.Where = where
		}

		todo.GetTasks(filter)
		return
	}

	slog.Error("unknown command", "command", os.Args[1])
	printHelp(os.Args[0])
}
