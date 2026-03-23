package engine

import (
	"os"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	state := GameState{
		Seed: 12345,
		Run: RunState{
			Ante:       2,
			BlindIndex: 1,
			Gold:       7,
			JokerSlots: 5,
			Jokers:     []OwnedJoker{{DefID: "greedy"}},
			Round: RoundState{
				HandsLeft:    3,
				DiscardsLeft: 2,
				HandSize:     8,
				Score:        500,
				Target:       1200,
			},
		},
	}

	tmpFile := t.TempDir() + "/save.json"

	if err := SaveGame(state, tmpFile); err != nil {
		t.Fatalf("SaveGame error: %v", err)
	}

	// 验证文件存在且为合法 JSON
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if len(data) == 0 {
		t.Error("saved file is empty")
	}

	loaded, err := LoadGame(tmpFile)
	if err != nil {
		t.Fatalf("LoadGame error: %v", err)
	}

	if loaded.Seed != state.Seed {
		t.Errorf("Seed = %d, want %d", loaded.Seed, state.Seed)
	}
	if loaded.Run.Ante != state.Run.Ante {
		t.Errorf("Ante = %d, want %d", loaded.Run.Ante, state.Run.Ante)
	}
	if loaded.Run.Gold != state.Run.Gold {
		t.Errorf("Gold = %d, want %d", loaded.Run.Gold, state.Run.Gold)
	}
	if len(loaded.Run.Jokers) != 1 || loaded.Run.Jokers[0].DefID != "greedy" {
		t.Errorf("Jokers = %v, want [{greedy}]", loaded.Run.Jokers)
	}
	if loaded.Run.Round.Score != state.Run.Round.Score {
		t.Errorf("Score = %d, want %d", loaded.Run.Round.Score, state.Run.Round.Score)
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	tmpFile := t.TempDir() + "/corrupt.json"
	if err := os.WriteFile(tmpFile, []byte("not valid json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadGame(tmpFile)
	if err == nil {
		t.Error("LoadGame corrupt JSON should return error")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := LoadGame("/nonexistent/path/save.json")
	if err == nil {
		t.Error("LoadGame missing file should return error")
	}
}

func TestSavePathNonEmpty(t *testing.T) {
	p := SavePath()
	if p == "" {
		t.Error("SavePath() should return non-empty string")
	}
}
