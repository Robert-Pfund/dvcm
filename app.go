package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type App struct {
	Workspace string
	Origin    string
	Name      string
	DvcFolder string
}

func (app *App) loadFromLocal() {

	source := path.Join(app.Origin, app.Name)
	target := path.Join(app.Workspace, app.DvcFolder)

	transferFilesBetweenDirectories(target, source)
}

func (app *App) saveToLocal() {

	source := path.Join(app.Workspace, app.DvcFolder)
	target := path.Join(app.Origin, app.Name)

	transferFilesBetweenDirectories(target, source)
}

func (app *App) loadFromRemote() {

	err := Load()
	if err != nil {

		fmt.Printf("failed to load configuration from file: %s\n", err)
		os.Exit(1)
	}
	Cfg.Name = app.Name

	target := path.Join(app.Workspace, app.DvcFolder)

	ghr := new(GithubDownloadResponse)
	var files []GithubDownloadedFile
	ghr.Files = files

	ghRepo := &GithubRepository{
		DownloadResponse: ghr,
	}

	url := ghRepo.getRepositoryInfoUrl(Cfg)

	request, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {

		fmt.Printf("failed to build request: %s\n", err)
		os.Exit(1)
	}

	ghRepo.addHeaders(*request)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {

		fmt.Printf("failed to send request to github api: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	err = ghRepo.DownloadResponse.setData(*response)
	if err != nil {

		fmt.Printf("failed to create a request: %s\n", err)
		os.Exit(1)
	}

	// TODO: use fileInfo to check if specified target is file/directory
	_, err = os.Stat(target)
	if err != nil {

		if !os.IsNotExist(err) {

			fmt.Printf("failed to get fileInfo for target %s: %s\n", target, err)
			os.Exit(1)
		} else {

			err = os.Mkdir(target, os.FileMode(0744))
			if err != nil {

				fmt.Printf("failed to create directory for target %s: %s\n", target, err)
				os.Exit(1)
			}
		}
	}

	fileIndex := 0
	for fileIndex < ghRepo.DownloadResponse.getFileNumber() {

		file := ghRepo.DownloadResponse.getFileAtIndex(fileIndex)

		downloads, err := http.Get(file.getUrl())
		if err != nil {
			log.Fatalln("Failed to send download request:", err)
		}
		defer downloads.Body.Close()

		var data []byte

		data, err = io.ReadAll(downloads.Body)
		if err != nil {
			log.Fatal(err)
		}

		file.setData(data)

		err = os.WriteFile(path.Join(app.Workspace, app.DvcFolder, file.getFilename()), file.getData(), 0666)
		if err != nil {
			log.Printf("Error writing file: %s\n", err)
			return
		}
		fileIndex++
	}
}

func (app *App) saveToRemote() {

	err := Load()
	if err != nil {
		fmt.Printf("failed to load configuration from file: %s\n", err)
		os.Exit(1)
	}
	Cfg.Name = app.Name

	source := path.Join(app.Workspace, app.DvcFolder)

	fileNames, err := getFilesByNameInDirectory(source)
	if err != nil {
		fmt.Printf("failed to get files for source directory: %s\n", source)
		os.Exit(1)
	}

	ghRepo := &GithubRepository{
		UploadBody: &GithubUploadBody{},
	}

	client := http.Client{}

	for _, filename := range fileNames {

		url := fmt.Sprintf(ghRepo.getRepositoryFileUrl(Cfg), filename)
		ghRepo.UploadBody.setMessage(fmt.Sprintf("uploading contents of file: %s in config for %s\n", filename, app.Name))

		contentBytes, err := os.ReadFile(path.Join(source, filename))
		if err != nil {
			fmt.Printf("failed to read file %s: %s\n", filename, err)
			os.Exit(1)
		}

		ghRepo.UploadBody.setContent(base64.StdEncoding.EncodeToString(contentBytes))

		bodyJSON, err := ghRepo.UploadBody.getJson()
		if err != nil {
			fmt.Printf("failed to marshal upload data to json: %s\n", err)
			os.Exit(1)
		}

		request, err := http.NewRequest(
			"PUT",
			url,
			bytes.NewReader(bodyJSON),
		)
		if err != nil {
			fmt.Printf("failed to build request: %s\n", err)
			os.Exit(1)
		}

		ghRepo.addHeaders(*request)

		resp, err := client.Do(request)
		if err != nil {
			fmt.Printf("failed to send create request: %s\n", err)
			os.Exit(1)
		}

		if resp.StatusCode != 201 {

			fmt.Printf("failed to create file in remote origin: %s\n", resp.Status)
			os.Exit(1)
		}
		defer resp.Body.Close()
	}
}
