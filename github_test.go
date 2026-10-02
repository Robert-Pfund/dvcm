package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGithubRepository_UrlsAndMethod(t *testing.T) {
	Cfg = Config{}
	Cfg.Github.RepoOwner = "acme"
	Cfg.Github.RepoName = "docs"
	Cfg.Name = "devcontainer.json"

	repo := &GithubRepository{}

	if got, want := repo.getRepositoryInfoUrl(Cfg), "https://api.github.com/repos/acme/docs/contents/devcontainer.json"; got != want {
		t.Fatalf("getRepositoryInfoUrl() = %q, want %q", got, want)
	}

	if got, want := repo.getRepositoryFileUrl(Cfg), "https://api.github.com/repos/acme/docs/contents/devcontainer.json/"; got != want {
		t.Fatalf("getRepositoryFileUrl() = %q, want %q", got, want)
	}

	if got, want := repo.getFileUploadHttpMethod(), "PUT"; got != want {
		t.Fatalf("getFileUploadHttpMethod() = %q, want %q", got, want)
	}
}

func TestGithubUploadBody_getJson(t *testing.T) {
	body := &GithubUploadBody{}
	body.setMessage("commit message")
	body.setContent([]byte("hello"))

	jsonBytes, err := body.getJson()
	if err != nil {
		t.Fatalf("getJson() returned error: %v", err)
	}

	got := string(jsonBytes)
	if !strings.Contains(got, `"message":"commit message"`) {
		t.Fatalf("JSON missing message: %s", got)
	}
	if !strings.Contains(got, `"content":"aGVsbG8="`) {
		t.Fatalf("JSON missing base64 content: %s", got)
	}
}

func TestGithubDownloadResponse_setData(t *testing.T) {
	body := io.NopCloser(strings.NewReader(`
    {
        "entries": [
            {"name":"a.txt","download_url":"https://example.com/a.txt","type":"file"},
            {"name":"b.txt","download_url":"https://example.com/b.txt","type":"file"}
        ]
    }
    `))

	resp := http.Response{Body: body}

	var response GithubDownloadResponse
	if err := response.setData(resp); err != nil {
		t.Fatalf("setData() returned error: %v", err)
	}

	if got, want := response.getFileNumber(), 2; got != want {
		t.Fatalf("getFileNumber() = %d, want %d", got, want)
	}

	first := response.getFileAtIndex(0)
	if got, want := first.getFilename(), "a.txt"; got != want {
		t.Fatalf("first file name = %q, want %q", got, want)
	}
	if got, want := first.getUrl(), "https://example.com/a.txt"; got != want {
		t.Fatalf("first file URL = %q, want %q", got, want)
	}
}
