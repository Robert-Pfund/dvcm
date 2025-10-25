package main

type GithubResponse struct {
	Files []struct {
		FileName    string `json:"name"`
		Path        string `json:"path"`
		Type        string `json:"type"`
		Sha         string `json:"sha"`
		DownloadURL string `json:"download_url"`
		Data        []byte
	} `json:"entries"`
}
