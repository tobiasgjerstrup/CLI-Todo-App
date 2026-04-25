package store

import "fmt"

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
}

type JsonObject struct {
	Tasks []Task `json:"tasks"`
}

func WriteTask(storageType string, task Task) error {
	if storageType == "json" {
		jsonObj, err := readJSONFromFile("store.json")
		if err != nil {
			return err
		}
		jsonObj.Tasks = append(jsonObj.Tasks, task)
		err = writeJSONToFile("store.json", jsonObj)
		if err != nil {
			return err
		}
		return nil
	}

	if storageType == "sqlite" {
		return writeTaskToSQLite(task)
	}

	return fmt.Errorf("unsupported storage type: %s", storageType)
}

func ReadTasks(storageType string) ([]Task, error) {
	if storageType == "json" {
		jsonObj, err := readJSONFromFile("store.json")
		if err != nil {
			return nil, err
		}
		return jsonObj.Tasks, nil
	}

	if storageType == "sqlite" {
		return readTasksFromSQLite()
	}

	return nil, fmt.Errorf("unsupported storage type: %s", storageType)
}
