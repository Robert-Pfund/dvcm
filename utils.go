package main

import (
	"fmt"
	"os"
	"path"
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

func transferFilesBetweenDirectories(target, source string) {

	fileNames, err := getFilesByNameInDirectory(source)
	if err != nil {
		fmt.Printf("failed to get files for source directory: %s\n", source)
		os.Exit(1)
	}

	// TODO: use fileInfo to check if specified target is file/directory
	_, err = os.Stat(target)
	if err != nil {

		if !os.IsNotExist(err) {

			fmt.Printf("failed to get fileInfo for target %s: %s\n", target, err)
			os.Exit(1)
		} else {

			err = os.Mkdir(target, 0744)
			if err != nil {

				fmt.Printf("failed to create directory for target %s: %s\n", target, err)
				os.Exit(1)
			}
		}
	}

	for i := range fileNames {

		sourceFile := path.Join(source, fileNames[i])
		targetFile := path.Join(target, fileNames[i])

		err := os.Link(sourceFile, targetFile)
		if err != nil {

			fmt.Printf("failed to create hard link: %s\n", err)
			os.Exit(1)
		}
	}
}
