package ETEHelper

import "os"

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
