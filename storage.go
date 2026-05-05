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
	BlockHeading           // ## heading
	BlockTodo              // checkbox
	BlockCode              // code block
	BlockQuote             // blockquote
	BlockDivider           // horizontal rule
	BlockCard              // markdown card / topic
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
				{ID: newID(), Type: BlockHeading, Content: "Bem-vindo ao **Tauos**"},
				{ID: newID(), Type: BlockText, Content: "Um espaço minimalista para notas e tarefas. Suporta *markdown* inline."},
				{ID: newID(), Type: BlockQuote, Content: "Pressione `a` para adicionar um bloco, `e` para editar."},
				{ID: newID(), Type: BlockDivider},
				{ID: newID(), Type: BlockTodo, Content: "Explorar o app"},
				{ID: newID(), Type: BlockTodo, Content: "Criar minha primeira nota", Done: true},
				{ID: newID(), Type: BlockCard, Content: "## Integrações\n\nConfigure `~/.config/tauos/config.json` com seus tokens do **Notion** e **GitHub** para enviar blocos com `ctrl+n` e `ctrl+g`."},
			},
		},
		{
			ID:    newID(),
			Title: "Ideias",
			Blocks: []Block{
				{ID: newID(), Type: BlockHeading, Content: "Minhas Ideias"},
				{ID: newID(), Type: BlockText, Content: "Escreva aqui suas ideias em *markdown*..."},
				{ID: newID(), Type: BlockCode, Content: "// Exemplo de código\nfmt.Println(\"Hello, Tauos!\")"},
			},
		},
	}
}
