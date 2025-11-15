package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var ErrHelp = errors.New("flag: help requested")

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	var (
		workspace    string
		origin       string
		isRemoteMode bool
	)

	params := map[string]string{}

	// parse flags
	flag.StringVar(&workspace, "workspace", ".", "directory where to load to/save from")
	flag.StringVar(&origin, "origin", "", "directory where to load from/save to")
	flag.BoolVar(&isRemoteMode, "r", false, "")
	flag.Parse()

	// handle flag values
	params["workspace"] = workspace
	params["origin"] = origin

	// handle arguments
	amountOfParams := len(flag.Args())
	if amountOfParams > 1 && amountOfParams < 3 {
		params["cmd"] = flag.Arg(0)
		params["name"] = flag.Arg(1)
	} else {
		fmt.Println("expected 2 (load/save, name) arguments to set but found:", amountOfParams)
		os.Exit(1)
	}

	// load config to use as fall-back values
	err := Load()
	if err != nil {
		fmt.Printf("failed to load configuration from file: %s\n", err)
		os.Exit(1)
	}
	Cfg.Name = params["name"]

	var remote RemoteRepository
	if isRemoteMode {
		// for now also use github as default case (if no origin is set)
		if strings.Contains(params["origin"], "github") {

			repoOwner, repoName := splitGithubOriginIntoComponents(params["origin"])
			Cfg.Github.RepoOwner = repoOwner
			Cfg.Github.RepoName = repoName

			var files []GithubDownloadedFile
			remote = &GithubRepository{
				DownloadResponse: &GithubDownloadResponse{
					Files: files,
				},
				UploadBody: &GithubUploadBody{},
			}
		} else if params["origin"] == "" {

			var files []GithubDownloadedFile
			remote = &GithubRepository{
				DownloadResponse: &GithubDownloadResponse{
					Files: files,
				},
				UploadBody: &GithubUploadBody{},
			}
		} else {

			fmt.Printf("found %s to be unknown source for remote origin\n", params["origin"])
			os.Exit(1)
		}
	}

	app := App{
		Workspace: params["workspace"],
		Origin:    params["origin"],
		Name:      params["name"],
		DvcFolder: ".devcontainer",
		Config:    Cfg,
		Remote:    remote,
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
