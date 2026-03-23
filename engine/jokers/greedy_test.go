package jokers

import (
	"balatro-cli/engine"
	"testing"
)

func TestGreedyJokerDiamond(t *testing.T) {
	// 打出 [A♦]，OnCardScored 触发，Mult +3
	cards := []engine.Card{
		{Suit: engine.SuitDiamond, Rank: engine.RankAce},
	}
	greedy, ok := engine.DefaultRegistry.ByID("greedy")
	if !ok {
		t.Fatal("greedy joker not registered")
	}
	jokers := []engine.JokerHook{greedy}
	result := engine.ScoreHand(cards, jokers)
	// HighCard A: base chips=5, base mult=1
	// OnCardScored: mult+3 → FinalMult=4
	// FinalChips = 5+11 = 16, Total = 16*4 = 64
	if result.FinalMult != 4 {
		t.Errorf("[A♦] with greedy: FinalMult = %d, want 4", result.FinalMult)
	}
}

func TestGreedyJokerNonDiamond(t *testing.T) {
	// 打出 [A♠]，Mult 不变
	cards := []engine.Card{
		{Suit: engine.SuitSpade, Rank: engine.RankAce},
	}
	greedy, ok := engine.DefaultRegistry.ByID("greedy")
	if !ok {
		t.Fatal("greedy joker not registered")
	}
	jokers := []engine.JokerHook{greedy}
	result := engine.ScoreHand(cards, jokers)
	// HighCard A: base mult=1，无触发
	if result.FinalMult != 1 {
		t.Errorf("[A♠] with greedy: FinalMult = %d, want 1", result.FinalMult)
	}
}

func TestGreedyJokerThreeDiamonds(t *testing.T) {
	// 打出 3 张♦，Mult +9
	cards := []engine.Card{
		{Suit: engine.SuitDiamond, Rank: engine.RankAce},
		{Suit: engine.SuitDiamond, Rank: engine.RankKing},
		{Suit: engine.SuitDiamond, Rank: engine.RankQueen},
	}
	greedy, ok := engine.DefaultRegistry.ByID("greedy")
	if !ok {
		t.Fatal("greedy joker not registered")
	}
	jokers := []engine.JokerHook{greedy}
	result := engine.ScoreHand(cards, jokers)
	// ThreeOfAKind? No — 三张不同点数 → HighCard
	// HighCard: ScoringCards=[A♦], base mult=1, +3 per diamond scoring card
	// OnCardScored for A♦: mult+3; K♦ and Q♦ are kickers, not scored
	// FinalMult = 1+3 = 4
	if result.FinalMult != 4 {
		t.Errorf("3 diamonds (diff rank) greedy FinalMult = %d, want 4", result.FinalMult)
	}
}

func TestGreedyJokerThreeDiamondsOnePair(t *testing.T) {
	// 3 张同点数♦构成三条，3 张都是 ScoringCards，Mult +9
	cards := []engine.Card{
		{Suit: engine.SuitDiamond, Rank: engine.Rank7},
		{Suit: engine.SuitDiamond, Rank: engine.Rank7},
		{Suit: engine.SuitDiamond, Rank: engine.Rank7},
	}
	greedy, ok := engine.DefaultRegistry.ByID("greedy")
	if !ok {
		t.Fatal("greedy joker not registered")
	}
	jokers := []engine.JokerHook{greedy}
	result := engine.ScoreHand(cards, jokers)
	// ThreeOfAKind: base mult=3, 3 scoring cards each +3 → FinalMult = 3+9 = 12
	if result.FinalMult != 12 {
		t.Errorf("3x7♦ three-of-a-kind greedy FinalMult = %d, want 12", result.FinalMult)
	}
}
