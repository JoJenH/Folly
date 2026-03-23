package tui

import (
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestMain(m *testing.M) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	os.Exit(m.Run())
}

func TestColorConstantsNonEmpty(t *testing.T) {
	colors := []lipgloss.Color{
		ColorCursor, ColorSelected, ColorRed,
		ColorJoker, ColorScore, ColorGold, ColorWarning,
	}
	for _, c := range colors {
		if string(c) == "" {
			t.Errorf("color constant is empty")
		}
	}
}

func TestStylesDistinguishable(t *testing.T) {
	cursor := CursorStyle.Render("X")
	selected := SelectedStyle.Render("X")
	red := RedSuitStyle.Render("X")

	if cursor == selected {
		t.Error("CursorStyle and SelectedStyle should be distinguishable")
	}
	if cursor == red {
		t.Error("CursorStyle and RedSuitStyle should be distinguishable")
	}
	if selected == red {
		t.Error("SelectedStyle and RedSuitStyle should be distinguishable")
	}
}
