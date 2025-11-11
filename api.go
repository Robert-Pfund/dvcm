package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RemoteRepository interface {
	addHeaders(http.Request) http.Request
	getRepositoryInfoUrl(Config) string
	getRepositoryFileUrl(Config) string
}

type DownloadedFile interface {
	getUrl() string
	getData() []byte
	setData([]byte)
	getFilename() string
}

type RemoteDownloadResponse interface {
	getFileAtIndex(int) DownloadedFile
	getFileNumber() int
	setData(http.Response) error
}

type RemoteUploadBody interface {
	setMessage(string)
	setContent(string)
	getJson() ([]byte, error)
}

type GithubRepository struct {
	DownloadResponse RemoteDownloadResponse
	UploadBody       RemoteUploadBody
}

func (repo *GithubRepository) addHeaders(request http.Request) http.Request {

	request.Header.Add("Accept", "application/vnd.github.object+json")
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Add("Authorization", "Bearer "+Cfg.Github.Token)

	return request
}

func (repo *GithubRepository) getRepositoryInfoUrl(cfg Config) string {

	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", cfg.Github.RepoOwner, cfg.Github.RepoName, cfg.Name)
}

func (repo *GithubRepository) getRepositoryFileUrl(cfg Config) string {

	// add "/%s" for specific file name which is set when looping through files in devcontainer folder
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", cfg.Github.RepoOwner, cfg.Github.RepoName, cfg.Name) + "/%s"
}

type GithubDownloadedFile struct {
	FileName    string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Sha         string `json:"sha"`
	DownloadURL string `json:"download_url"`
	Data        []byte
}

func (f GithubDownloadedFile) getUrl() string { return f.DownloadURL }

func (f GithubDownloadedFile) getData() []byte { return f.Data }

func (f *GithubDownloadedFile) setData(data []byte) { f.Data = data }

func (f GithubDownloadedFile) getFilename() string { return f.FileName }

type GithubDownloadResponse struct {
	Files []GithubDownloadedFile `json:"entries"`
}

func (ghr GithubDownloadResponse) getFileAtIndex(idx int) DownloadedFile {

	return &ghr.Files[idx]
}

func (ghr GithubDownloadResponse) getFileNumber() int {

	return len(ghr.Files)
}

func (ghr *GithubDownloadResponse) setData(response http.Response) error {

	err := json.NewDecoder(response.Body).Decode(&ghr) // TODO: make struct be returned instead of interface
	if err != nil {
		return err
	}

	return nil
}

type GithubUploadBody struct {
	Message string `json:"message"`
	Content string `json:"content"`
}

func (b *GithubUploadBody) setMessage(msg string) { b.Message = msg }

func (b *GithubUploadBody) setContent(content string) { b.Content = content }

func (b GithubUploadBody) getJson() ([]byte, error) { return json.Marshal(b) }
