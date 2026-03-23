package engine

import "fmt"

// Suit 花色
type Suit int

const (
	SuitSpade   Suit = iota // ♠
	SuitHeart               // ♥
	SuitDiamond             // ♦
	SuitClub                // ♣
)

func (s Suit) String() string {
	switch s {
	case SuitSpade:
		return "♠"
	case SuitHeart:
		return "♥"
	case SuitDiamond:
		return "♦"
	case SuitClub:
		return "♣"
	}
	return "?"
}

// Rank 点数，2–14（14 = A）
type Rank int

const (
	Rank2  Rank = 2
	Rank3  Rank = 3
	Rank4  Rank = 4
	Rank5  Rank = 5
	Rank6  Rank = 6
	Rank7  Rank = 7
	Rank8  Rank = 8
	Rank9  Rank = 9
	Rank10  Rank = 10
	RankJack  Rank = 11
	RankQueen Rank = 12
	RankKing  Rank = 13
	RankAce   Rank = 14
)

func (r Rank) String() string {
	switch r {
	case RankJack:
		return "J"
	case RankQueen:
		return "Q"
	case RankKing:
		return "K"
	case RankAce:
		return "A"
	}
	return fmt.Sprintf("%d", int(r))
}

// Enhancement 牌面增强
type Enhancement int

const (
	EnhancementNone Enhancement = iota
	EnhancementGold
)

// Card 单张牌
type Card struct {
	Suit        Suit
	Rank        Rank
	Enhancement Enhancement
}

// Chips 返回牌面基础筹码值
func (c Card) Chips() int {
	if c.Rank >= Rank10 && c.Rank <= RankKing {
		return 10
	}
	if c.Rank == RankAce {
		return 11
	}
	return int(c.Rank)
}

// String 返回 "[A♠]" 格式；Gold Card 含 "$" 标记
func (c Card) String() string {
	if c.Enhancement == EnhancementGold {
		return fmt.Sprintf("[$%s%s]", c.Rank.String(), c.Suit.String())
	}
	return fmt.Sprintf("[%s%s]", c.Rank.String(), c.Suit.String())
}

// IsRed 红心/方块为红色
func (c Card) IsRed() bool {
	return c.Suit == SuitHeart || c.Suit == SuitDiamond
}

// IsGold 是否为 Gold Card
func (c Card) IsGold() bool {
	return c.Enhancement == EnhancementGold
}
