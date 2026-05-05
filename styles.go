package main

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha
var (
	colorBg      = lipgloss.Color("#1e1e2e")
	colorSurface = lipgloss.Color("#181825")
	colorOverlay = lipgloss.Color("#313244")
	colorText    = lipgloss.Color("#cdd6f4")
	colorSubtext = lipgloss.Color("#a6adc8")
	colorMuted   = lipgloss.Color("#6c7086")
	colorAccent  = lipgloss.Color("#cba6f7") // mauve
	colorBlue    = lipgloss.Color("#89b4fa")
	colorGreen   = lipgloss.Color("#a6e3a1")
	colorRed     = lipgloss.Color("#f38ba8")
	colorYellow  = lipgloss.Color("#f9e2af")
	colorTeal    = lipgloss.Color("#94e2d5")
	colorPeach   = lipgloss.Color("#fab387")
	colorBorder  = lipgloss.Color("#45475a")
	colorCodeBg  = lipgloss.Color("#1e1e2e")
	colorCardBg  = lipgloss.Color("#1e1e2e")

	sidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorSubtext)

	sidebarItemActiveStyle = lipgloss.NewStyle().
				Foreground(colorAccent)

	sidebarItemSelectedStyle = lipgloss.NewStyle().
					Foreground(colorBg).
					Background(colorAccent).
					Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			MarginBottom(0)

	headingBlockStyle = lipgloss.NewStyle().
				Foreground(colorBlue).
				Bold(true)

	headingBlockSelectedStyle = lipgloss.NewStyle().
					Foreground(colorBlue).
					Background(colorOverlay).
					Bold(true)

	textBlockStyle = lipgloss.NewStyle().
			Foreground(colorText)

	textBlockSelectedStyle = lipgloss.NewStyle().
				Foreground(colorText).
				Background(colorOverlay)

	codeBlockStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Background(colorCodeBg)

	codeBlockSelectedStyle = lipgloss.NewStyle().
				Foreground(colorYellow).
				Background(colorOverlay)

	quoteBarStyle = lipgloss.NewStyle().
			Foreground(colorTeal)

	quoteTextStyle = lipgloss.NewStyle().
			Foreground(colorTeal).
			Italic(true)

	quoteTextSelectedStyle = lipgloss.NewStyle().
				Foreground(colorTeal).
				Background(colorOverlay).
				Italic(true)

	cardStyle = lipgloss.NewStyle().
			Background(colorCardBg).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	cardSelectedStyle = lipgloss.NewStyle().
				Background(colorCardBg).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)

	todoUncheckedStyle = lipgloss.NewStyle().
				Foreground(colorSubtext)

	todoDoneStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	checkboxOn  = lipgloss.NewStyle().Foreground(colorGreen).Render("●")
	checkboxOff = lipgloss.NewStyle().Foreground(colorMuted).Render("○")

	keyHintKey = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	keyHintDesc = lipgloss.NewStyle().
			Foreground(colorMuted)

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	successStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	inlineCodeStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Background(colorSurface)

	_ = colorPeach // reserved
)
