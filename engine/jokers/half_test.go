package jokers

import (
	"balatro-cli/engine"
	"testing"
)

func TestHalfJokerThreeCards(t *testing.T) {
	// 打出 3 张牌时 Mult +20
	cards := []engine.Card{
		{Suit: engine.SuitSpade, Rank: engine.Rank7},
		{Suit: engine.SuitHeart, Rank: engine.Rank7},
		{Suit: engine.SuitDiamond, Rank: engine.Rank7},
	}
	half, ok := engine.DefaultRegistry.ByID("half")
	if !ok {
		t.Fatal("half joker not registered")
	}
	jokers := []engine.JokerHook{half}
	result := engine.ScoreHand(cards, jokers)
	// ThreeOfAKind: base mult=3, +20 → FinalMult=23
	if result.FinalMult != 23 {
		t.Errorf("3 cards half joker: FinalMult = %d, want 23", result.FinalMult)
	}
}

func TestHalfJokerFiveCards(t *testing.T) {
	// 打出 5 张牌时 Mult 不变
	cards := []engine.Card{
		{Suit: engine.SuitSpade, Rank: engine.RankAce},
		{Suit: engine.SuitHeart, Rank: engine.RankKing},
		{Suit: engine.SuitDiamond, Rank: engine.RankQueen},
		{Suit: engine.SuitClub, Rank: engine.RankJack},
		{Suit: engine.SuitSpade, Rank: engine.Rank9},
	}
	half, ok := engine.DefaultRegistry.ByID("half")
	if !ok {
		t.Fatal("half joker not registered")
	}
	jokers := []engine.JokerHook{half}
	result := engine.ScoreHand(cards, jokers)
	// HighCard: base mult=1，5 张不触发 half
	if result.FinalMult != 1 {
		t.Errorf("5 cards half joker: FinalMult = %d, want 1", result.FinalMult)
	}
}

func TestHalfJokerOneCard(t *testing.T) {
	// 打出 1 张牌时 Mult +20
	cards := []engine.Card{
		{Suit: engine.SuitSpade, Rank: engine.RankAce},
	}
	half, ok := engine.DefaultRegistry.ByID("half")
	if !ok {
		t.Fatal("half joker not registered")
	}
	jokers := []engine.JokerHook{half}
	result := engine.ScoreHand(cards, jokers)
	// HighCard: base mult=1, +20 → FinalMult=21
	if result.FinalMult != 21 {
		t.Errorf("1 card half joker: FinalMult = %d, want 21", result.FinalMult)
	}
}
