package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGitlabRepository_UrlsAndMethod(t *testing.T) {
	Cfg = Config{}
	Cfg.Gitlab.ProjectId = "12345"
	Cfg.Name = "images"

	repo := &GitlabRepository{}

	if got, want := repo.getRepositoryInfoUrl(Cfg), "https://gitlab.com/api/v4/projects/12345/repository/tree?path=images"; got != want {
		t.Fatalf("getRepositoryInfoUrl() = %q, want %q", got, want)
	}

	if got, want := repo.getRepositoryFileUrl(Cfg), "https://gitlab.com/api/v4/projects/12345/repository/files/images%2F"; got != want {
		t.Fatalf("getRepositoryFileUrl() = %q, want %q", got, want)
	}

	if got, want := repo.getFileUploadHttpMethod(), "POST"; got != want {
		t.Fatalf("getFileUploadHttpMethod() = %q, want %q", got, want)
	}
}

func TestGitlabUploadBody_getJson(t *testing.T) {
	Cfg = Config{}
	Cfg.Gitlab.Branch = "main"
	Cfg.Gitlab.ProjectId = "987"

	body := &GitlabUploadBody{}
	body.setMessage("uploading file")
	body.setContent([]byte("hello"))

	jsonBytes, err := body.getJson()
	if err != nil {
		t.Fatalf("getJson() returned error: %v", err)
	}

	got := string(jsonBytes)
	if !strings.Contains(got, `"branch":"main"`) {
		t.Fatalf("JSON missing branch: %s", got)
	}
	if !strings.Contains(got, `"id":"987"`) {
		t.Fatalf("JSON missing project id: %s", got)
	}
	if !strings.Contains(got, `"commit_message":"uploading file"`) {
		t.Fatalf("JSON missing commit message: %s", got)
	}
	if !strings.Contains(got, `"content":"hello"`) {
		t.Fatalf("JSON missing content: %s", got)
	}
}

func TestGitlabDownloadedFile_setData(t *testing.T) {
	file := &GitlabDownloadedFile{}

	payload := `{"file_name":"hello.txt","content":"aGVsbG8="}`
	err := file.setData([]byte(payload))
	if err != nil {
		t.Fatalf("setData() returned error: %v", err)
	}

	if got, want := string(file.Data), "hello"; got != want {
		t.Fatalf("file.Data = %q, want %q", got, want)
	}
}

func TestGitlabDownloadResponse_setData(t *testing.T) {
	body := io.NopCloser(strings.NewReader(`
    [
        {"name":"a.txt","path":"a.txt","id":"1"},
        {"name":"b.txt","path":"b.txt","id":"2"}
    ]
    `))

	resp := http.Response{Body: body}

	var response GitlabDownloadResponse
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
}
