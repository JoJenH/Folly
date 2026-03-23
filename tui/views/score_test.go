package views

import (
	"strings"
	"testing"

	"balatro-cli/engine"
)

func TestRenderScoreStepsCount(t *testing.T) {
	result := engine.ScoreResult{
		HandType:   engine.OnePair,
		FinalChips: 32,
		FinalMult:  2,
		Total:      64,
		Steps: []engine.ScoreStep{
			{Description: "一对（基础 Chips+10, Mult×2）", ChipsAfter: 10, MultAfter: 2},
			{Description: "[A♠] +11 chips", ChipsAfter: 21, MultAfter: 2},
			{Description: "[A♥] +11 chips", ChipsAfter: 32, MultAfter: 2},
		},
	}
	out := RenderScore(result, 100)
	for _, step := range result.Steps {
		if !strings.Contains(out, step.Description) {
			t.Errorf("step %q not found in output", step.Description)
		}
	}
}

func TestRenderScoreTotalLine(t *testing.T) {
	result := engine.ScoreResult{
		HandType:   engine.OnePair,
		FinalChips: 32,
		FinalMult:  2,
		Total:      64,
		Steps:      []engine.ScoreStep{},
	}
	out := RenderScore(result, 100)
	if !strings.Contains(out, "64") {
		t.Error("total score not found in output")
	}
	if !strings.Contains(out, "Total") && !strings.Contains(out, "总分") {
		t.Error("total label not found in output")
	}
}

func TestRenderScoreVictoryMark(t *testing.T) {
	result := engine.ScoreResult{Total: 100}
	out := RenderScore(result, 99) // score > target
	if !strings.Contains(out, "✓") {
		t.Error("expected ✓ mark when score exceeds target")
	}
}

func TestRenderScoreDefeatMark(t *testing.T) {
	result := engine.ScoreResult{Total: 50}
	out := RenderScore(result, 100) // score < target
	if !strings.Contains(out, "✗") {
		t.Error("expected ✗ mark when score is below target")
	}
}
