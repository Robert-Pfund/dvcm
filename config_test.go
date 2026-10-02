package main_test

import (
	"reflect"
	"testing"

	main "github.com/Robert-Pfund/dvcm"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		wantErr    bool
		wantCfg    main.Config
	}{
		{
			name:       "happy path",
			configFile: "testdata/config_testdata.json",
			wantErr:    false,
			wantCfg: main.Config{
				Default: "github",
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "my_testuser",
					RepoName:  "my_testrepo",
					Token:     "my_secrettoken1",
				},
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "11111141",
					Branch:    "main",
					Token:     "my_secrettoken2",
				},
			},
		},
		{
			name:       "malformed json",
			configFile: "testdata/config_malformed_testdata.json",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main.Cfg = main.Config{}
			gotErr := main.Load(tt.configFile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Load() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Load() succeeded unexpectedly")
			} else {

				if !reflect.DeepEqual(main.Cfg, tt.wantCfg) {
					t.Errorf("Load() = %v, want %v", main.Cfg, tt.wantCfg)
				}
			}
		})
	}
}
