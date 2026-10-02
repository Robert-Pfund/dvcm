package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type App struct {
	Workspace string
	Origin    string
	Name      string
	DvcFolder string
	Config    Config
	Remote    RemoteRepository
}

func (app *App) loadFromLocal() error {

	source := path.Join(app.Origin, app.Name)
	target := path.Join(app.Workspace, app.DvcFolder)

	err := transferFilesBetweenDirectories(target, source)
	if err != nil {

		return err
	}

	return nil
}

func (app *App) saveToLocal() error {

	source := path.Join(app.Workspace, app.DvcFolder)
	target := path.Join(app.Origin, app.Name)

	err := transferFilesBetweenDirectories(target, source)
	if err != nil {

		return err
	}

	return nil
}

func (app *App) loadFromRemote() error {

	target := path.Join(app.Workspace, app.DvcFolder)

	client := http.Client{}
	url := app.Remote.getRepositoryInfoUrl(app.Config)
	request, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {

		return fmt.Errorf("failed to build request: %s\n", err)
	}
	app.Remote.addHeaders(*request)

	response, err := client.Do(request)
	if err != nil {

		return fmt.Errorf("failed to send request to remote api: %s\n", err)
	}
	defer response.Body.Close()

	err = app.Remote.getDownloadResponse().setData(*response)
	if err != nil {

		return fmt.Errorf("failed to create a request: %s\n", err)
	}

	// TODO: use fileInfo to check if specified target is file/directory
	_, err = os.Stat(target)
	if err != nil {

		if !os.IsNotExist(err) {

			return fmt.Errorf("failed to get fileInfo for target %s: %s\n", target, err)
		} else {

			err = os.Mkdir(target, os.FileMode(0744))
			if err != nil {

				return fmt.Errorf("failed to create directory for target %s: %s\n", target, err)
			}
		}
	}

	fileIndex := 0
	for fileIndex < app.Remote.getDownloadResponse().getFileNumber() {

		file := app.Remote.getDownloadResponse().getFileAtIndex(fileIndex)
		downloads, err := http.Get(file.getUrl())
		if err != nil {

			return fmt.Errorf("Failed to send download request: %s\n", err)
		}
		defer downloads.Body.Close()

		data, err := io.ReadAll(downloads.Body)
		if err != nil {

			return fmt.Errorf("failed to read response body for %s: %s\n", file.getFilename(), err)
		}

		err = file.setData(data)
		if err != nil {

			return fmt.Errorf("failed to set file data for %s: %s\n", file.getFilename(), err)
		}

		err = os.WriteFile(path.Join(app.Workspace, app.DvcFolder, file.getFilename()), file.getData(), 0666)
		if err != nil {

			return fmt.Errorf("Error writing file: %s\n", err)
		}
		fileIndex++
	}

	return nil
}

func (app *App) saveToRemote() error {

	source := path.Join(app.Workspace, app.DvcFolder)
	fileNames, err := getFilesByNameInDirectory(source)
	if err != nil {
		return fmt.Errorf("failed to get files for source directory: %s\n", source)
	}

	client := http.Client{}

	for _, filename := range fileNames {

		url := fmt.Sprint(app.Remote.getRepositoryFileUrl(app.Config), filename)
		app.Remote.getUploadBody().setMessage(fmt.Sprintf("uploading contents of file: %s in config for %s\n", filename, app.Name))

		contentBytes, err := os.ReadFile(path.Join(source, filename))
		if err != nil {

			return fmt.Errorf("failed to read file %s: %s\n", filename, err)
		}

		app.Remote.getUploadBody().setContent(contentBytes)
		bodyJSON, err := app.Remote.getUploadBody().getJson()
		if err != nil {
			return fmt.Errorf("failed to marshal upload data to json: %s\n", err)
		}

		request, err := http.NewRequest(
			app.Remote.getFileUploadHttpMethod(),
			url,
			bytes.NewReader(bodyJSON),
		)
		if err != nil {

			return fmt.Errorf("failed to build request: %s\n", err)
		}
		app.Remote.addHeaders(*request)

		resp, err := client.Do(request)
		if err != nil {

			return fmt.Errorf("failed to send create request: %s\n", err)
		}

		if resp.StatusCode != 201 {

			return fmt.Errorf("failed to create file in remote origin: %s\n", resp.Status)
		}
		defer resp.Body.Close()
	}

	return nil
}
