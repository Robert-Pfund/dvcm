package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Github struct {
		RepoOwner string `json:"repoowner"`
		RepoName  string `json:"reponame"`
		Token     string `json:"token"`
	} `json:"github"`
}

var Cfg Config

func Load() (err error) {
	rawData, err := os.ReadFile("config.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(rawData, &Cfg)
	if err != nil {
		return
	}
	return
}
