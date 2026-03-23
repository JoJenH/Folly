package game

import (
	"balatro-cli/engine"
	"testing"
)

func TestPlayHandDecreasesHandsLeft(t *testing.T) {
	state := NewGame(42)
	// 选前 5 张
	_, err := PlayHand(&state, []int{0, 1, 2, 3, 4})
	if err != nil {
		t.Fatalf("PlayHand error: %v", err)
	}
	if state.Run.Round.HandsLeft != 3 {
		t.Errorf("HandsLeft = %d, want 3", state.Run.Round.HandsLeft)
	}
}

func TestPlayHandReplenishesHand(t *testing.T) {
	state := NewGame(42)
	_, err := PlayHand(&state, []int{0, 1, 2, 3, 4})
	if err != nil {
		t.Fatalf("PlayHand error: %v", err)
	}
	if len(state.Run.Round.Hand) != 8 {
		t.Errorf("Hand size = %d, want 8", len(state.Run.Round.Hand))
	}
}

func TestPlayHandIncreasesScore(t *testing.T) {
	state := NewGame(42)
	result, err := PlayHand(&state, []int{0, 1, 2, 3, 4})
	if err != nil {
		t.Fatalf("PlayHand error: %v", err)
	}
	if state.Run.Round.Score != result.Total {
		t.Errorf("Score = %d, want %d", state.Run.Round.Score, result.Total)
	}
}

func TestIsBlindComplete(t *testing.T) {
	state := NewGame(42)
	state.Run.Round.Score = state.Run.Round.Target
	if !IsBlindComplete(state) {
		t.Error("score >= target should be blind complete")
	}
}

func TestIsBlindCompleteNotYet(t *testing.T) {
	state := NewGame(42)
	if IsBlindComplete(state) {
		t.Error("fresh game should not be blind complete")
	}
}

func TestGameOverWhenHandsExhausted(t *testing.T) {
	state := NewGame(42)
	state.Run.Round.HandsLeft = 1
	// Score < Target
	_, err := PlayHand(&state, []int{0, 1, 2, 3, 4})
	if err != nil {
		t.Fatalf("PlayHand error: %v", err)
	}
	if state.Run.Round.HandsLeft != 0 {
		t.Fatalf("HandsLeft should be 0")
	}
	// Score unlikely to reach 300 with 5 random cards from seed 42,
	// but to be deterministic, set score manually
	state.Run.Round.Score = 0
	if !IsGameOver(state) {
		t.Error("hands=0, score<target should be game over")
	}
}

func TestPlayHandTooManyCards(t *testing.T) {
	state := NewGame(42)
	_, err := PlayHand(&state, []int{0, 1, 2, 3, 4, 5})
	if err == nil {
		t.Error("selecting 6 cards should return error")
	}
}

func TestPlayHandZeroCards(t *testing.T) {
	state := NewGame(42)
	_, err := PlayHand(&state, []int{})
	if err == nil {
		t.Error("selecting 0 cards should return error")
	}
}

func TestDiscardDecreasesDiscardsLeft(t *testing.T) {
	state := NewGame(42)
	err := Discard(&state, []int{0, 1, 2})
	if err != nil {
		t.Fatalf("Discard error: %v", err)
	}
	if state.Run.Round.DiscardsLeft != 3 {
		t.Errorf("DiscardsLeft = %d, want 3", state.Run.Round.DiscardsLeft)
	}
	if len(state.Run.Round.Hand) != 8 {
		t.Errorf("Hand size = %d, want 8", len(state.Run.Round.Hand))
	}
}

func TestDiscardExhausted(t *testing.T) {
	state := NewGame(42)
	state.Run.Round.DiscardsLeft = 0
	err := Discard(&state, []int{0})
	if err == nil {
		t.Error("discard with DiscardsLeft=0 should return error")
	}
}

func TestDiscardZeroCards(t *testing.T) {
	state := NewGame(42)
	err := Discard(&state, []int{})
	if err == nil {
		t.Error("discard 0 cards should return error")
	}
}

func TestGoldCardAwardsGold(t *testing.T) {
	state := NewGame(42)
	// 将手牌第 0 张设为 Gold Card
	state.Run.Round.Hand[0].Enhancement = engine.EnhancementGold
	initialGold := state.Run.Gold
	_, err := PlayHand(&state, []int{0})
	if err != nil {
		t.Fatalf("PlayHand error: %v", err)
	}
	if state.Run.Gold != initialGold+1 {
		t.Errorf("Gold = %d, want %d after playing Gold Card", state.Run.Gold, initialGold+1)
	}
}
