package store

import (
	"encoding/json"
	"os"
)

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
