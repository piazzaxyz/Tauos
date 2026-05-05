package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const notionAPIBase = "https://api.notion.com/v1"

type NotionClient struct {
	Token      string
	DatabaseID string
}

type notionCreatePageReq struct {
	Parent     notionParent               `json:"parent"`
	Properties map[string]notionPropValue `json:"properties"`
	Children   []notionBlock              `json:"children,omitempty"`
}

type notionParent struct {
	DatabaseID string `json:"database_id"`
}

type notionPropValue struct {
	Title []notionRichText `json:"title,omitempty"`
}

type notionRichText struct {
	Text notionText `json:"text"`
}

type notionText struct {
	Content string `json:"content"`
}

type notionBlock struct {
	Object    string           `json:"object"`
	Type      string           `json:"type"`
	Paragraph *notionParagraph `json:"paragraph,omitempty"`
}

type notionParagraph struct {
	RichText []notionRichText `json:"rich_text"`
}

func (c *NotionClient) CreateTask(title, body string) (string, error) {
	payload := notionCreatePageReq{
		Parent: notionParent{DatabaseID: c.DatabaseID},
		Properties: map[string]notionPropValue{
			"Name": {Title: []notionRichText{{Text: notionText{Content: title}}}},
		},
		Children: []notionBlock{
			{
				Object: "block",
				Type:   "paragraph",
				Paragraph: &notionParagraph{
					RichText: []notionRichText{{Text: notionText{Content: body}}},
				},
			},
		},
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", notionAPIBase+"/pages", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("notion API: %s", resp.Status)
	}

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	if id, ok := result["id"].(string); ok {
		return "https://notion.so/" + id, nil
	}
	return "", nil
}
