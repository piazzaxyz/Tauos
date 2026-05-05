package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const sidebarWidth = 28

func (m Model) View() string {
	if m.width == 0 {
		return "Carregando..."
	}

	header := m.viewHeader()
	status := m.viewStatus()

	headerH := lipgloss.Height(header)
	statusH := lipgloss.Height(status)
	bodyH := m.height - headerH - statusH
	if bodyH < 1 {
		bodyH = 1
	}

	sidebar := m.viewSidebar(bodyH)
	content := m.viewContent(bodyH)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, status)
}

// ─── Header ───────────────────────────────────────────────────────────────────

func (m Model) viewHeader() string {
	appName := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Render("  τ  tauos")

	var mid string
	if len(m.pages) > 0 {
		crumb := lipgloss.NewStyle().Foreground(colorMuted).Render("  ›  ")
		pageTitle := lipgloss.NewStyle().Foreground(colorSubtext).Render(m.pages[m.currentPage].Title)
		mid = crumb + pageTitle
	}

	date := time.Now().Format("02 Jan 2006")
	right := lipgloss.NewStyle().Foreground(colorMuted).Render(date + "  ")

	leftW := lipgloss.Width(appName + mid)
	rightW := lipgloss.Width(right)
	gap := m.width - leftW - rightW
	if gap < 0 {
		gap = 0
	}

	line := appName + mid + strings.Repeat(" ", gap) + right
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
	if m.statusMsg != "" {
		return lipgloss.NewStyle().
			Background(colorSurface).
			Width(m.width).
			Render("  " + m.statusMsg)
	}

	var hints []string

	switch m.state {
	case StateNewPage:
		hints = []string{key("enter", "criar"), key("esc", "cancelar")}
	case StateRename:
		hints = []string{key("enter", "renomear"), key("esc", "cancelar")}
	case StateNewBlock, StateEditBlock:
		hints = []string{
			key("ctrl+s", "salvar"),
			key("esc", "cancelar"),
			keyHintDesc.Render("  markdown suportado"),
		}
	case StateChooseBlockType:
		hints = []string{
			key("t", "texto"),
			key("h", "título"),
			key("c", "check"),
			key("C", "código"),
			key("q", "citação"),
			key("k", "card"),
			key("/", "divisor"),
			key("esc", "cancelar"),
		}
	default:
		if m.focus == FocusSidebar {
			hints = []string{
				key("↑↓/jk", "navegar"),
				key("enter", "abrir"),
				key("n", "nova"),
				key("r", "renomear"),
				key("d", "deletar"),
				key("tab", "conteúdo"),
				key("q", "sair"),
			}
		} else {
			hints = []string{
				key("↑↓/jk", "navegar"),
				key("enter", "toggle"),
				key("e", "editar"),
				key("a", "adicionar"),
				key("K/J", "mover"),
				key("d", "deletar"),
				key("^n", "notion"),
				key("^g", "github"),
			}
		}
	}

	bar := strings.Join(hints, dimStyle.Render("  ·  "))
	barW := lipgloss.Width(bar)
	pad := m.width - barW - 2
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
	contentW := sidebarWidth - 2 // subtract border

	var rows []string

	label := lipgloss.NewStyle().
		Foreground(colorMuted).
		Bold(true).
		Width(contentW).
		Render(" PÁGINAS")
	rows = append(rows, label)

	divider := lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(strings.Repeat("─", contentW))
	rows = append(rows, divider)

	for i := m.sideOffset; i < len(m.pages); i++ {
		title := truncate(m.pages[i].Title, contentW-4)
		var row string
		switch {
		case i == m.currentPage && m.focus == FocusSidebar:
			row = lipgloss.NewStyle().
				Background(colorAccent).
				Foreground(colorBg).
				Bold(true).
				Width(contentW).
				Render(" ▸ " + title)
		case i == m.currentPage:
			row = lipgloss.NewStyle().
				Foreground(colorAccent).
				Width(contentW).
				Render(" · " + title)
		default:
			row = sidebarItemStyle.
				Width(contentW).
				Render("   " + title)
		}
		rows = append(rows, row)
	}

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(colorBorder).
		Width(contentW).
		Render(fmt.Sprintf(" %d página(s)", len(m.pages)))
	// Pad rows to fill height minus footer (1 line)
	for len(rows) < height-1 {
		rows = append(rows, strings.Repeat(" ", contentW))
	}
	rows = append(rows, footer)

	content := strings.Join(rows, "\n")

	return lipgloss.NewStyle().
		Background(colorSurface).
		Width(sidebarWidth).
		Height(height).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(colorBorder).
		Render(content)
}

// ─── Content ──────────────────────────────────────────────────────────────────

func (m Model) viewContent(height int) string {
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
			Render("Nenhuma página\n\n" + keyHintKey.Render("n") + keyHintDesc.Render(" para criar"))
		return lipgloss.NewStyle().Background(colorBg).Render(empty)
	}

	page := m.pages[m.currentPage]
	pad := "   "
	innerW := cw - 6
	if innerW < 10 {
		innerW = 10
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, pad+titleStyle.Width(innerW).Render(page.Title))
	lines = append(lines, pad+lipgloss.NewStyle().Foreground(colorBorder).Render(strings.Repeat("─", innerW)))
	lines = append(lines, "")

	// Inline input for new page / rename
	if m.state == StateNewPage || m.state == StateRename {
		label := "Nova página: "
		if m.state == StateRename {
			label = "Renomear: "
		}
		box := lipgloss.NewStyle().
			Background(colorOverlay).
			Width(innerW).
			Padding(0, 1).
			Render(inputPromptStyle.Render(label) + m.input.View())
		lines = append(lines, pad+box, "")
	}

	// Blocks
	for i, block := range page.Blocks {
		selected := i == m.currentBlock && m.focus == FocusContent

		// Show textarea when editing this block
		if (m.state == StateNewBlock || m.state == StateEditBlock) && i == m.currentBlock {
			taStyle := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Width(innerW).
				Padding(0, 1)
			prefix := blockTypeLabel(block.Type)
			if prefix != "" {
				lines = append(lines, pad+inputPromptStyle.Render(prefix))
			}
			lines = append(lines, pad+taStyle.Render(m.textarea.View()))
			lines = append(lines, "")
			continue
		}

		blockLines := m.renderBlock(block, selected, innerW)
		lines = append(lines, blockLines...)
		lines = append(lines, "")
	}

	// Block type chooser overlay
	if m.state == StateChooseBlockType {
		chooser := lipgloss.NewStyle().
			Background(colorOverlay).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Width(innerW).
			Padding(0, 1).
			Render(
				keyHintKey.Render("t") + keyHintDesc.Render(" texto") +
					"  " + keyHintKey.Render("h") + keyHintDesc.Render(" título") +
					"  " + keyHintKey.Render("c") + keyHintDesc.Render(" check") +
					"  " + keyHintKey.Render("C") + keyHintDesc.Render(" código") +
					"  " + keyHintKey.Render("q") + keyHintDesc.Render(" citação") +
					"  " + keyHintKey.Render("k") + keyHintDesc.Render(" card") +
					"  " + keyHintKey.Render("/") + keyHintDesc.Render(" divisor") +
					"  " + dimStyle.Render("esc"),
			)
		lines = append(lines, pad+chooser, "")
	}

	body := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Background(colorBg).
		Width(cw).
		Height(height).
		Render(body)
}

// ─── Block rendering ──────────────────────────────────────────────────────────

func (m Model) renderBlock(block Block, selected bool, innerW int) []string {
	pad := "   "
	var lines []string

	switch block.Type {
	case BlockDivider:
		color := colorBorder
		if selected {
			color = colorAccent
		}
		rule := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("─", innerW))
		lines = append(lines, pad+rule)

	case BlockHeading:
		content := block.Content
		if content == "" {
			content = lipgloss.NewStyle().Foreground(colorMuted).Render("Título vazio")
		}
		prefix := lipgloss.NewStyle().Foreground(colorBlue).Render("## ")
		style := headingBlockStyle
		if selected {
			style = headingBlockSelectedStyle
		}
		rendered := renderInlineMarkdown(content, style)
		for _, seg := range splitLines(rendered) {
			lines = append(lines, pad+prefix+lipgloss.NewStyle().Width(innerW-3).Render(seg))
		}

	case BlockTodo:
		cb := checkboxOff
		textStyle := todoUncheckedStyle
		if block.Done {
			cb = checkboxOn
			textStyle = todoDoneStyle
		}
		rendered := renderInlineMarkdown(block.Content, textStyle)
		row := cb + " " + rendered
		if selected {
			row = lipgloss.NewStyle().Background(colorOverlay).Width(innerW).Render(cb + " " + block.Content)
		}
		lines = append(lines, pad+row)

	case BlockCode:
		segs := splitLines(block.Content)
		if len(segs) == 0 {
			segs = []string{""}
		}
		style := codeBlockStyle
		if selected {
			style = codeBlockSelectedStyle
		}
		barColor := colorYellow
		if selected {
			barColor = colorAccent
		}
		bar := lipgloss.NewStyle().Foreground(barColor).Render("▎")
		for _, seg := range segs {
			lines = append(lines, pad+bar+style.Width(innerW-1).Render(" "+seg))
		}

	case BlockQuote:
		segs := splitLines(block.Content)
		if len(segs) == 0 {
			segs = []string{""}
		}
		style := quoteTextStyle
		if selected {
			style = quoteTextSelectedStyle
		}
		bar := quoteBarStyle.Render("│ ")
		for _, seg := range segs {
			rendered := renderInlineMarkdown(seg, style)
			lines = append(lines, pad+bar+rendered)
		}

	case BlockCard:
		content := block.Content
		if content == "" {
			content = lipgloss.NewStyle().Foreground(colorMuted).Render("Card vazio")
		}
		style := cardStyle
		if selected {
			style = cardSelectedStyle
		}
		rendered := renderCardContent(content, innerW-4)
		lines = append(lines, pad+style.Width(innerW).Render(rendered))

	default: // BlockText
		content := block.Content
		if content == "" {
			if selected {
				lines = append(lines, pad+textBlockSelectedStyle.Width(innerW).Render(" "))
			} else {
				lines = append(lines, pad+dimStyle.Render("·"))
			}
			return lines
		}
		style := textBlockStyle
		if selected {
			style = textBlockSelectedStyle
		}
		for _, seg := range splitLines(content) {
			rendered := renderInlineMarkdown(seg, style)
			lines = append(lines, pad+rendered)
		}
	}

	return lines
}

// blockTypeLabel returns a small label shown above the textarea.
func blockTypeLabel(bt BlockType) string {
	switch bt {
	case BlockHeading:
		return "## Título"
	case BlockTodo:
		return "○ Tarefa"
	case BlockCode:
		return "▎ Código"
	case BlockQuote:
		return "│ Citação"
	case BlockCard:
		return "╭ Card"
	default:
		return ""
	}
}

// ─── Inline markdown renderer ─────────────────────────────────────────────────

var (
	reBold          = regexp.MustCompile(`\*\*(.+?)\*\*|__(.+?)__`)
	reItalic        = regexp.MustCompile(`\*([^*]+?)\*|_([^_]+?)_`)
	reStrike        = regexp.MustCompile(`~~(.+?)~~`)
	reInlineCode    = regexp.MustCompile("`([^`]+?)`")
)

// renderInlineMarkdown applies bold/italic/strike/code styling to a string,
// then wraps the result in baseStyle.
func renderInlineMarkdown(s string, baseStyle lipgloss.Style) string {
	// We build segments: each segment is (text, style)
	type segment struct {
		text  string
		style lipgloss.Style
	}

	bold := lipgloss.NewStyle().Bold(true).Foreground(baseStyle.GetForeground()).Background(baseStyle.GetBackground())
	italic := lipgloss.NewStyle().Italic(true).Foreground(baseStyle.GetForeground()).Background(baseStyle.GetBackground())
	strike := lipgloss.NewStyle().Strikethrough(true).Foreground(colorMuted).Background(baseStyle.GetBackground())

	// Process inline code first (highest priority)
	type span struct {
		start, end int
		render     string
	}

	var spans []span

	for _, m := range reInlineCode.FindAllStringSubmatchIndex(s, -1) {
		inner := s[m[2]:m[3]]
		spans = append(spans, span{m[0], m[1], inlineCodeStyle.Render(inner)})
	}
	for _, m := range reBold.FindAllStringSubmatchIndex(s, -1) {
		inner := s[m[2]:m[3]]
		if inner == "" {
			inner = s[m[4]:m[5]]
		}
		spans = append(spans, span{m[0], m[1], bold.Render(inner)})
	}
	for _, m := range reItalic.FindAllStringSubmatchIndex(s, -1) {
		inner := s[m[2]:m[3]]
		if inner == "" && m[4] >= 0 {
			inner = s[m[4]:m[5]]
		}
		spans = append(spans, span{m[0], m[1], italic.Render(inner)})
	}
	for _, m := range reStrike.FindAllStringSubmatchIndex(s, -1) {
		inner := s[m[2]:m[3]]
		spans = append(spans, span{m[0], m[1], strike.Render(inner)})
	}

	if len(spans) == 0 {
		return baseStyle.Render(s)
	}

	// Sort spans by start, resolve overlaps (first wins)
	// Simple greedy: build output left to right
	type sortedSpan struct{ start, end int; render string }
	sorted := make([]sortedSpan, 0, len(spans))
	for _, sp := range spans {
		sorted = append(sorted, sortedSpan{sp.start, sp.end, sp.render})
	}
	// bubble sort by start
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].start < sorted[i].start {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	var out strings.Builder
	pos := 0
	for _, sp := range sorted {
		if sp.start < pos {
			continue // overlapping, skip
		}
		if sp.start > pos {
			out.WriteString(baseStyle.Render(s[pos:sp.start]))
		}
		out.WriteString(sp.render)
		pos = sp.end
	}
	if pos < len(s) {
		out.WriteString(baseStyle.Render(s[pos:]))
	}
	return out.String()
}

// renderCardContent renders multiline markdown content for a card block.
func renderCardContent(content string, w int) string {
	lines := splitLines(content)
	var rendered []string
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "## "):
			text := strings.TrimPrefix(line, "## ")
			rendered = append(rendered, lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Width(w).Render(text))
		case strings.HasPrefix(line, "# "):
			text := strings.TrimPrefix(line, "# ")
			rendered = append(rendered, lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Width(w).Render(text))
		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* "):
			text := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
			bullet := lipgloss.NewStyle().Foreground(colorAccent).Render("• ")
			rendered = append(rendered, bullet+renderInlineMarkdown(text, textBlockStyle))
		case strings.HasPrefix(line, "> "):
			text := strings.TrimPrefix(line, "> ")
			bar := quoteBarStyle.Render("│ ")
			rendered = append(rendered, bar+quoteTextStyle.Render(text))
		case line == "---" || line == "___":
			rendered = append(rendered, lipgloss.NewStyle().Foreground(colorBorder).Render(strings.Repeat("─", w)))
		default:
			rendered = append(rendered, renderInlineMarkdown(line, textBlockStyle))
		}
	}
	return strings.Join(rendered, "\n")
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func splitLines(s string) []string {
	if s == "" {
		return []string{""}
	}
	return strings.Split(s, "\n")
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}
