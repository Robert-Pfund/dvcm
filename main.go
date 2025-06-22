package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var ErrHelp = errors.New("flag: help requested")

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

type File struct {
	Name           string
	Content        string
	EncodedContent string
}

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
	fmt.Printf("loading %s from %s to %s", app.Name, app.Origin, app.Workspace)
}
func (app *App) saveToLocal() {
	fmt.Printf("saving from %s to %s as %s", app.Workspace, app.Origin, app.Name)
}
