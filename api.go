package main

import (
	"net/http"
)

type RemoteRepository interface {
	addHeaders(http.Request) http.Request
	getDownloadResponse() RemoteDownloadResponse
	setDownloadResponse(RemoteDownloadResponse)
	getUploadBody() RemoteUploadBody
	setUploadBody(RemoteUploadBody)
	getRepositoryInfoUrl(Config) string // URL to request general information about contents in repository
	getRepositoryFileUrl(Config) string // URL to request file contents in repository
	getFileUploadHttpMethod() string
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
