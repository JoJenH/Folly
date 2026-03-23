package jokers

import (
	"balatro-cli/engine"
	"testing"
)

func TestRetriggerLastCard(t *testing.T) {
	// 单张 A，Retrigger Joker 让最后一张 ScoringCard 再触发一次
	// HighCard: ScoringCards=[A♠], base chips=5, base mult=1
	// 正常: chips = 5+11 = 16, total = 16*1 = 16
	// Retrigger: A♠ 额外触发一次 → chips = 5+11+11 = 27, total = 27*1 = 27
	cards := []engine.Card{
		{Suit: engine.SuitSpade, Rank: engine.RankAce},
	}
	retrig, ok := engine.DefaultRegistry.ByID("retrigger")
	if !ok {
		t.Fatal("retrigger joker not registered")
	}
	jokers := []engine.JokerHook{retrig}
	result := engine.ScoreHand(cards, jokers)
	if result.FinalChips != 27 {
		t.Errorf("Retrigger FinalChips = %d, want 27", result.FinalChips)
	}
}

func TestRetriggerStepsContainRetrigger(t *testing.T) {
	cards := []engine.Card{
		{Suit: engine.SuitSpade, Rank: engine.RankAce},
	}
	retrig, ok := engine.DefaultRegistry.ByID("retrigger")
	if !ok {
		t.Fatal("retrigger joker not registered")
	}
	jokers := []engine.JokerHook{retrig}
	result := engine.ScoreHand(cards, jokers)
	// Steps: 1 hand type + 1 retrigger + 1 normal card = 3
	if len(result.Steps) != 3 {
		t.Errorf("Retrigger Steps = %d, want 3", len(result.Steps))
	}
}
