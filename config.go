package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	NotionToken      string `json:"notion_token"`
	NotionDatabaseID string `json:"notion_database_id"`
	GitHubToken      string `json:"github_token"`
	GitHubOwner      string `json:"github_owner"`
	GitHubRepo       string `json:"github_repo"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tauos")
}

func configFilePath() string {
	return filepath.Join(configDir(), "config.json")
}

func loadConfig() Config {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return Config{}
	}
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

func saveConfigTemplate() {
	if err := os.MkdirAll(configDir(), 0755); err != nil {
		return
	}
	if _, err := os.Stat(configFilePath()); err == nil {
		return // already exists
	}
	template := Config{}
	data, _ := json.MarshalIndent(template, "", "  ")
	_ = os.WriteFile(configFilePath(), data, 0600)
}
