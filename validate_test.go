package main

import "testing"

func Test_validateGithubConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid github config",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "my_testuser",
					RepoName:  "my_testrepo",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: false,
		},
		{
			name: "github missing config value 1",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "",
					RepoName:  "my_testrepo",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},
		{
			name: "github missing config value 2",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "my_testuser",
					RepoName:  "",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},
		{
			name: "github missing config value 3",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "my_testuser",
					RepoName:  "my_testrepo",
					Token:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "github missing config value 4",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "",
					RepoName:  "",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},

		{
			name: "github missing config value 5",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "",
					RepoName:  "my_testrepo",
					Token:     "",
				},
			},
			wantErr: true,
		},

		{
			name: "github missing config value 6",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "my_testuser",
					RepoName:  "",
					Token:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "github missing config value 7",
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{
					RepoOwner: "",
					RepoName:  "",
					Token:     "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validateGithubConfig(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("validateGithubConfig() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("validateGithubConfig() succeeded unexpectedly")
			}
		})
	}
}

func Test_validateGitlabConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid gitlab config",
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "12345678",
					Branch:    "main",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: false,
		},
		{
			name: "gitlab missing config value 1",
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "",
					Branch:    "main",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},
		{
			name: "gitlab missing config value 2",
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "12345678",
					Branch:    "",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},
		{
			name: "gitlab missing config value 3",
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "12345678",
					Branch:    "main",
					Token:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "gitlab invalid project id", // project id must be 8 characters long (too short)
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "1234567",
					Branch:    "main",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},
		{
			name: "gitlab invalid project id 2", // project id must be 8 characters long (too long)
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{
					ProjectId: "123456789",
					Branch:    "main",
					Token:     "my_secrettoken1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validateGitlabConfig(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("validateGitlabConfig() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("validateGitlabConfig() succeeded unexpectedly")
			}
		})
	}
}
