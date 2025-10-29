package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var ErrHelp = errors.New("flag: help requested")

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

type App struct {
	Workspace string
	Origin    string
	Name      string
}

func main() {

	var (
		workspace    string
		origin       string
		isRemoteMode bool
	)

	params := map[string]string{}

	flag.StringVar(&workspace, "workspace", ".", "directory where to load to/save from")
	flag.StringVar(&origin, "origin", "", "directory where to load from/save to")
	flag.BoolVar(&isRemoteMode, "r", false, "")
	flag.Parse()

	params["workspace"] = workspace
	params["origin"] = origin

	amountOfParams := len(flag.Args())
	if amountOfParams > 1 && amountOfParams < 3 {
		params["cmd"] = flag.Arg(0)
		params["name"] = flag.Arg(1)
	} else {
		fmt.Println("expected 2 (load/save, name) arguments to set but found:", amountOfParams)
		os.Exit(1)
	}

	requiredParamNotSet := false
	for param, value := range params {

		if !(param == "workspace") && value == "" {
			fmt.Printf("%s has not been set\n", param)
			requiredParamNotSet = true
		}
	}
	if requiredParamNotSet {
		os.Exit(1)
	}

	app := App{
		Workspace: params["workspace"],
		Origin:    params["origin"],
		Name:      params["name"],
	}

	switch params["cmd"] {
	case "load":
		fmt.Printf("loading %s from %s to %s\n", app.Name, app.Origin, app.Workspace)
		if isRemoteMode {
			fmt.Println("loading from remote")
			app.loadFromRemote()
		} else {
			app.loadFromLocal()
		}
	case "save":
		fmt.Printf("saving from %s to %s as %s\n", app.Workspace, app.Origin, app.Name)
		if isRemoteMode {
			fmt.Println("saving to remote")
			app.saveToRemote()
		} else {
			app.saveToLocal()
		}
	default:
		fmt.Println("unknown command set: \n", params["cmd"])
		os.Exit(1)
	}
}

func (app *App) loadFromLocal() {

	source := path.Join(app.Origin, app.Name)
	target := path.Join(app.Workspace, ".devcontainer")

	transferFiles(target, source)
}

func (app *App) saveToLocal() {

	source := path.Join(app.Workspace, ".devcontainer")
	target := path.Join(app.Origin, app.Name)

	transferFiles(target, source)
}

func (app *App) loadFromRemote() {

	err := Load()
	if err != nil {

		fmt.Printf("failed to load configuration from file: %s\n", err)
	}

	target := path.Join(app.Workspace, ".devcontainer")

	urlTemplate := "https://api.github.com/repos/%s/%s/contents/%s"
	url := fmt.Sprintf(urlTemplate, Cfg.Github.RepoOwner, Cfg.Github.RepoName, app.Name)

	request, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {

		fmt.Printf("failed to build request: %s\n", err)
		os.Exit(1)
	}

	request.Header.Add("Accept", "application/vnd.github.object+json")
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Add("Authorization", "Bearer "+Cfg.Github.Token)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {

		fmt.Printf("failed to send request to github api: %s\n", err)
	}
	defer response.Body.Close()

	respContents := &GithubResponse{}

	err = json.NewDecoder(response.Body).Decode(&respContents)
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

	for _, file := range respContents.Files {

		downloads, err := http.Get(file.DownloadURL)
		if err != nil {
			log.Fatalln("Failed to send download request:", err)
		}
		defer downloads.Body.Close()

		file.Data, err = io.ReadAll(downloads.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(path.Join(app.Workspace, ".devcontainer", file.FileName), file.Data, 0666)
		if err != nil {
			log.Printf("Error writing file: %s\n", err)
			return
		}
	}
}

func (app *App) saveToRemote() {

	err := Load()
	if err != nil {
		fmt.Printf("failed to load configuration from file: %s\n", err)
		os.Exit(1)
	}

	source := path.Join(app.Workspace, ".devcontainer")

	fileNames, err := getFilesByNameInDirectory(source)
	if err != nil {
		fmt.Printf("failed to get files for source directory: %s\n", source)
		os.Exit(1)
	}

	urlTemplate := "https://api.github.com/repos/%s/%s/contents/%s/%s"

	body := GithubUploadBody{}
	client := http.Client{}

	for _, filename := range fileNames {

		url := fmt.Sprintf(urlTemplate, Cfg.Github.RepoOwner, Cfg.Github.RepoName, app.Name, filename)
		body.Message = fmt.Sprintf("uploading contents of file: %s in config for %s\n", filename, app.Name)

		contentBytes, err := os.ReadFile(path.Join(source, filename))
		if err != nil {
			fmt.Printf("failed to read file %s: %s\n", filename, err)
			os.Exit(1)
		}

		body.Content = base64.StdEncoding.EncodeToString(contentBytes)

		bodyJSON, err := json.Marshal(body)
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

		request.Header.Add("Accept", "application/vnd.github.object+json")
		request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
		request.Header.Add("Authorization", "Bearer "+Cfg.Github.Token)

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

func transferFiles(target, source string) {

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
