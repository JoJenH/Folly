package tui

import (
	"testing"

	"balatro-cli/game"
	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel() Model {
	return NewModel(game.NewGame(42))
}

func sendKey(m Model, key tea.KeyMsg) (Model, tea.Cmd) {
	updated, cmd := m.Update(key)
	return updated.(Model), cmd
}

func TestInitialViewIsHand(t *testing.T) {
	m := newTestModel()
	if m.view != ViewHand {
		t.Errorf("expected ViewHand, got %v", m.view)
	}
}

func TestQuestionMarkSwitchesToHelp(t *testing.T) {
	m := newTestModel()
	m, _ = sendKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if m.view != ViewHelp {
		t.Errorf("expected ViewHelp after ?, got %v", m.view)
	}
}

func TestAnyKeyFromHelpReturnsToHand(t *testing.T) {
	m := newTestModel()
	m, _ = sendKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m, _ = sendKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if m.view != ViewHand {
		t.Errorf("expected ViewHand after any key from help, got %v", m.view)
	}
}

func TestPlayWithNoSelectionShowsWarning(t *testing.T) {
	m := newTestModel()
	if len(m.selected) != 0 {
		t.Fatal("expected no selected cards initially")
	}
	m, _ = sendKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if m.view != ViewHand {
		t.Errorf("expected to stay on ViewHand, got %v", m.view)
	}
	if m.lastWarning == "" {
		t.Error("expected a warning message when playing with no selection")
	}
}

func TestFirstCtrlCSetsPendingQuit(t *testing.T) {
	m := newTestModel()
	m, _ = sendKey(m, tea.KeyMsg{Type: tea.KeyCtrlC})
	if !m.pendingQuit {
		t.Error("expected pendingQuit=true after first Ctrl+C")
	}
}

func TestSecondCtrlCQuitsWhenPending(t *testing.T) {
	m := newTestModel()
	m.pendingQuit = true
	_, cmd := sendKey(m, tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("expected a quit command on second Ctrl+C")
	}
	// tea.Quit is a function; verify cmd is not nil (it returns tea.Quit msg)
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", msg)
	}
}
