package engine

import (
	"testing"
)

func TestDeckSize(t *testing.T) {
	d := NewDeck(42)
	if len(d.cards) != 52 {
		t.Errorf("NewDeck should have 52 cards, got %d", len(d.cards))
	}
	// 验证不重复
	seen := make(map[string]bool)
	for _, c := range d.cards {
		key := c.String()
		if seen[key] {
			t.Errorf("Duplicate card: %s", key)
		}
		seen[key] = true
	}
}

func TestDeckShuffleSameSeed(t *testing.T) {
	d1 := NewDeck(12345)
	d1.Shuffle()
	d2 := NewDeck(12345)
	d2.Shuffle()
	for i := range d1.cards {
		if d1.cards[i].String() != d2.cards[i].String() {
			t.Errorf("Same seed should produce same shuffle at index %d: %s vs %s",
				i, d1.cards[i].String(), d2.cards[i].String())
		}
	}
}

func TestDeckShuffleDifferentSeed(t *testing.T) {
	d1 := NewDeck(1)
	d1.Shuffle()
	d2 := NewDeck(2)
	d2.Shuffle()
	diff := false
	for i := range d1.cards {
		if d1.cards[i].String() != d2.cards[i].String() {
			diff = true
			break
		}
	}
	if !diff {
		t.Error("Different seeds should produce different shuffles")
	}
}

func TestDeckDeal(t *testing.T) {
	d := NewDeck(42)
	d.Shuffle()
	dealt := d.Deal(5)
	if len(dealt) != 5 {
		t.Errorf("Deal(5) should return 5 cards, got %d", len(dealt))
	}
	if len(d.cards) != 47 {
		t.Errorf("After Deal(5), deck should have 47 cards, got %d", len(d.cards))
	}
}

func TestDeckDealEmpty(t *testing.T) {
	d := NewDeck(42)
	// 发完所有牌
	d.Deal(52)
	// 空牌堆再发不 panic
	dealt := d.Deal(5)
	if len(dealt) != 0 {
		t.Errorf("Deal from empty deck should return empty slice, got %d cards", len(dealt))
	}
}

func TestDeckDealPartial(t *testing.T) {
	d := NewDeck(42)
	d.Deal(50)
	// 只剩 2 张，请求发 5 张
	dealt := d.Deal(5)
	if len(dealt) != 2 {
		t.Errorf("Deal(5) from 2-card deck should return 2 cards, got %d", len(dealt))
	}
}
