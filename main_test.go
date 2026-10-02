package main

import "testing"

func Test_buildApp(t *testing.T) {
	tests := []struct {
		name        string
		runtimeCfg  runtimeConfig
		cfg         Config
		want        App
		wantErr     bool
		wantRemote  string
		wantOwner   string
		wantRepo    string
		wantProject string
	}{
		{
			name: "local load command builds app without remote",
			runtimeCfg: runtimeConfig{
				workspace: "/tmp/workspace",
				origin:    "/tmp/origin",
				cmd:       "load",
				name:      "devcontainer",
			},
			want: App{
				Workspace: "/tmp/workspace",
				Origin:    "/tmp/origin",
				Name:      "devcontainer",
				DvcFolder: defaultFolderName,
			},
		},
		{
			name: "remote github origin builds github remote client",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://github.com/acme/docs",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			cfg: Config{
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{Token: "test-token"},
			},
			want: App{
				Workspace: "/tmp/workspace",
				Origin:    "https://github.com/acme/docs",
				Name:      "devcontainer",
				DvcFolder: defaultFolderName,
			},
			wantRemote: "github",
			wantOwner:  "acme",
			wantRepo:   "docs",
		},
		{
			name: "remote gitlab origin builds gitlab remote client",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://gitlab.com/api/v4/projects/10312419",
				isRemoteMode: true,
				cmd:          "save",
				name:         "devcontainer",
			},
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{Branch: "main", Token: "test-token"},
			},
			want: App{
				Workspace: "/tmp/workspace",
				Origin:    "https://gitlab.com/api/v4/projects/10312419",
				Name:      "devcontainer",
				DvcFolder: defaultFolderName,
			},
			wantRemote:  "gitlab",
			wantProject: "10312419",
		},
		{
			name: "remote mode rejects unknown origin",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://example.com/project",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			wantErr: true,
		},
		{
			name: "default github config is used when origin empty",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			cfg: Config{
				Default: "github",
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{RepoOwner: "acme", RepoName: "docs", Token: "test-token"},
			},
			want: App{
				Workspace: "/tmp/workspace",
				Origin:    "",
				Name:      "devcontainer",
				DvcFolder: defaultFolderName,
			},
			wantRemote: "github",
			wantOwner:  "acme",
			wantRepo:   "docs",
		},
		{
			name: "remote github origin rejects missing token",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://github.com/acme/docs",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			wantErr: true,
		},
		{
			name: "remote gitlab origin rejects missing token",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://gitlab.com/api/v4/projects/10312419",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			cfg: Config{
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{Branch: "main"},
			},
			wantErr: true,
		},
		{
			name: "default github rejects missing token",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			cfg: Config{
				Default: "github",
				Github: struct {
					RepoOwner string `json:"repoowner"`
					RepoName  string `json:"reponame"`
					Token     string `json:"token"`
				}{RepoOwner: "acme", RepoName: "docs"},
			},
			wantErr: true,
		},
		{
			name: "default gitlab rejects missing token",
			runtimeCfg: runtimeConfig{
				workspace:    "/tmp/workspace",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
			cfg: Config{
				Default: "gitlab",
				Gitlab: struct {
					ProjectId string `json:"projectid"`
					Branch    string `json:"branch"`
					Token     string `json:"token"`
				}{ProjectId: "10312419", Branch: "main"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := buildApp(tt.runtimeCfg, tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Fatalf("buildApp() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("buildApp() succeeded unexpectedly")
			}

			if got.Workspace != tt.want.Workspace {
				t.Fatalf("Workspace = %q, want %q", got.Workspace, tt.want.Workspace)
			}
			if got.Origin != tt.want.Origin {
				t.Fatalf("Origin = %q, want %q", got.Origin, tt.want.Origin)
			}
			if got.Name != tt.want.Name {
				t.Fatalf("Name = %q, want %q", got.Name, tt.want.Name)
			}
			if got.DvcFolder != tt.want.DvcFolder {
				t.Fatalf("DvcFolder = %q, want %q", got.DvcFolder, tt.want.DvcFolder)
			}

			switch tt.wantRemote {
			case "github":
				if _, ok := got.Remote.(*GithubRepository); !ok {
					t.Fatalf("Remote type = %T, want *GithubRepository", got.Remote)
				}
				if got.Config.Github.RepoOwner != tt.wantOwner {
					t.Fatalf("Config.Github.RepoOwner = %q, want %q", got.Config.Github.RepoOwner, tt.wantOwner)
				}
				if got.Config.Github.RepoName != tt.wantRepo {
					t.Fatalf("Config.Github.RepoName = %q, want %q", got.Config.Github.RepoName, tt.wantRepo)
				}
			case "gitlab":
				if _, ok := got.Remote.(*GitlabRepository); !ok {
					t.Fatalf("Remote type = %T, want *GitlabRepository", got.Remote)
				}
				if got.Config.Gitlab.ProjectId != tt.wantProject {
					t.Fatalf("Config.Gitlab.ProjectId = %q, want %q", got.Config.Gitlab.ProjectId, tt.wantProject)
				}
			case "":
				if got.Remote != nil {
					t.Fatalf("Remote = %#v, want nil", got.Remote)
				}
			default:
				t.Fatalf("unsupported remote expectation: %q", tt.wantRemote)
			}
		})
	}
}

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		isRemoteMode bool
		want         runtimeConfig
		wantErr      bool
	}{
		{
			name: "local load command is parsed",
			args: []string{
				"/tmp/workspace",
				"/tmp/origin",
				"load",
				"devcontainer",
			},
			want: runtimeConfig{
				workspace: "/tmp/workspace",
				origin:    "/tmp/origin",
				cmd:       "load",
				name:      "devcontainer",
			},
		},
		{
			name: "remote github url is accepted",
			args: []string{
				"/tmp/workspace",
				"https://github.com/acme/docs",
				"load",
				"devcontainer",
			},
			isRemoteMode: true,
			want: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://github.com/acme/docs",
				isRemoteMode: true,
				cmd:          "load",
				name:         "devcontainer",
			},
		},
		{
			name: "remote gitlab url is accepted",
			args: []string{
				"/tmp/workspace",
				"https://gitlab.com/api/v4/projects/10312419",
				"save",
				"devcontainer",
			},
			isRemoteMode: true,
			want: runtimeConfig{
				workspace:    "/tmp/workspace",
				origin:       "https://gitlab.com/api/v4/projects/10312419",
				isRemoteMode: true,
				cmd:          "save",
				name:         "devcontainer",
			},
		},
		{
			name: "invalid remote origin is rejected",
			args: []string{
				"/tmp/workspace",
				"https://example.com/project",
				"load",
				"devcontainer",
			},
			isRemoteMode: true,
			wantErr:      true,
		},
		{
			name: "invalid command is rejected",
			args: []string{
				"/tmp/workspace",
				"/tmp/origin",
				"unknown",
				"devcontainer",
			},
			wantErr: true,
		},
		{
			name: "wrong argument count is rejected",
			args: []string{
				"/tmp/workspace",
				"load",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parseArgs(tt.args, tt.isRemoteMode)
			if gotErr != nil {
				if !tt.wantErr {
					t.Fatalf("parseArgs() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseArgs() succeeded unexpectedly")
			}

			if got.workspace != tt.want.workspace {
				t.Fatalf("workspace = %q, want %q", got.workspace, tt.want.workspace)
			}
			if got.origin != tt.want.origin {
				t.Fatalf("origin = %q, want %q", got.origin, tt.want.origin)
			}
			if got.isRemoteMode != tt.want.isRemoteMode {
				t.Fatalf("isRemoteMode = %v, want %v", got.isRemoteMode, tt.want.isRemoteMode)
			}
			if got.cmd != tt.want.cmd {
				t.Fatalf("cmd = %q, want %q", got.cmd, tt.want.cmd)
			}
			if got.name != tt.want.name {
				t.Fatalf("name = %q, want %q", got.name, tt.want.name)
			}
		})
	}
}
