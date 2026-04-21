package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 24

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	header := m.viewHeader()
	status := m.viewStatus()

	bodyH := m.height - lipgloss.Height(header) - lipgloss.Height(status)
	if bodyH < 1 {
		bodyH = 1
	}

	sidebar := m.viewSidebar(bodyH)
	content := m.viewContent(bodyH)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, status)
}

// ─── Header ──────────────────────────────────────────────────────────────────

func (m Model) viewHeader() string {
	left := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Render("  τ  tau")

	date := time.Now().Format("02 Jan 2006")
	right := lipgloss.NewStyle().
		Foreground(colorMuted).
		Render(date + "  ")

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	line := left + strings.Repeat(" ", gap) + right
	return lipgloss.NewStyle().
		Background(colorSurface).
		Width(m.width).
		Render(line)
}

// ─── Status bar ───────────────────────────────────────────────────────────────

func key(k, desc string) string {
	return keyHintKey.Render(k) + keyHintDesc.Render(" "+desc)
}

func (m Model) viewStatus() string {
	var hints []string

	switch m.state {
	case StateNewPage:
		hints = []string{
			key("enter", "confirmar"),
			key("esc", "cancelar"),
		}
	case StateEditBlock, StateNewBlock:
		hints = []string{
			key("enter", "salvar"),
			key("esc", "cancelar"),
		}
	case StateChooseBlockType:
		hints = []string{
			key("t", "texto"),
			key("h", "título"),
			key("c", "checkbox"),
			key("esc", "cancelar"),
		}
	default:
		if m.focus == FocusSidebar {
			hints = []string{
				key("↑↓", "navegar"),
				key("enter", "abrir"),
				key("n", "nova página"),
				key("d", "deletar"),
				key("tab", "conteúdo"),
				key("q", "sair"),
			}
		} else {
			hints = []string{
				key("↑↓", "navegar"),
				key("enter", "ativar"),
				key("e", "editar"),
				key("a", "adicionar"),
				key("d", "deletar"),
				key("tab", "páginas"),
			}
		}
	}

	bar := strings.Join(hints, dimStyle.Render("  ·  "))
	barWidth := lipgloss.Width(bar)
	pad := m.width - barWidth - 2
	if pad < 0 {
		pad = 0
	}

	return lipgloss.NewStyle().
		Background(colorSurface).
		Width(m.width).
		Render("  " + bar + strings.Repeat(" ", pad))
}

// ─── Sidebar ──────────────────────────────────────────────────────────────────

func (m Model) viewSidebar(height int) string {
	inner := sidebarWidth - 2 // account for right border char

	var sb strings.Builder

	// Section label
	label := lipgloss.NewStyle().
		Foreground(colorMuted).
		Bold(true).
		Width(inner).
		Render(" PÁGINAS")
	sb.WriteString(label + "\n")

	divider := lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(strings.Repeat("─", inner))
	sb.WriteString(divider + "\n")

	linesUsed := 2

	maxVisible := height - linesUsed
	start := m.sideOffset
	if start < 0 {
		start = 0
	}

	for i := start; i < len(m.pages) && linesUsed < height; i++ {
		page := m.pages[i]
		title := page.Title

		maxLen := inner - 3
		if len([]rune(title)) > maxLen {
			runes := []rune(title)
			title = string(runes[:maxLen-1]) + "…"
		}

		var line string
		switch {
		case i == m.currentPage && m.focus == FocusSidebar:
			line = sidebarItemSelectedStyle.Width(inner).Render("▸ " + title)
		case i == m.currentPage:
			line = sidebarItemActiveStyle.Width(inner).Render("· " + title)
		default:
			line = sidebarItemStyle.Width(inner).Render("  " + title)
		}
		sb.WriteString(line + "\n")
		linesUsed++
		_ = maxVisible
	}

	// Fill remaining height
	for linesUsed < height {
		sb.WriteString(strings.Repeat(" ", inner) + "\n")
		linesUsed++
	}

	rendered := strings.TrimRight(sb.String(), "\n")

	return lipgloss.NewStyle().
		Background(colorSurface).
		Width(sidebarWidth).
		Height(height).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Render(rendered)
}

// ─── Content ──────────────────────────────────────────────────────────────────

func (m Model) viewContent(height int) string {
	// sidebarWidth includes the border char rendered by BorderRight
	cw := m.width - sidebarWidth - 1
	if cw < 10 {
		cw = 10
	}

	if len(m.pages) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(colorMuted).
			Width(cw).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Nenhuma página\n\n" +
				keyHintKey.Render("n") + keyHintDesc.Render(" para criar uma"))
		return lipgloss.NewStyle().Background(colorBg).Render(empty)
	}

	page := m.pages[m.currentPage]
	pad := "  "
	innerW := cw - 4

	var lines []string

	// Page title
	lines = append(lines, "")
	lines = append(lines, pad+titleStyle.Width(innerW).Render(page.Title))
	lines = append(lines,
		pad+lipgloss.NewStyle().Foreground(colorBorder).Render(strings.Repeat("─", innerW)))
	lines = append(lines, "")

	// Blocks
	for i, block := range page.Blocks {
		selected := i == m.currentBlock && m.focus == FocusContent

		// Editing this block
		if (m.state == StateEditBlock || m.state == StateNewBlock) && i == m.currentBlock {
			prefix := inputPromptStyle.Render("❯ ")
			var blockTypePrefix string
			switch block.Type {
			case BlockHeading:
				blockTypePrefix = lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("## ")
			case BlockTodo:
				blockTypePrefix = checkboxOff + " "
			}
			inputLine := pad + prefix + blockTypePrefix + m.input.View()
			lines = append(lines, inputLine)
			lines = append(lines, "")
			continue
		}

		var line string
		switch block.Type {
		case BlockHeading:
			text := "## " + block.Content
			if selected {
				line = pad + headingBlockSelectedStyle.Width(innerW).Render(text)
			} else {
				line = pad + headingBlockStyle.Render(text)
			}

		case BlockTodo:
			var cb, text string
			if block.Done {
				cb = checkboxOn
				text = todoDoneStyle.Render(block.Content)
			} else {
				cb = checkboxOff
				text = todoUncheckedStyle.Render(block.Content)
			}
			content := cb + " " + text
			if selected {
				line = pad + lipgloss.NewStyle().Background(colorOverlay).Width(innerW).Render(cb+" "+block.Content)
				// re-render with proper colors
				if block.Done {
					line = pad + lipgloss.NewStyle().Background(colorOverlay).Width(innerW).Render(checkboxOn + " " + block.Content)
				} else {
					line = pad + lipgloss.NewStyle().Background(colorOverlay).Width(innerW).Render(checkboxOff + " " + block.Content)
				}
			} else {
				line = pad + content
			}

		default: // BlockText
			if block.Content == "" {
				if selected {
					line = pad + textBlockSelectedStyle.Width(innerW).Render(" ")
				} else {
					line = pad + dimStyle.Render("·")
				}
			} else {
				if selected {
					line = pad + textBlockSelectedStyle.Width(innerW).Render(block.Content)
				} else {
					line = pad + textBlockStyle.Render(block.Content)
				}
			}
		}

		lines = append(lines, line)
		lines = append(lines, "") // spacing between blocks
	}

	// Block type chooser overlay
	if m.state == StateChooseBlockType {
		lines = append(lines, "")
		chooser := lipgloss.NewStyle().
			Background(colorOverlay).
			Foreground(colorText).
			Width(innerW).
			Padding(0, 1).
			Render(
				keyHintKey.Render("t") + keyHintDesc.Render(" texto") +
					"  " + keyHintKey.Render("h") + keyHintDesc.Render(" título") +
					"  " + keyHintKey.Render("c") + keyHintDesc.Render(" checkbox") +
					"  " + dimStyle.Render("esc cancelar"),
			)
		lines = append(lines, pad+chooser)
	}

	// New page input overlay (shown in content area)
	if m.state == StateNewPage {
		lines = append(lines, "")
		inputBox := lipgloss.NewStyle().
			Background(colorOverlay).
			Width(innerW).
			Padding(0, 1).
			Render(
				inputPromptStyle.Render("Nova página: ") + m.input.View(),
			)
		lines = append(lines, pad+inputBox)
	}

	// Render with scrolling: take visible slice
	visible := lines
	if m.contOffset > 0 && m.contOffset < len(lines) {
		visible = lines[m.contOffset:]
	}

	body := strings.Join(visible, "\n")

	return lipgloss.NewStyle().
		Background(colorBg).
		Width(cw).
		Height(height).
		Render(body)
}
