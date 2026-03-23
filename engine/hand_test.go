package engine

import "testing"

// 辅助：快速构造牌
func card(suit Suit, rank Rank) Card {
	return Card{Suit: suit, Rank: rank}
}

func TestEvaluateStraightFlush(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitSpade, RankKing),
		card(SuitSpade, RankQueen), card(SuitSpade, RankJack), card(SuitSpade, Rank10),
	}
	r := Evaluate(cards)
	if r.Type != StraightFlush {
		t.Errorf("expected StraightFlush, got %v", r.Type)
	}
	if len(r.ScoringCards) != 5 {
		t.Errorf("StraightFlush ScoringCards = %d, want 5", len(r.ScoringCards))
	}
}

func TestEvaluateFourOfAKind(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankAce),
		card(SuitDiamond, RankAce), card(SuitClub, RankAce), card(SuitSpade, Rank2),
	}
	r := Evaluate(cards)
	if r.Type != FourOfAKind {
		t.Errorf("expected FourOfAKind, got %v", r.Type)
	}
	if len(r.ScoringCards) != 4 {
		t.Errorf("FourOfAKind ScoringCards = %d, want 4", len(r.ScoringCards))
	}
}

func TestEvaluateFullHouse(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankKing), card(SuitHeart, RankKing), card(SuitDiamond, RankKing),
		card(SuitClub, RankAce), card(SuitSpade, RankAce),
	}
	r := Evaluate(cards)
	if r.Type != FullHouse {
		t.Errorf("expected FullHouse, got %v", r.Type)
	}
	if len(r.ScoringCards) != 5 {
		t.Errorf("FullHouse ScoringCards = %d, want 5", len(r.ScoringCards))
	}
}

func TestEvaluateFlush(t *testing.T) {
	cards := []Card{
		card(SuitSpade, Rank2), card(SuitSpade, Rank5),
		card(SuitSpade, Rank7), card(SuitSpade, Rank9), card(SuitSpade, RankJack),
	}
	r := Evaluate(cards)
	if r.Type != Flush {
		t.Errorf("expected Flush, got %v", r.Type)
	}
	if len(r.ScoringCards) != 5 {
		t.Errorf("Flush ScoringCards = %d, want 5", len(r.ScoringCards))
	}
}

func TestEvaluateStraightAceLow(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, Rank2),
		card(SuitDiamond, Rank3), card(SuitClub, Rank4), card(SuitSpade, Rank5),
	}
	r := Evaluate(cards)
	if r.Type != Straight {
		t.Errorf("expected Straight (A-low), got %v", r.Type)
	}
}

func TestEvaluateStraightAceHigh(t *testing.T) {
	cards := []Card{
		card(SuitSpade, Rank10), card(SuitHeart, RankJack),
		card(SuitDiamond, RankQueen), card(SuitClub, RankKing), card(SuitSpade, RankAce),
	}
	r := Evaluate(cards)
	if r.Type != Straight {
		t.Errorf("expected Straight (A-high), got %v", r.Type)
	}
}

func TestEvaluateThreeOfAKind(t *testing.T) {
	cards := []Card{
		card(SuitSpade, Rank7), card(SuitHeart, Rank7), card(SuitDiamond, Rank7),
		card(SuitClub, Rank2), card(SuitSpade, Rank5),
	}
	r := Evaluate(cards)
	if r.Type != ThreeOfAKind {
		t.Errorf("expected ThreeOfAKind, got %v", r.Type)
	}
	if len(r.ScoringCards) != 3 {
		t.Errorf("ThreeOfAKind ScoringCards = %d, want 3", len(r.ScoringCards))
	}
}

func TestEvaluateTwoPair(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankJack), card(SuitHeart, RankJack),
		card(SuitDiamond, Rank9), card(SuitClub, Rank9), card(SuitSpade, Rank3),
	}
	r := Evaluate(cards)
	if r.Type != TwoPair {
		t.Errorf("expected TwoPair, got %v", r.Type)
	}
	if len(r.ScoringCards) != 4 {
		t.Errorf("TwoPair ScoringCards = %d, want 4", len(r.ScoringCards))
	}
}

func TestEvaluateOnePair(t *testing.T) {
	cards := []Card{
		card(SuitSpade, Rank4), card(SuitHeart, Rank4),
		card(SuitDiamond, RankKing), card(SuitClub, Rank7), card(SuitSpade, Rank2),
	}
	r := Evaluate(cards)
	if r.Type != OnePair {
		t.Errorf("expected OnePair, got %v", r.Type)
	}
	if len(r.ScoringCards) != 2 {
		t.Errorf("OnePair ScoringCards = %d, want 2", len(r.ScoringCards))
	}
}

func TestEvaluateHighCard(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankKing),
		card(SuitDiamond, RankQueen), card(SuitClub, RankJack), card(SuitSpade, Rank9),
	}
	r := Evaluate(cards)
	if r.Type != HighCard {
		t.Errorf("expected HighCard, got %v", r.Type)
	}
	if len(r.ScoringCards) != 1 {
		t.Errorf("HighCard ScoringCards = %d, want 1", len(r.ScoringCards))
	}
}

func TestEvaluateTwoCards(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankAce),
	}
	r := Evaluate(cards)
	if r.Type != OnePair {
		t.Errorf("2 cards AA: expected OnePair, got %v", r.Type)
	}
}

func TestEvaluateOneCard(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce),
	}
	r := Evaluate(cards)
	if r.Type != HighCard {
		t.Errorf("1 card A: expected HighCard, got %v", r.Type)
	}
}

func TestEvaluateScoringVsKicker(t *testing.T) {
	// 三条 7，踢牌是 2 和 5
	cards := []Card{
		card(SuitSpade, Rank7), card(SuitHeart, Rank7), card(SuitDiamond, Rank7),
		card(SuitClub, Rank2), card(SuitSpade, Rank5),
	}
	r := Evaluate(cards)
	if len(r.KickerCards) != 2 {
		t.Errorf("ThreeOfAKind KickerCards = %d, want 2", len(r.KickerCards))
	}
}

func TestEvaluateBaseValues(t *testing.T) {
	if BaseChips(StraightFlush) != 100 {
		t.Errorf("StraightFlush BaseChips = %d, want 100", BaseChips(StraightFlush))
	}
	if BaseMult(StraightFlush) != 8 {
		t.Errorf("StraightFlush BaseMult = %d, want 8", BaseMult(StraightFlush))
	}
	if BaseChips(HighCard) != 5 {
		t.Errorf("HighCard BaseChips = %d, want 5", BaseChips(HighCard))
	}
	if BaseMult(HighCard) != 1 {
		t.Errorf("HighCard BaseMult = %d, want 1", BaseMult(HighCard))
	}
}
