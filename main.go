package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

var ErrHelp = errors.New("flag: help requested")

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

const defaultFolderName string = ".devcontainer"
const expectedArguments int = 4

var knownCmds = []string{"load", "save"}
var knownRemoteOrigins = []string{"github", "gitlab"}

type runtimeConfig struct {
	workspace    string
	origin       string
	isRemoteMode bool
	cmd          string
	name         string
}

func (r *runtimeConfig) verifyCmd() error {

	if !slices.Contains(knownCmds, r.cmd) {

		return fmt.Errorf("received unexpected command (%s) - not in known commands: %v", r.cmd, knownCmds)

	}

	return nil
}

// for now only check if given origin input is link to http resource and if host is either github or gitlab
func (r *runtimeConfig) verifyOrigin() error {

	if strings.Contains(r.origin, "http") {

		for _, knownOrigin := range knownRemoteOrigins {

			if strings.Contains(r.origin, knownOrigin) {

				return nil
			}
		}

		return fmt.Errorf("received unknown origin (%s) - currently supported remote origins: %v", r.cmd, knownRemoteOrigins)
	}

	return nil
}

func main() {

	var (
		workspace    string
		origin       string
		isRemoteMode bool

		cmd  string
		name string
	)

	// parse flags
	flag.StringVar(&workspace, "workspace", ".", "directory where to load to/save from")
	flag.StringVar(&origin, "origin", "", "directory where to load from/save to")
	flag.BoolVar(&isRemoteMode, "r", false, "toggle remote mode")
	flag.Parse()

	// handle arguments
	amountOfParams := len(flag.Args())
	if amountOfParams > 1 && amountOfParams < 3 {
		cmd = flag.Arg(0)
		name = flag.Arg(1)
	} else {
		fmt.Println("expected 2 (load/save, name) arguments to set but found:", amountOfParams)
		os.Exit(1)
	}

	// create runtime config
	runtime, err := parseArgs([]string{
		workspace,
		origin,
		cmd,
		name,
	},
		isRemoteMode,
	)

	// load config to use as fall-back values
	err = Load()
	if err != nil {
		fmt.Printf("failed to load configuration from file: %s\n", err)
		os.Exit(1)
	}
	Cfg.Name = runtime.name

	// unify all configuration and build app
	app, err := buildApp(runtime, Cfg)
	if err != nil {
		fmt.Printf("failed to build app: %s\n", err)
		os.Exit(1)
	}

	// run app
	err = run(app, runtime)
	if err != nil {
		fmt.Printf("error occured while running app: %s", err)
	}
}

func buildApp(runtimeCfg runtimeConfig, cfg Config) (App, error) {

	var remote RemoteRepository
	if runtimeCfg.isRemoteMode {

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
			return App{}, fmt.Errorf("found %s to be unknown source for remote origin\n", runtimeCfg.origin)
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

func parseArgs(args []string, isRemoteMode bool) (runtimeConfig, error) {

	var runtime runtimeConfig

	if len(args) != expectedArguments {

		return runtime, fmt.Errorf("received unexpected number of arguments. expected %d but got %d", expectedArguments, len(args))
	}

	runtime.workspace = args[0]
	runtime.origin = args[1]
	err := runtime.verifyOrigin()
	if err != nil {

		return runtime, err
	}

	runtime.cmd = args[2]
	err = runtime.verifyCmd()
	if err != nil {

		return runtime, err
	}

	runtime.name = args[3]

	runtime.isRemoteMode = isRemoteMode

	return runtime, nil
}

func run(app App, runtime runtimeConfig) error {

	switch runtime.cmd {
	case "load":
		fmt.Printf("loading %s from %s to %s\n", app.Name, app.Origin, app.Workspace)
		if runtime.isRemoteMode {
			fmt.Println("loading from remote")
			err := app.loadFromRemote()
			if err != nil {

				return err
			}
		} else {
			err := app.loadFromLocal()
			if err != nil {

				return err
			}
		}
	case "save":
		fmt.Printf("saving from %s to %s as %s\n", app.Workspace, app.Origin, app.Name)
		if runtime.isRemoteMode {
			fmt.Println("saving to remote")
			err := app.saveToRemote()
			if err != nil {

				return err
			}
		} else {
			err := app.saveToLocal()
			if err != nil {

				return err
			}
		}
	default:
		return fmt.Errorf("unknown command set: %s\n", runtime.cmd)
	}

	return nil
}
