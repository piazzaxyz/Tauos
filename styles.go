package main

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha palette
var (
	colorBg      = lipgloss.Color("#1e1e2e")
	colorSurface = lipgloss.Color("#181825")
	colorOverlay = lipgloss.Color("#313244")
	colorText    = lipgloss.Color("#cdd6f4")
	colorSubtext = lipgloss.Color("#a6adc8")
	colorMuted   = lipgloss.Color("#6c7086")
	colorAccent  = lipgloss.Color("#cba6f7")
	colorBlue    = lipgloss.Color("#89b4fa")
	colorGreen   = lipgloss.Color("#a6e3a1")
	colorRed     = lipgloss.Color("#f38ba8")
	colorBorder  = lipgloss.Color("#45475a")

	sidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorSubtext).
				PaddingLeft(1)

	sidebarItemActiveStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				PaddingLeft(1)

	sidebarItemSelectedStyle = lipgloss.NewStyle().
					Foreground(colorAccent).
					Background(colorOverlay).
					Bold(true).
					PaddingLeft(1)

	titleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

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

	todoUncheckedStyle = lipgloss.NewStyle().
				Foreground(colorSubtext)

	todoDoneStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	checkboxOn  = lipgloss.NewStyle().Foreground(colorGreen).Render("☑")
	checkboxOff = lipgloss.NewStyle().Foreground(colorMuted).Render("☐")

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

	dimStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	_ = colorRed // used via errorStyle
)
