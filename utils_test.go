package main

import "testing"

func Test_splitGithubOriginIntoComponents(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		wantUser string
		wantRepo string
	}{
		{
			name:     "happy path",
			origin:   "https://github.com/myuser/myrepo",
			wantUser: "myuser",
			wantRepo: "myrepo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotRepo := splitGithubOriginIntoComponents(tt.origin)

			if gotUser != tt.wantUser {
				t.Errorf("splitGithubOriginIntoComponents() = %v, want %v", gotUser, tt.wantUser)
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("splitGithubOriginIntoComponents() = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}

func Test_splitGitlabOriginIntoComponents(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   string
	}{
		{
			name:   "happy path",
			origin: "https://gitlab.com/api/v4/projects/10312419",
			want:   "10312419",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitGitlabOriginIntoComponents(tt.origin)

			if got != tt.want {
				t.Errorf("splitGitlabOriginIntoComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}
