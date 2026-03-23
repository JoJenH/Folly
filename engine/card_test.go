package engine

import (
	"strings"
	"testing"
)

func TestCardString(t *testing.T) {
	tests := []struct {
		card Card
		want string
	}{
		{Card{Suit: SuitSpade, Rank: RankAce}, "[A♠]"},
		{Card{Suit: SuitHeart, Rank: Rank10}, "[10♥]"},
		{Card{Suit: SuitDiamond, Rank: Rank2}, "[2♦]"},
		{Card{Suit: SuitClub, Rank: RankKing}, "[K♣]"},
		{Card{Suit: SuitSpade, Rank: RankJack}, "[J♠]"},
		{Card{Suit: SuitHeart, Rank: RankQueen}, "[Q♥]"},
	}
	for _, tt := range tests {
		got := tt.card.String()
		if got != tt.want {
			t.Errorf("Card{%v,%v}.String() = %q, want %q", tt.card.Suit, tt.card.Rank, got, tt.want)
		}
	}
}

func TestCardIsRed(t *testing.T) {
	heart := Card{Suit: SuitHeart, Rank: RankAce}
	diamond := Card{Suit: SuitDiamond, Rank: RankAce}
	spade := Card{Suit: SuitSpade, Rank: RankAce}
	club := Card{Suit: SuitClub, Rank: RankAce}
	if !heart.IsRed() {
		t.Error("Heart should be red")
	}
	if !diamond.IsRed() {
		t.Error("Diamond should be red")
	}
	if spade.IsRed() {
		t.Error("Spade should not be red")
	}
	if club.IsRed() {
		t.Error("Club should not be red")
	}
}

func TestGoldCardString(t *testing.T) {
	c := Card{Suit: SuitSpade, Rank: RankAce, Enhancement: EnhancementGold}
	s := c.String()
	if !strings.Contains(s, "$") {
		t.Errorf("Gold card String() = %q, want '$' marker", s)
	}
	if !c.IsGold() {
		t.Error("Card with EnhancementGold should return IsGold() = true")
	}
	normal := Card{Suit: SuitSpade, Rank: RankAce}
	if normal.IsGold() {
		t.Error("Normal card should return IsGold() = false")
	}
}

func TestCardChips(t *testing.T) {
	tests := []struct {
		rank Rank
		want int
	}{
		{Rank2, 2}, {Rank3, 3}, {Rank4, 4}, {Rank5, 5},
		{Rank6, 6}, {Rank7, 7}, {Rank8, 8}, {Rank9, 9},
		{Rank10, 10}, {RankJack, 10}, {RankQueen, 10}, {RankKing, 10},
		{RankAce, 11},
	}
	for _, tt := range tests {
		c := Card{Suit: SuitSpade, Rank: tt.rank}
		if got := c.Chips(); got != tt.want {
			t.Errorf("Card{Rank:%v}.Chips() = %d, want %d", tt.rank, got, tt.want)
		}
	}
}
