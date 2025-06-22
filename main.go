package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
)

var ErrHelp = errors.New("flag: help requested")

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

type File struct {
	Name           string
	Content        []byte
	EncodedContent string
}

var DevcontainerFolderName string = ".devcontainer"

type App struct {
	Workspace string
	Origin    string
	Name      string
	Data      []File
}

func main() {

	var (
		workspace string
		origin    string
	)

	params := map[string]string{}

	flag.StringVar(&workspace, "workspace", ".", "directory where to load to/save from")
	flag.StringVar(&origin, "origin", "", "directory where to load from/save to")
	flag.Parse()

	params["workspace"] = workspace
	params["origin"] = origin

	amountOfParams := len(flag.Args())
	if amountOfParams > 1 && amountOfParams < 3 {
		params["cmd"] = flag.Arg(0)
		params["name"] = flag.Arg(1)
	} else {
		fmt.Println("expected 2 arguments to set but found:", amountOfParams)
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
		app.loadFromLocal()
	case "save":
		app.saveToLocal()
	default:
		fmt.Println("unknown command set: ", params["cmd"])
		os.Exit(1)
	}
}

func (app *App) loadFromLocal() {
	fmt.Printf("loading %s from %s to %s\n", app.Name, app.Origin, app.Workspace)

	// check if workspace path already contains ".devcontainer" folder
	var pathWithDevcontainerFolderName string
	_, filePath := path.Split(app.Workspace)
	fmt.Println(filePath, DevcontainerFolderName)
	if filePath == DevcontainerFolderName {
		pathWithDevcontainerFolderName = app.Workspace
	} else {
		pathWithDevcontainerFolderName = path.Join(app.Workspace, DevcontainerFolderName)
	}

	fmt.Println(pathWithDevcontainerFolderName)
	// check if workspace contains devcontainer configuration already
	workspaceInfo, err := os.Stat(pathWithDevcontainerFolderName)
	if err != nil {

		if !os.IsNotExist(err) {
			fmt.Printf("failed to check for existing devcontainer configuration in %s: %s", app.Workspace, err)
			os.Exit(1)
		} else {
			fmt.Println("no existing devcontainer configuration found")
		}
	} else {
		if workspaceInfo.IsDir() {
			fmt.Println("found existing devcontainer configuration in workspace: ", app.Workspace)
			os.Exit(1)
		}
	}

	// check if origin exists
	fullOriginPath := path.Join(app.Origin, app.Name)
	originInfo, err := os.Stat(fullOriginPath)
	if err != nil {

		if !os.IsNotExist(err) {
			fmt.Printf("failed to check for existing devcontainer configuration in %s: %s", app.Workspace, err)
			os.Exit(1)
		} else {
			fmt.Println("failed to find the specified configuration in origin:", fullOriginPath)
			os.Exit(1)
		}
	}
	if !originInfo.IsDir() {
		fmt.Println("failed to find the specified configuration in origin: ", fullOriginPath)
		os.Exit(1)
	}

	// try read contents from origin
	dir, err := os.ReadDir(fullOriginPath)
	if err != nil {
		fmt.Printf("failed to read local source dir %s: %s\n", fullOriginPath, err)
		os.Exit(1)
	}
	for _, file := range dir {

		fileName := file.Name()
		fileContents, err := os.ReadFile(path.Join(fullOriginPath, fileName))
		if err != nil {
			fmt.Printf("failed to read file %s: %s", fileName, err)
		}

		app.Data = append(app.Data, File{
			Name:           fileName,
			Content:        fileContents,
			EncodedContent: base64.StdEncoding.EncodeToString(fileContents),
		})
	}

	// create folder for devcontainer configuration
	err = os.Mkdir(DevcontainerFolderName, 0777)
	if err != nil {
		fmt.Println("failed to create folder:", err)
	}

	// write all files to devcontainer configuration folder
	for _, file := range app.Data {

		path := path.Join(pathWithDevcontainerFolderName, file.Name)
		err := os.WriteFile(path, file.Content, 0666)
		if err != nil {
			fmt.Printf("error writing file %s: %s\n", file.Name, err)
			os.Exit(1)
		}
	}
}

func (app *App) saveToLocal() {
	fmt.Printf("saving from %s to %s as %s\n", app.Workspace, app.Origin, app.Name)
}
