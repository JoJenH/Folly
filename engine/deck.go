package engine

import "math/rand/v2"

// Deck 牌堆
type Deck struct {
	cards []Card
	rng   *rand.Rand
}

// NewDeck 创建一副标准 52 张牌，使用 seed 初始化 rng
func NewDeck(seed int64) Deck {
	var cards []Card
	for _, suit := range []Suit{SuitSpade, SuitHeart, SuitDiamond, SuitClub} {
		for rank := Rank2; rank <= RankAce; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	return Deck{
		cards: cards,
		rng:   rand.New(rand.NewPCG(uint64(seed), 0)),
	}
}

// Shuffle Fisher-Yates 洗牌
func (d *Deck) Shuffle() {
	n := len(d.cards)
	for i := n - 1; i > 0; i-- {
		j := int(d.rng.Int64N(int64(i + 1)))
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

// Deal 从顶部取 n 张牌；牌堆不足时返回剩余全部
func (d *Deck) Deal(n int) []Card {
	if len(d.cards) == 0 {
		return []Card{}
	}
	if n > len(d.cards) {
		n = len(d.cards)
	}
	dealt := make([]Card, n)
	copy(dealt, d.cards[:n])
	d.cards = d.cards[n:]
	return dealt
}
