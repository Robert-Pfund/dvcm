package main

import "fmt"

const expectedProjectIdLength int = 8

func validateGithubConfig(cfg Config) error {

	if cfg.Github.RepoOwner == "" || cfg.Github.RepoName == "" || cfg.Github.Token == "" {

		return fmt.Errorf("github config is incomplete - please check your config file")
	}

	return nil
}

func validateGitlabConfig(cfg Config) error {

	if cfg.Gitlab.ProjectId == "" || cfg.Gitlab.Branch == "" || cfg.Gitlab.Token == "" {

		return fmt.Errorf("gitlab config is incomplete - please check your config file")
	}

	if len(cfg.Gitlab.ProjectId) != expectedProjectIdLength {

		return fmt.Errorf("gitlab project id is invalid (expected %d characters but got %d) - please check your config file", expectedProjectIdLength, len(cfg.Gitlab.ProjectId))
	}

	return nil
}
