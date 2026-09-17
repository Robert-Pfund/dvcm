package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name    string
	Default string `json:"default_remote"`
	Github  struct {
		RepoOwner string `json:"repoowner"`
		RepoName  string `json:"reponame"`
		Token     string `json:"token"`
	} `json:"github"`
	Gitlab struct {
		ProjectId string `json:"projectid"`
		Branch    string `json:"branch"`
		Token     string `json:"token"`
	} `json:"gitlab"`
}

var Cfg Config

func Load(configFile string) (err error) {
	rawData, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(rawData, &Cfg)
	if err != nil {
		return
	}
	return
}
