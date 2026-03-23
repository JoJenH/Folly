package views

import (
	"os"
	"strings"
	"testing"

	"balatro-cli/game"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestMain(m *testing.M) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	os.Exit(m.Run())
}

func TestRenderHandShowsAllCards(t *testing.T) {
	state := game.NewGame(42)
	hand := state.Run.Round.Hand
	if len(hand) != 8 {
		t.Fatalf("expected 8 cards, got %d", len(hand))
	}
	out := RenderHand(state, 0, nil)
	for _, c := range hand {
		if !strings.Contains(out, c.Rank.String()) {
			t.Errorf("card %s not found in rendered output", c.String())
		}
	}
}

func TestRenderHandCursorStyle(t *testing.T) {
	state := game.NewGame(42)
	out0 := RenderHand(state, 0, nil)
	out1 := RenderHand(state, 1, nil)
	// cursor at 0 vs cursor at 1 should produce different output
	if out0 == out1 {
		t.Error("output should differ when cursor position changes")
	}
	// cursor render should contain ANSI escape sequences (color enabled)
	if !strings.Contains(out0, "\x1b[") {
		t.Error("expected ANSI escape sequences in cursor-styled output")
	}
}

func TestRenderHandSelectedStyle(t *testing.T) {
	state := game.NewGame(42)
	noSel := RenderHand(state, 0, nil)
	withSel := RenderHand(state, 0, []int{0})
	if noSel == withSel {
		t.Error("output should differ when card is selected")
	}
}

func TestRenderHandStatusBar(t *testing.T) {
	state := game.NewGame(42)
	out := RenderHand(state, 0, nil)
	r := state.Run
	if !strings.Contains(out, "Ante") {
		t.Error("status bar missing Ante")
	}
	_ = r.Gold
	if !strings.Contains(out, "$") {
		t.Error("status bar missing gold indicator")
	}
	if !strings.Contains(out, "出牌") {
		t.Error("status bar missing hands left")
	}
	if !strings.Contains(out, "弃牌") {
		t.Error("status bar missing discards left")
	}
}
