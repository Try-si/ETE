package ETEHelper

import (
	"encoding/json"
	"os"
)

func GetAllFilesInDir(path string) []string {
	// TODO: implement
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
}

func SaveFile[T any](path string, data T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0o644)
	if err != nil {
		return err
	}

	return nil
}
