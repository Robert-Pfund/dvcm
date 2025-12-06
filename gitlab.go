package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GitlabRepository struct {
	DownloadResponse RemoteDownloadResponse
	UploadBody       RemoteUploadBody
}

func (repo *GitlabRepository) addHeaders(request http.Request) http.Request {

	request.Header.Add("PRIVATE-TOKEN", Cfg.Gitlab.Token)
	request.Header.Add("Content-Type", "application/json")

	return request
}

// TODO: check if repository info request is required on gitlab
func (repo *GitlabRepository) getRepositoryInfoUrl(cfg Config) string {

	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/files/", cfg.Gitlab.ProjectId)
}

func (repo *GitlabRepository) getRepositoryFileUrl(cfg Config) string {

	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/files/", cfg.Gitlab.ProjectId) + Cfg.Name + "%2F"
}

func (repo *GitlabRepository) getUploadBody() RemoteUploadBody {

	return repo.UploadBody
}

func (repo *GitlabRepository) setUploadBody(body RemoteUploadBody) {

	repo.UploadBody = body
}

func (repo *GitlabRepository) getDownloadResponse() RemoteDownloadResponse {

	return repo.DownloadResponse
}

func (repo *GitlabRepository) setDownloadResponse(response RemoteDownloadResponse) {

	repo.DownloadResponse = response
}

func (repo *GitlabRepository) getFileUploadHttpMethod() string {

	return "POST"
}

type GitlabDownloadedFile struct {
	FileName string `json:"file_name"`
	Path     string `json:"file_path"`
	Sha      string `json:"blob_id"`
	Data     []byte
}

func (f GitlabDownloadedFile) getUrl() string {

	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/repository/files/", Cfg.Gitlab.ProjectId) + Cfg.Name + "%2F"
}

func (f GitlabDownloadedFile) getData() []byte {

	return f.Data
}

func (f GitlabDownloadedFile) setData(data []byte) {

	f.Data = data
}

func (f GitlabDownloadedFile) getFilename() string {

	return f.FileName
}

type GitlabDownloadResponse struct {
	Files []GitlabDownloadedFile
}

func (glr GitlabDownloadResponse) getFileAtIndex(idx int) DownloadedFile {

	return &glr.Files[idx]
}

func (glr GitlabDownloadResponse) getFileNumber() int {

	return len(glr.Files)
}

func (glr *GitlabDownloadResponse) setData(response http.Response) error {

	err := json.NewDecoder(response.Body).Decode(&glr)
	if err != nil {
		return err
	}

	return nil
}

type GitlabUploadBody struct {
	Branch    string `json:"branch"`
	FilePath  string `json:"file_path"`
	ProjectId string `json:"id"`
	Message   string `json:"commit_message"`
	Content   string `json:"content"`
}

func (b *GitlabUploadBody) setMessage(msg string) {

	b.Message = msg
}

func (b *GitlabUploadBody) setContent(content []byte) {

	b.Content = string(content)
}

func (b GitlabUploadBody) getJson() ([]byte, error) {

	b.Branch = Cfg.Gitlab.Branch
	b.ProjectId = Cfg.Gitlab.ProjectId

	return json.Marshal(b)
}
