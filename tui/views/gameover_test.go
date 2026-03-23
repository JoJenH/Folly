package views

import (
	"strings"
	"testing"

	"balatro-cli/game"
)

func TestRenderGameOverVictory(t *testing.T) {
	state := game.NewGame(12345)
	out := RenderGameOver(state, true)
	if !strings.Contains(out, "胜利") && !strings.Contains(out, "Victory") {
		t.Error("victory text not found in game over screen")
	}
}

func TestRenderGameOverDefeat(t *testing.T) {
	state := game.NewGame(12345)
	out := RenderGameOver(state, false)
	if !strings.Contains(out, "失败") && !strings.Contains(out, "Game Over") {
		t.Error("defeat text not found in game over screen")
	}
}

func TestRenderGameOverShowsAnte(t *testing.T) {
	state := game.NewGame(42)
	out := RenderGameOver(state, false)
	if !strings.Contains(out, "Ante") {
		t.Error("ante not shown in game over screen")
	}
}

func TestRenderGameOverShowsSeed(t *testing.T) {
	state := game.NewGame(99999)
	out := RenderGameOver(state, false)
	if !strings.Contains(out, "99999") {
		t.Error("seed not shown in game over screen")
	}
}
