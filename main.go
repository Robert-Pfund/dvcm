package main

import (
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

type App struct {
	Workspace string
	Origin    string
	Name      string
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
		app.loadFromLocal()
	case "save":
		fmt.Printf("saving from %s to %s as %s\n", app.Workspace, app.Origin, app.Name)
	default:
		fmt.Println("unknown command set: \n", params["cmd"])
		os.Exit(1)
	}
}

func (app *App) loadFromLocal() {

	source := path.Join(app.Origin, app.Name)
	target := path.Join(app.Workspace, ".devcontainer")

	fileNames, err := getFilesByNameInDirectory(source)
	if err != nil {
		fmt.Printf("failed to get files for source directory: %s", source)
		os.Exit(1)
	}

	// TODO: use fileInfo to check if specified target is file/directory
	_, err = os.Stat(target)
	if err != nil {

		if !os.IsNotExist(err) {

			fmt.Printf("failed to get fileInfo for target %s: %s\n", target, err)
			os.Exit(1)
		} else {

			err = os.Mkdir(target, os.ModeAppend)
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
