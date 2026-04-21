package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type BlockType int

const (
	BlockText    BlockType = iota
	BlockHeading BlockType = iota
	BlockTodo    BlockType = iota
)

type Block struct {
	ID      string    `json:"id"`
	Type    BlockType `json:"type"`
	Content string    `json:"content"`
	Done    bool      `json:"done,omitempty"`
}

type Page struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Blocks []Block `json:"blocks"`
}

func newID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func storageDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "tauos")
}

func storagePath() string {
	return filepath.Join(storageDir(), "pages.json")
}

func loadPages() []Page {
	data, err := os.ReadFile(storagePath())
	if err != nil {
		return defaultPages()
	}
	var pages []Page
	if err := json.Unmarshal(data, &pages); err != nil {
		return defaultPages()
	}
	return pages
}

func savePages(pages []Page) {
	if err := os.MkdirAll(storageDir(), 0755); err != nil {
		return
	}
	data, _ := json.MarshalIndent(pages, "", "  ")
	_ = os.WriteFile(storagePath(), data, 0644)
}

func defaultPages() []Page {
	return []Page{
		{
			ID:    newID(),
			Title: "Bem-vindo",
			Blocks: []Block{
				{ID: newID(), Type: BlockHeading, Content: "Bem-vindo ao Tau"},
				{ID: newID(), Type: BlockText, Content: "Um espaço simples para suas anotações e tarefas."},
				{ID: newID(), Type: BlockTodo, Content: "Explorar o app"},
				{ID: newID(), Type: BlockTodo, Content: "Criar minha primeira nota", Done: true},
			},
		},
		{
			ID:    newID(),
			Title: "Ideias",
			Blocks: []Block{
				{ID: newID(), Type: BlockHeading, Content: "Minhas Ideias"},
				{ID: newID(), Type: BlockText, Content: "Escreva aqui suas ideias e pensamentos..."},
			},
		},
	}
}
