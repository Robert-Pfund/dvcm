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
		fmt.Printf("loading from: %s to %s", app.Origin, app.Workspace)
	case "save":
		fmt.Printf("saving from %s to: %s", app.Workspace, app.Origin)
	default:
		fmt.Println("unknown command set: ", params["cmd"])
		os.Exit(1)
	}
}
