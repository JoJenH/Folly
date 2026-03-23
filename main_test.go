package main

import (
	"balatro-cli/engine"
	"balatro-cli/game"
	"os"
	"testing"
)

func TestSeedReproducibility(t *testing.T) {
	// 相同 seed 两次 NewGame 手牌相同
	s1 := game.NewGame(12345)
	s2 := game.NewGame(12345)
	if len(s1.Run.Round.Hand) != len(s2.Run.Round.Hand) {
		t.Fatal("hand size mismatch")
	}
	for i := range s1.Run.Round.Hand {
		if s1.Run.Round.Hand[i].String() != s2.Run.Round.Hand[i].String() {
			t.Errorf("hand[%d]: %s vs %s", i,
				s1.Run.Round.Hand[i].String(), s2.Run.Round.Hand[i].String())
		}
	}
}

func TestPlaySequenceReproducible(t *testing.T) {
	// 相同 seed 打出相同的牌，得分相同
	play := func(seed int64) int {
		state := game.NewGame(seed)
		result, err := game.PlayHand(&state, []int{0, 1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}
		return result.Total
	}
	if play(12345) != play(12345) {
		t.Error("same seed should produce same score")
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	state := game.NewGame(99)
	// 打一手牌改变状态
	game.PlayHand(&state, []int{0, 1, 2, 3, 4})

	tmpFile := t.TempDir() + "/save.json"
	if err := engine.SaveGame(state, tmpFile); err != nil {
		t.Fatalf("SaveGame: %v", err)
	}

	loaded, err := engine.LoadGame(tmpFile)
	if err != nil {
		t.Fatalf("LoadGame: %v", err)
	}

	if loaded.Seed != state.Seed {
		t.Errorf("Seed: %d vs %d", loaded.Seed, state.Seed)
	}
	if loaded.Run.Round.HandsLeft != state.Run.Round.HandsLeft {
		t.Errorf("HandsLeft: %d vs %d", loaded.Run.Round.HandsLeft, state.Run.Round.HandsLeft)
	}
	if loaded.Run.Round.Score != state.Run.Round.Score {
		t.Errorf("Score: %d vs %d", loaded.Run.Round.Score, state.Run.Round.Score)
	}
}

func TestSaveLoadCorrupt(t *testing.T) {
	tmpFile := t.TempDir() + "/bad.json"
	os.WriteFile(tmpFile, []byte("corrupted"), 0644)
	_, err := engine.LoadGame(tmpFile)
	if err == nil {
		t.Error("LoadGame corrupt should return error")
	}
}
