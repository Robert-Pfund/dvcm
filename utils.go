package main

import (
	"fmt"
	"os"
)

func getFilesByNameInDirectory(directory string) ([]string, error) {

	var fileNames []string

	dirEntries, err := os.ReadDir(directory)
	if err != nil {

		fmt.Printf("failed to read directory: %s", err)
		return fileNames, err
	}

	for _, entry := range dirEntries {

		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	return fileNames, nil
}
