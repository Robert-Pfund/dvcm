package main

import (
	"fmt"
	"os"
	"path"
	"strings"
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

func transferFilesBetweenDirectories(target, source string) error {

	fileNames, err := getFilesByNameInDirectory(source)
	if err != nil {

		return fmt.Errorf("failed to get files for source directory: %s\n", source)
	}

	// TODO: use fileInfo to check if specified target is file/directory
	_, err = os.Stat(target)
	if err != nil {

		if !os.IsNotExist(err) {

			return fmt.Errorf("failed to get fileInfo for target %s: %s\n", target, err)

		} else {

			err = os.Mkdir(target, 0744)
			if err != nil {

				return fmt.Errorf("failed to create directory for target %s: %s\n", target, err)
			}
		}
	}

	for i := range fileNames {

		sourceFile := path.Join(source, fileNames[i])
		targetFile := path.Join(target, fileNames[i])

		err := os.Link(sourceFile, targetFile)
		if err != nil {

			return fmt.Errorf("failed to create hard link for %s: %s\n", sourceFile, err)
		}
	}

	return nil
}

// expecting url like "https://github.com/Robert-Pfund/devcs"
func splitGithubOriginIntoComponents(origin string) (string, string) {

	components := strings.Split(origin, "/")
	repoOwner := components[len(components)-2]
	repoName := components[len(components)-1]

	return repoOwner, repoName
}

// expecting url like "https://gitlab.com/api/v4/projects/:id"
func splitGitlabOriginIntoComponents(origin string) string {

	components := strings.Split(origin, "/")
	projectId := components[len(components)-1]
	return projectId
}
