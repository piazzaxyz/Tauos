package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focus int

const (
	FocusSidebar Focus = iota
	FocusContent
)

type AppState int

const (
	StateNormal          AppState = iota
	StateNewPage                  // typing new page title
	StateEditBlock                // editing existing block content
	StateNewBlock                 // typing new block content (type already chosen)
	StateChooseBlockType          // waiting for t/h/c key
)

type Model struct {
	pages        []Page
	currentPage  int
	currentBlock int
	sideOffset   int // sidebar scroll offset
	contOffset   int // content scroll offset

	focus Focus
	state AppState

	newBlockType BlockType
	input        textinput.Model

	width  int
	height int
}

func newInput() textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	ti.CursorStyle = lipgloss.NewStyle().Foreground(colorAccent)
	return ti
}

func initialModel() Model {
	pages := loadPages()
	return Model{
		pages: pages,
		input: newInput(),
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
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateNewPage:
			return m.updateNewPage(msg)
		case StateEditBlock, StateNewBlock:
			return m.updateEditBlock(msg)
		case StateChooseBlockType:
			return m.updateChooseBlockType(msg)
		default:
			return m.updateNormal(msg)
		}
	}
	return m, nil
}

// ─── Normal navigation ────────────────────────────────────────────────────────

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	case "enter":
		if m.focus == FocusSidebar {
			m.focus = FocusContent
		} else {
			m.activateBlock()
		}

	case "n":
		// New page
		m.input = newInput()
		m.input.Placeholder = "Título da página..."
		m.input.Focus()
		m.state = StateNewPage

	case "N":
		// New page (same)
		m.input = newInput()
		m.input.Placeholder = "Título da página..."
		m.input.Focus()
		m.state = StateNewPage

	case "a":
		// Add block
		if len(m.pages) > 0 {
			m.state = StateChooseBlockType
		}

	case "e":
		// Edit current block
		if m.focus == FocusContent && len(m.pages) > 0 {
			m.enterEditBlock()
		}

	case "d", "backspace":
		if m.focus == FocusSidebar {
			m.deleteCurrentPage()
		} else if m.focus == FocusContent {
			m.deleteCurrentBlock()
		}

	case "D":
		// Delete page regardless of focus
		m.deleteCurrentPage()

	case "esc":
		m.focus = FocusSidebar
	}

	return m, nil
}

// ─── New page input ───────────────────────────────────────────────────────────

func (m Model) updateNewPage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		title := m.input.Value()
		if title == "" {
			title = "Nova Página"
		}
		page := Page{
			ID:    newID(),
			Title: title,
			Blocks: []Block{
				{ID: newID(), Type: BlockText, Content: ""},
			},
		}
		m.pages = append(m.pages, page)
		m.currentPage = len(m.pages) - 1
		m.currentBlock = 0
		m.contOffset = 0
		m.focus = FocusContent
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

// ─── Edit block input ─────────────────────────────────────────────────────────

func (m Model) updateEditBlock(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		content := m.input.Value()
		page := &m.pages[m.currentPage]
		if m.state == StateNewBlock {
			if content == "" {
				// remove the placeholder empty block we added
				if len(page.Blocks) > 0 {
					page.Blocks = page.Blocks[:len(page.Blocks)-1]
				}
			} else {
				page.Blocks[m.currentBlock].Content = content
			}
		} else {
			if m.currentBlock < len(page.Blocks) {
				page.Blocks[m.currentBlock].Content = content
			}
		}
		m.state = StateNormal
		savePages(m.pages)
		return m, nil

	case "esc":
		if m.state == StateNewBlock {
			// remove the placeholder block
			page := &m.pages[m.currentPage]
			if len(page.Blocks) > 0 {
				page.Blocks = page.Blocks[:len(page.Blocks)-1]
				if m.currentBlock >= len(page.Blocks) {
					m.currentBlock = max(0, len(page.Blocks)-1)
				}
			}
		}
		m.state = StateNormal
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// ─── Choose block type ────────────────────────────────────────────────────────

func (m Model) updateChooseBlockType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "t":
		m.addNewBlock(BlockText)
	case "h":
		m.addNewBlock(BlockHeading)
	case "c":
		m.addNewBlock(BlockTodo)
	case "esc", "q":
		m.state = StateNormal
	}
	return m, nil
}

func (m *Model) addNewBlock(bt BlockType) {
	page := &m.pages[m.currentPage]
	insertAt := m.currentBlock + 1
	if insertAt > len(page.Blocks) {
		insertAt = len(page.Blocks)
	}

	newBlock := Block{
		ID:   newID(),
		Type: bt,
	}

	// Insert at position
	page.Blocks = append(page.Blocks, Block{})
	copy(page.Blocks[insertAt+1:], page.Blocks[insertAt:])
	page.Blocks[insertAt] = newBlock

	m.currentBlock = insertAt
	m.newBlockType = bt

	m.input = newInput()
	switch bt {
	case BlockHeading:
		m.input.Placeholder = "Título..."
	case BlockTodo:
		m.input.Placeholder = "Tarefa..."
	default:
		m.input.Placeholder = "Escreva algo..."
	}
	m.input.Focus()
	m.state = StateNewBlock
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
	} else {
		m.enterEditBlock()
	}
}

func (m *Model) enterEditBlock() {
	if len(m.pages) == 0 {
		return
	}
	page := &m.pages[m.currentPage]
	if m.currentBlock >= len(page.Blocks) {
		return
	}
	block := page.Blocks[m.currentBlock]
	m.input = newInput()
	m.input.SetValue(block.Content)
	m.input.CursorEnd()
	m.input.Focus()
	m.state = StateEditBlock
}

// ─── Navigation helpers ───────────────────────────────────────────────────────

func (m *Model) moveSidebarUp() {
	if m.currentPage > 0 {
		m.currentPage--
		m.currentBlock = 0
		m.contOffset = 0
		if m.currentPage < m.sideOffset {
			m.sideOffset = m.currentPage
		}
	}
}

func (m *Model) moveSidebarDown() {
	if m.currentPage < len(m.pages)-1 {
		m.currentPage++
		m.currentBlock = 0
		m.contOffset = 0
		visibleLines := m.height - 6 // approx header+status
		if m.currentPage >= m.sideOffset+visibleLines {
			m.sideOffset = m.currentPage - visibleLines + 1
		}
	}
}

func (m *Model) moveBlockUp() {
	if m.currentBlock > 0 {
		m.currentBlock--
		if m.currentBlock < m.contOffset {
			m.contOffset = m.currentBlock
		}
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
	m.contOffset = 0
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
