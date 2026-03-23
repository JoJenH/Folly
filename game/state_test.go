package game

import "testing"

func TestNewGameInitialState(t *testing.T) {
	state := NewGame(42)

	if state.Run.Ante != 1 {
		t.Errorf("Ante = %d, want 1", state.Run.Ante)
	}
	if state.Run.BlindIndex != 0 {
		t.Errorf("BlindIndex = %d, want 0", state.Run.BlindIndex)
	}
	if state.Run.Gold != 4 {
		t.Errorf("Gold = %d, want 4", state.Run.Gold)
	}
	if state.Run.Round.HandsLeft != 4 {
		t.Errorf("HandsLeft = %d, want 4", state.Run.Round.HandsLeft)
	}
	if state.Run.Round.DiscardsLeft != 4 {
		t.Errorf("DiscardsLeft = %d, want 4", state.Run.Round.DiscardsLeft)
	}
	if state.Run.Round.HandSize != 8 {
		t.Errorf("HandSize = %d, want 8", state.Run.Round.HandSize)
	}
	if len(state.Run.Round.Hand) != 8 {
		t.Errorf("Hand size = %d, want 8", len(state.Run.Round.Hand))
	}
	if len(state.Run.Round.Deck) != 44 {
		t.Errorf("Deck size = %d, want 44", len(state.Run.Round.Deck))
	}
	if state.Run.JokerSlots != 5 {
		t.Errorf("JokerSlots = %d, want 5", state.Run.JokerSlots)
	}
	if len(state.Run.Jokers) != 0 {
		t.Errorf("Jokers = %d, want 0", len(state.Run.Jokers))
	}
}

func TestNewGameSameSeed(t *testing.T) {
	s1 := NewGame(12345)
	s2 := NewGame(12345)
	if len(s1.Run.Round.Hand) != len(s2.Run.Round.Hand) {
		t.Fatal("hand size mismatch")
	}
	for i := range s1.Run.Round.Hand {
		c1 := s1.Run.Round.Hand[i]
		c2 := s2.Run.Round.Hand[i]
		if c1.String() != c2.String() {
			t.Errorf("hand[%d]: %s vs %s", i, c1.String(), c2.String())
		}
	}
}

func TestIsGameOverFalse(t *testing.T) {
	state := NewGame(42)
	if IsGameOver(state) {
		t.Error("fresh game should not be game over")
	}
}

func TestIsVictoryFalse(t *testing.T) {
	state := NewGame(42)
	if IsVictory(state) {
		t.Error("fresh game should not be victory")
	}
}

func TestCurrentTargetAnte1Small(t *testing.T) {
	state := NewGame(42)
	// Ante 1, BlindIndex 0 (Small Blind) = 300
	target := CurrentTarget(state)
	if target != 300 {
		t.Errorf("Ante1 Small Blind target = %d, want 300", target)
	}
}
