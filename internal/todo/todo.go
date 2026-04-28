package todo

import (
	"cli-todo-app/internal/store"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const StorageType = "sqlite"

const (
	defaultTerminalWidth = 100
	minTerminalWidth     = 60
	minNameWidth         = 12
	minStateWidth        = 8
	minDescWidth         = 20
	maxNameWidth         = 24
	maxStateWidth        = 12
	maxDescriptionLines  = 2
)

func CreateTask(name string, description string, state string) {
	tasks, err := store.ReadTasks(StorageType, store.TaskFilter{})
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

func GetTasks(filter store.TaskFilter) {
	tasks, err := store.ReadTasks(StorageType, filter)
	if err != nil {
		slog.Error("failed to read tasks", "storage", StorageType, "error", err)
		return
	}

	if len(tasks) == 0 {
		slog.Info("No tasks found.")
		return
	}

	printTasks(tasks)
}

func printTasks(tasks []store.Task) {
	idWidth := 4
	for _, task := range tasks {
		current := len(strconv.Itoa(task.ID))
		if current > idWidth {
			idWidth = current
		}
	}

	terminalWidth := getTerminalWidth()
	nameWidth := maxNameWidth
	stateWidth := maxStateWidth
	separatorWidth := len(" | ") * 3
	descWidth := terminalWidth - idWidth - nameWidth - stateWidth - separatorWidth

	if descWidth < minDescWidth {
		missing := minDescWidth - descWidth

		shrinkName := nameWidth - minNameWidth
		if shrinkName > 0 {
			take := minInt(shrinkName, missing)
			nameWidth -= take
			missing -= take
		}

		shrinkState := stateWidth - minStateWidth
		if shrinkState > 0 && missing > 0 {
			take := minInt(shrinkState, missing)
			stateWidth -= take
			missing -= take
		}

		descWidth = minDescWidth
	}

	header := fmt.Sprintf(
		"%-*s | %-*s | %-*s | %s",
		idWidth,
		"ID",
		nameWidth,
		"Name",
		stateWidth,
		"State",
		"Description",
	)
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(header)))

	for _, task := range tasks {
		name := truncateText(task.Name, nameWidth)
		state := truncateText(task.State, stateWidth)
		descriptionLines := wrapAndClampText(task.Description, descWidth, maxDescriptionLines)

		for index, line := range descriptionLines {
			if index == 0 {
				fmt.Printf(
					"%-*d | %-*s | %-*s | %s\n",
					idWidth,
					task.ID,
					nameWidth,
					name,
					stateWidth,
					state,
					line,
				)
				continue
			}

			fmt.Printf(
				"%-*s | %-*s | %-*s | %s\n",
				idWidth,
				"",
				nameWidth,
				"",
				stateWidth,
				"",
				line,
			)
		}
	}
}

func getTerminalWidth() int {
	fd := int(os.Stdout.Fd())

	if !term.IsTerminal(fd) {
		return defaultTerminalWidth
	}

	width, _, err := term.GetSize(fd)
	if err != nil || width < minTerminalWidth {
		return defaultTerminalWidth
	}

	return width
}

func truncateText(value string, width int) string {
	if width <= 0 {
		return ""
	}

	runes := []rune(value)
	if len(runes) <= width {
		return value
	}

	if width <= 3 {
		return string(runes[:width])
	}

	return string(runes[:width-3]) + "..."
}

func wrapAndClampText(value string, width int, maxLines int) []string {
	if width <= 0 {
		return []string{""}
	}

	if strings.TrimSpace(value) == "" {
		return []string{""}
	}

	words := strings.Fields(value)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, maxLines)
	current := ""

	for _, word := range words {
		remaining := word
		for utf8.RuneCountInString(remaining) > width {
			prefix := string([]rune(remaining)[:width])
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
			lines = append(lines, prefix)
			remaining = string([]rune(remaining)[width:])
		}

		if remaining == "" {
			continue
		}

		candidate := remaining
		if current != "" {
			candidate = current + " " + remaining
		}

		if utf8.RuneCountInString(candidate) <= width {
			current = candidate
		} else {
			if current != "" {
				lines = append(lines, current)
			}
			current = remaining
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	if len(lines) <= maxLines {
		return lines
	}

	trimmed := lines[:maxLines]
	trimmed[maxLines-1] = truncateText(trimmed[maxLines-1], width)
	if !strings.HasSuffix(trimmed[maxLines-1], "...") {
		trimmed[maxLines-1] = truncateText(trimmed[maxLines-1]+"...", width)
	}

	return trimmed
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
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
