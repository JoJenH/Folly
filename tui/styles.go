package tui

import "github.com/charmbracelet/lipgloss"

const (
	ColorCursor   = lipgloss.Color("#00FFFF")
	ColorSelected = lipgloss.Color("#FFD700")
	ColorRed      = lipgloss.Color("#FF4444")
	ColorJoker    = lipgloss.Color("#CC44FF")
	ColorScore    = lipgloss.Color("#44FF88")
	ColorGold     = lipgloss.Color("#FFD700")
	ColorWarning  = lipgloss.Color("#8B0000")
)

var (
	CursorStyle   = lipgloss.NewStyle().Foreground(ColorCursor).Bold(true)
	SelectedStyle = lipgloss.NewStyle().Foreground(ColorSelected).Bold(true)
	RedSuitStyle  = lipgloss.NewStyle().Foreground(ColorRed)
	JokerStyle    = lipgloss.NewStyle().Foreground(ColorJoker)
	ScoreStyle    = lipgloss.NewStyle().Foreground(ColorScore)
	GoldStyle     = lipgloss.NewStyle().Foreground(ColorGold)
	WarningStyle  = lipgloss.NewStyle().Foreground(ColorWarning)
)
