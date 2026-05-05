package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tauos/integrations"
)

type Focus int

const (
	FocusSidebar Focus = iota
	FocusContent
)

type AppState int

const (
	StateNormal          AppState = iota
	StateNewPage                  // textinput: nova página
	StateRename                   // textinput: renomear
	StateChooseBlockType          // escolher tipo de bloco
	StateEditBlock                // textarea: editar bloco existente
	StateNewBlock                 // textarea: novo bloco
)

type notionDoneMsg struct{ url string }
type githubDoneMsg struct{ url string }
type errMsg struct{ err error }

type Model struct {
	pages        []Page
	currentPage  int
	currentBlock int
	sideOffset   int

	focus Focus
	state AppState

	newBlockType BlockType
	input        textinput.Model
	textarea     textarea.Model

	cfg       Config
	statusMsg string
	width     int
	height    int
}

func newInput() textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	ti.CursorStyle = lipgloss.NewStyle().Foreground(colorAccent)
	return ti
}

func newTextarea() textarea.Model {
	ta := textarea.New()
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.SetHeight(6)
	ta.CharLimit = 0
	ta.FocusedStyle.Base = lipgloss.NewStyle().Foreground(colorText)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(colorOverlay)
	ta.BlurredStyle.Base = lipgloss.NewStyle().Foreground(colorText)
	return ta
}

func initialModel() Model {
	cfg := loadConfig()
	saveConfigTemplate()
	return Model{
		pages:    loadPages(),
		input:    newInput(),
		textarea: newTextarea(),
		cfg:      cfg,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(m.width - sidebarWidth - 10)
		return m, nil

	case notionDoneMsg:
		m.statusMsg = successStyle.Render("✓ Notion: " + msg.url)
		return m, nil

	case githubDoneMsg:
		m.statusMsg = successStyle.Render("✓ GitHub: " + msg.url)
		return m, nil

	case errMsg:
		m.statusMsg = errorStyle.Render("✗ " + msg.err.Error())
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateNewPage, StateRename:
			return m.updateTextInput(msg)
		case StateNewBlock, StateEditBlock:
			return m.updateTextarea(msg)
		case StateChooseBlockType:
			return m.updateChooseBlockType(msg)
		default:
			return m.updateNormal(msg)
		}
	}

	// Forward to textarea when active
	if m.state == StateNewBlock || m.state == StateEditBlock {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

// ─── Normal navigation ────────────────────────────────────────────────────────

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.statusMsg = ""
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		if m.focus == FocusSidebar {
			m.focus = FocusContent
		} else {
			m.focus = FocusSidebar
		}

	case "left", "h":
		m.focus = FocusSidebar

	case "right", "l":
		if len(m.pages) > 0 {
			m.focus = FocusContent
		}

	case "up", "k":
		if m.focus == FocusSidebar {
			m.moveSidebarUp()
		} else {
			m.moveBlockUp()
		}

	case "down", "j":
		if m.focus == FocusSidebar {
			m.moveSidebarDown()
		} else {
			m.moveBlockDown()
		}

	case "K":
		if m.focus == FocusContent {
			m.swapBlockUp()
		}

	case "J":
		if m.focus == FocusContent {
			m.swapBlockDown()
		}

	case "enter":
		if m.focus == FocusSidebar {
			if len(m.pages) > 0 {
				m.focus = FocusContent
			}
		} else {
			m.activateBlock()
		}

	case "n":
		m.input = newInput()
		m.input.Placeholder = "Título da página..."
		m.input.Focus()
		m.state = StateNewPage

	case "r":
		if len(m.pages) > 0 {
			m.input = newInput()
			m.input.SetValue(m.pages[m.currentPage].Title)
			m.input.CursorEnd()
			m.input.Focus()
			m.state = StateRename
		}

	case "a":
		if len(m.pages) > 0 {
			m.state = StateChooseBlockType
		}

	case "e":
		if m.focus == FocusContent && len(m.pages) > 0 {
			return m.openBlockEditor()
		}

	case "d", "backspace":
		if m.focus == FocusSidebar {
			m.deleteCurrentPage()
		} else if m.focus == FocusContent {
			m.deleteCurrentBlock()
		}

	case "D":
		m.deleteCurrentPage()

	case "ctrl+n":
		if m.focus == FocusContent && len(m.pages) > 0 {
			return m, m.pushToNotion()
		}

	case "ctrl+g":
		if m.focus == FocusContent && len(m.pages) > 0 {
			return m, m.pushToGitHub()
		}

	case "esc":
		m.focus = FocusSidebar
	}

	return m, nil
}

// ─── Text input (new page / rename) ──────────────────────────────────────────

func (m Model) updateTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.input.Value())
		if m.state == StateNewPage {
			if val == "" {
				val = "Nova Página"
			}
			page := Page{
				ID:    newID(),
				Title: val,
				Blocks: []Block{
					{ID: newID(), Type: BlockText, Content: ""},
				},
			}
			m.pages = append(m.pages, page)
			m.currentPage = len(m.pages) - 1
			m.currentBlock = 0
			m.focus = FocusContent
		} else {
			if val == "" {
				val = "Sem Título"
			}
			m.pages[m.currentPage].Title = val
		}
		m.state = StateNormal
		savePages(m.pages)
		return m, nil

	case "esc":
		m.state = StateNormal
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// ─── Textarea (new block / edit block) ───────────────────────────────────────

func (m Model) updateTextarea(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		content := strings.TrimSpace(m.textarea.Value())
		if m.state == StateNewBlock {
			page := &m.pages[m.currentPage]
			if content == "" {
				// Remove the placeholder block
				if m.currentBlock < len(page.Blocks) {
					page.Blocks = append(page.Blocks[:m.currentBlock], page.Blocks[m.currentBlock+1:]...)
					if m.currentBlock >= len(page.Blocks) {
						m.currentBlock = max(0, len(page.Blocks)-1)
					}
				}
			} else {
				page.Blocks[m.currentBlock].Content = content
			}
		} else { // StateEditBlock
			if len(m.pages) > 0 && m.currentBlock < len(m.pages[m.currentPage].Blocks) {
				m.pages[m.currentPage].Blocks[m.currentBlock].Content = content
			}
		}
		m.state = StateNormal
		m.textarea.Blur()
		savePages(m.pages)
		return m, nil

	case "esc":
		if m.state == StateNewBlock {
			page := &m.pages[m.currentPage]
			if m.currentBlock < len(page.Blocks) {
				page.Blocks = append(page.Blocks[:m.currentBlock], page.Blocks[m.currentBlock+1:]...)
				if m.currentBlock >= len(page.Blocks) {
					m.currentBlock = max(0, len(page.Blocks)-1)
				}
			}
		}
		m.state = StateNormal
		m.textarea.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// ─── Choose block type ────────────────────────────────────────────────────────

func (m Model) updateChooseBlockType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "t":
		return m.addNewBlock(BlockText)
	case "h":
		return m.addNewBlock(BlockHeading)
	case "c":
		return m.addNewBlock(BlockTodo)
	case "C":
		return m.addNewBlock(BlockCode)
	case "q":
		return m.addNewBlock(BlockQuote)
	case "k":
		return m.addNewBlock(BlockCard)
	case "/":
		m.insertDivider()
	case "esc":
		m.state = StateNormal
	}
	return m, nil
}

func (m *Model) insertDivider() {
	page := &m.pages[m.currentPage]
	insertAt := m.currentBlock + 1
	if insertAt > len(page.Blocks) {
		insertAt = len(page.Blocks)
	}
	page.Blocks = append(page.Blocks, Block{})
	copy(page.Blocks[insertAt+1:], page.Blocks[insertAt:])
	page.Blocks[insertAt] = Block{ID: newID(), Type: BlockDivider}
	m.currentBlock = insertAt
	m.state = StateNormal
	savePages(m.pages)
}

func (m Model) addNewBlock(bt BlockType) (Model, tea.Cmd) {
	page := &m.pages[m.currentPage]
	insertAt := m.currentBlock + 1
	if insertAt > len(page.Blocks) {
		insertAt = len(page.Blocks)
	}

	newBlock := Block{ID: newID(), Type: bt}
	page.Blocks = append(page.Blocks, Block{})
	copy(page.Blocks[insertAt+1:], page.Blocks[insertAt:])
	page.Blocks[insertAt] = newBlock
	m.currentBlock = insertAt
	m.newBlockType = bt

	if bt == BlockDivider {
		m.state = StateNormal
		savePages(m.pages)
		return m, nil
	}

	m.textarea = newTextarea()
	m.textarea.SetWidth(m.width - sidebarWidth - 10)
	var placeholder string
	switch bt {
	case BlockHeading:
		placeholder = "Título... (ctrl+s para salvar)"
	case BlockTodo:
		placeholder = "Tarefa... (ctrl+s para salvar)"
	case BlockCode:
		placeholder = "Código... (ctrl+s para salvar)"
	case BlockQuote:
		placeholder = "Citação... (ctrl+s para salvar)"
	case BlockCard:
		placeholder = "Conteúdo do card em markdown... (ctrl+s para salvar)"
	default:
		placeholder = "Escreva em markdown... (ctrl+s para salvar)"
	}
	m.textarea.Placeholder = placeholder
	m.textarea.Focus()
	m.state = StateNewBlock
	return m, textarea.Blink
}

// ─── Block editor ─────────────────────────────────────────────────────────────

func (m Model) openBlockEditor() (Model, tea.Cmd) {
	if len(m.pages) == 0 {
		return m, nil
	}
	page := m.pages[m.currentPage]
	if m.currentBlock >= len(page.Blocks) {
		return m, nil
	}
	block := page.Blocks[m.currentBlock]
	if block.Type == BlockDivider {
		return m, nil
	}

	m.textarea = newTextarea()
	m.textarea.SetWidth(m.width - sidebarWidth - 10)
	m.textarea.SetValue(block.Content)
	m.textarea.Focus()
	m.state = StateEditBlock
	return m, textarea.Blink
}

// ─── Block activation ─────────────────────────────────────────────────────────

func (m *Model) activateBlock() {
	if len(m.pages) == 0 {
		return
	}
	page := &m.pages[m.currentPage]
	if m.currentBlock >= len(page.Blocks) {
		return
	}
	block := &page.Blocks[m.currentBlock]
	if block.Type == BlockTodo {
		block.Done = !block.Done
		savePages(m.pages)
	}
}

// ─── Integrations ─────────────────────────────────────────────────────────────

func (m Model) pushToNotion() tea.Cmd {
	if m.cfg.NotionToken == "" || m.cfg.NotionDatabaseID == "" {
		return func() tea.Msg {
			return errMsg{fmt.Errorf("configure notion_token e notion_database_id em ~/.config/tauos/config.json")}
		}
	}
	if len(m.pages) == 0 {
		return nil
	}
	page := m.pages[m.currentPage]
	var title, body string
	if m.currentBlock < len(page.Blocks) {
		b := page.Blocks[m.currentBlock]
		content := b.Content
		if len(content) > 80 {
			title = content[:80] + "..."
		} else {
			title = content
		}
		body = content
	} else {
		title = page.Title
		body = page.Title
	}
	token := m.cfg.NotionToken
	dbID := m.cfg.NotionDatabaseID
	return func() tea.Msg {
		client := integrations.NotionClient{Token: token, DatabaseID: dbID}
		url, err := client.CreateTask(title, body)
		if err != nil {
			return errMsg{err}
		}
		return notionDoneMsg{url}
	}
}

func (m Model) pushToGitHub() tea.Cmd {
	if m.cfg.GitHubToken == "" || m.cfg.GitHubOwner == "" || m.cfg.GitHubRepo == "" {
		return func() tea.Msg {
			return errMsg{fmt.Errorf("configure github_token, github_owner e github_repo em ~/.config/tauos/config.json")}
		}
	}
	if len(m.pages) == 0 {
		return nil
	}
	page := m.pages[m.currentPage]
	var title, body string
	if m.currentBlock < len(page.Blocks) {
		b := page.Blocks[m.currentBlock]
		content := b.Content
		if len(content) > 80 {
			title = content[:80] + "..."
		} else {
			title = content
		}
		body = "**Página:** " + page.Title + "\n\n" + content
	} else {
		title = page.Title
		body = page.Title
	}
	token := m.cfg.GitHubToken
	owner := m.cfg.GitHubOwner
	repo := m.cfg.GitHubRepo
	return func() tea.Msg {
		client := integrations.GitHubClient{Token: token, Owner: owner, Repo: repo}
		url, err := client.CreateIssue(title, body)
		if err != nil {
			return errMsg{err}
		}
		return githubDoneMsg{url}
	}
}

// ─── Navigation helpers ───────────────────────────────────────────────────────

func (m *Model) moveSidebarUp() {
	if m.currentPage > 0 {
		m.currentPage--
		m.currentBlock = 0
		if m.currentPage < m.sideOffset {
			m.sideOffset = m.currentPage
		}
	}
}

func (m *Model) moveSidebarDown() {
	if m.currentPage < len(m.pages)-1 {
		m.currentPage++
		m.currentBlock = 0
		visibleLines := m.height - 6
		if visibleLines < 1 {
			visibleLines = 1
		}
		if m.currentPage >= m.sideOffset+visibleLines {
			m.sideOffset = m.currentPage - visibleLines + 1
		}
	}
}

func (m *Model) moveBlockUp() {
	if m.currentBlock > 0 {
		m.currentBlock--
	}
}

func (m *Model) moveBlockDown() {
	if len(m.pages) == 0 {
		return
	}
	page := m.pages[m.currentPage]
	if m.currentBlock < len(page.Blocks)-1 {
		m.currentBlock++
	}
}

func (m *Model) swapBlockUp() {
	if len(m.pages) == 0 || m.currentBlock <= 0 {
		return
	}
	page := &m.pages[m.currentPage]
	i := m.currentBlock
	page.Blocks[i], page.Blocks[i-1] = page.Blocks[i-1], page.Blocks[i]
	m.currentBlock--
	savePages(m.pages)
}

func (m *Model) swapBlockDown() {
	if len(m.pages) == 0 {
		return
	}
	page := &m.pages[m.currentPage]
	i := m.currentBlock
	if i >= len(page.Blocks)-1 {
		return
	}
	page.Blocks[i], page.Blocks[i+1] = page.Blocks[i+1], page.Blocks[i]
	m.currentBlock++
	savePages(m.pages)
}

// ─── Delete helpers ───────────────────────────────────────────────────────────

func (m *Model) deleteCurrentPage() {
	if len(m.pages) == 0 {
		return
	}
	m.pages = append(m.pages[:m.currentPage], m.pages[m.currentPage+1:]...)
	if m.currentPage >= len(m.pages) {
		m.currentPage = max(0, len(m.pages)-1)
	}
	m.currentBlock = 0
	savePages(m.pages)
}

func (m *Model) deleteCurrentBlock() {
	if len(m.pages) == 0 {
		return
	}
	page := &m.pages[m.currentPage]
	if len(page.Blocks) == 0 {
		return
	}
	page.Blocks = append(page.Blocks[:m.currentBlock], page.Blocks[m.currentBlock+1:]...)
	if m.currentBlock >= len(page.Blocks) {
		m.currentBlock = max(0, len(page.Blocks)-1)
	}
	savePages(m.pages)
}
