package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
}

type JsonObject struct {
	Tasks []Task `json:"tasks"`
}

func readJSONFromFile(filePath string) (JsonObject, error) {
	var jsonObj JsonObject
	defaultJSON := []byte("{\n  \"tasks\": []\n}\n")

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(filePath, defaultJSON, 0o644)
			if err != nil {
				return jsonObj, err
			}

			return JsonObject{Tasks: []Task{}}, nil
		}

		return jsonObj, err
	}

	err = json.Unmarshal(bytes, &jsonObj)
	if err != nil {
		return jsonObj, err
	}

	return jsonObj, nil
}

func writeJSONToFile(filePath string, jsonObj JsonObject) error {
	bytes, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, bytes, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	println("Hello, World!")

	if len(os.Args) < 2 {
		println("Usage: <program> <args>")
		println("Use help for more info")
		return
	}

	if os.Args[1] == "help" {
		println("Usage: <program> <args>")
		println("Available commands: help, create, get, delete, update")
		return
	}

	if os.Args[1] == "create" {
		// println("Creating a new task...")

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

		jsonObj, err := readJSONFromFile("store.json")
		if err != nil {
			fmt.Println("Failed to read JSON file:", err)
			return
		}

		jsonObj.Tasks = append(jsonObj.Tasks, Task{
			ID:          len(jsonObj.Tasks) + 1,
			Name:        name,
			Description: description,
			State:       state,
		})

		err = writeJSONToFile("store.json", jsonObj)
		if err != nil {
			fmt.Println("Failed to write JSON file:", err)
			return
		}

		fmt.Printf("Created task: %s:%d\n", jsonObj.Tasks[len(jsonObj.Tasks)-1].Name, jsonObj.Tasks[len(jsonObj.Tasks)-1].ID)
		// fmt.Printf("Loaded %d tasks from store.json\n", len(jsonObj.Tasks))
		return
	}

	if os.Args[1] == "get" || os.Args[1] == "delete" || os.Args[1] == "update" {
		jsonObj, err := readJSONFromFile("store.json")
		if err != nil {
			fmt.Println("Failed to read JSON file:", err)
			return
		}

		fmt.Printf("%+v\n", jsonObj)
		return
	}

	firstArg := os.Args[1]

	println("First argument:", firstArg)
	// TODO: setup cli commands

	// TODO: setup data storage

	// TODO: setup logging

	// TODO: structure the project properly
}
