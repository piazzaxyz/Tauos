package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitHubClient struct {
	Token string
	Owner string
	Repo  string
}

type githubIssueReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *GitHubClient) CreateIssue(title, body string) (string, error) {
	payload := githubIssueReq{Title: title, Body: body}
	data, _ := json.Marshal(payload)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", c.Owner, c.Repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("github API: %s", resp.Status)
	}

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	if htmlURL, ok := result["html_url"].(string); ok {
		return htmlURL, nil
	}
	return "", nil
}
