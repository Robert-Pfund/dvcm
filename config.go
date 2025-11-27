package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name   string
	Github struct {
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
