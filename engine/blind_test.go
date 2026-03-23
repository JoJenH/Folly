package engine

import "testing"

func TestBlindTargetAnte1Small(t *testing.T) {
	b := BlindAt(1, 0)
	if b.Target != 300 {
		t.Errorf("Ante1 Small target = %d, want 300", b.Target)
	}
}

func TestBlindTargetAnte1Big(t *testing.T) {
	b := BlindAt(1, 1)
	if b.Target != 450 {
		t.Errorf("Ante1 Big target = %d, want 450", b.Target)
	}
}

func TestBlindTargetAnte1Boss(t *testing.T) {
	b := BlindAt(1, 2)
	if b.Target != 600 {
		t.Errorf("Ante1 Boss target = %d, want 600", b.Target)
	}
}

func TestBlindTargetAnte8Boss(t *testing.T) {
	b := BlindAt(8, 2)
	if b.Target != 100000 {
		t.Errorf("Ante8 Boss target = %d, want 100000", b.Target)
	}
}

func TestBlindReward(t *testing.T) {
	small := BlindAt(1, 0)
	if small.Reward != 3 {
		t.Errorf("Small Blind reward = %d, want 3", small.Reward)
	}
	big := BlindAt(1, 1)
	if big.Reward != 4 {
		t.Errorf("Big Blind reward = %d, want 4", big.Reward)
	}
	boss := BlindAt(1, 2)
	if boss.Reward != 5 {
		t.Errorf("Boss Blind reward = %d, want 5", boss.Reward)
	}
}

func TestBlindNames(t *testing.T) {
	small := BlindAt(1, 0)
	if small.Name == "" {
		t.Error("Small Blind name should not be empty")
	}
	big := BlindAt(1, 1)
	if big.Name == "" {
		t.Error("Big Blind name should not be empty")
	}
	boss := BlindAt(1, 2)
	if boss.Name == "" {
		t.Error("Boss Blind name should not be empty")
	}
}
