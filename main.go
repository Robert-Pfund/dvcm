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

const defaultFolderName string = ".devcontainer"

type runtimeConfig struct {
	workspace    string
	origin       string
	isRemoteMode bool
	cmd          string
	name         string
}

func main() {

	var runtime runtimeConfig

	// parse flags
	flag.StringVar(&runtime.workspace, "workspace", ".", "directory where to load to/save from")
	flag.StringVar(&runtime.origin, "origin", "", "directory where to load from/save to")
	flag.BoolVar(&runtime.isRemoteMode, "r", false, "toggle remote mode")
	flag.Parse()

	// handle arguments
	amountOfParams := len(flag.Args())
	if amountOfParams > 1 && amountOfParams < 3 {
		runtime.cmd = flag.Arg(0)
		runtime.name = flag.Arg(1)
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
	Cfg.Name = runtime.name

	app, err := buildApp(runtime, Cfg)

	switch runtime.cmd {
	case "load":
		fmt.Printf("loading %s from %s to %s\n", app.Name, app.Origin, app.Workspace)
		if runtime.isRemoteMode {
			fmt.Println("loading from remote")
			app.loadFromRemote()
		} else {
			app.loadFromLocal()
		}
	case "save":
		fmt.Printf("saving from %s to %s as %s\n", app.Workspace, app.Origin, app.Name)
		if runtime.isRemoteMode {
			fmt.Println("saving to remote")
			app.saveToRemote()
		} else {
			app.saveToLocal()
		}
	default:
		fmt.Println("unknown command set: \n", runtime.cmd)
		os.Exit(1)
	}
}

func buildApp(runtimeCfg runtimeConfig, cfg Config) (App, error) {

	var remote RemoteRepository
	if runtimeCfg.isRemoteMode {

		// for now also use github as default case (if no origin is set)
		if strings.Contains(runtimeCfg.origin, "github") {

			repoOwner, repoName := splitGithubOriginIntoComponents(runtimeCfg.origin)
			cfg.Github.RepoOwner = repoOwner
			cfg.Github.RepoName = repoName

			remote = getGithubConfig()

		} else if strings.Contains(runtimeCfg.origin, "gitlab") {

			projectId := splitGitlabOriginIntoComponents(runtimeCfg.origin)
			cfg.Gitlab.ProjectId = projectId

			remote = getGitlabConfig()

		} else if runtimeCfg.origin == "" && cfg.Default == "github" {

			remote = getGithubConfig()

		} else if runtimeCfg.origin == "" && cfg.Default == "gitlab" {

			remote = getGitlabConfig()

		} else {

			fmt.Printf("found %s to be unknown source for remote origin\n", runtimeCfg.origin)
			os.Exit(1)
		}
	}

	app := App{
		Workspace: runtimeCfg.workspace,
		Origin:    runtimeCfg.origin,
		Name:      runtimeCfg.name,
		DvcFolder: defaultFolderName,
		Config:    cfg,
		Remote:    remote,
	}

	return app, nil
}

func getGithubConfig() *GithubRepository {

	var files []GithubDownloadedFile
	remote := &GithubRepository{
		DownloadResponse: &GithubDownloadResponse{
			Files: files,
		},
		UploadBody: &GithubUploadBody{},
	}

	return remote
}

func getGitlabConfig() *GitlabRepository {

	var files []GitlabDownloadedFile
	remote := &GitlabRepository{
		DownloadResponse: &GitlabDownloadResponse{
			Files: files,
		},
		UploadBody: &GitlabUploadBody{},
	}

	return remote
}
