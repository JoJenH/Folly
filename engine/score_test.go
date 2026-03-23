package engine

import "testing"

func TestScoreOnePairAces(t *testing.T) {
	// 一对 A：基础 Chips=10，基础 Mult=2，两张 A 各贡献 11 chips
	// Total = (10 + 11 + 11) * 2 = 64
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankAce),
		card(SuitDiamond, RankKing), card(SuitClub, Rank7), card(SuitSpade, Rank2),
	}
	result := ScoreHand(cards, nil)
	if result.Total != 64 {
		t.Errorf("OnePair AA Total = %d, want 64", result.Total)
	}
	if result.FinalChips != 32 {
		t.Errorf("OnePair AA FinalChips = %d, want 32", result.FinalChips)
	}
	if result.FinalMult != 2 {
		t.Errorf("OnePair AA FinalMult = %d, want 2", result.FinalMult)
	}
}

func TestScoreStraightFlushAllScore(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitSpade, RankKing),
		card(SuitSpade, RankQueen), card(SuitSpade, RankJack), card(SuitSpade, Rank10),
	}
	result := ScoreHand(cards, nil)
	// 所有 5 张都是 ScoringCards，Steps 应有 5+1（5张牌 + 1个手牌类型）
	if len(result.Steps) != 6 {
		t.Errorf("StraightFlush Steps = %d, want 6 (5 cards + 1 hand type)", len(result.Steps))
	}
}

func TestScoreStepsCount(t *testing.T) {
	// 高牌：1 张 ScoringCard，Steps = 1 张牌 + 1 手牌类型 = 2
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankKing),
		card(SuitDiamond, RankQueen), card(SuitClub, RankJack), card(SuitSpade, Rank9),
	}
	result := ScoreHand(cards, nil)
	if len(result.Steps) != 2 {
		t.Errorf("HighCard Steps = %d, want 2 (1 card + 1 hand type)", len(result.Steps))
	}
}

func TestScoreTotalEqualsChipsTimesMult(t *testing.T) {
	cards := []Card{
		card(SuitSpade, Rank7), card(SuitHeart, Rank7), card(SuitDiamond, Rank7),
		card(SuitClub, Rank2), card(SuitSpade, Rank5),
	}
	result := ScoreHand(cards, nil)
	if result.Total != result.FinalChips*result.FinalMult {
		t.Errorf("Total(%d) != FinalChips(%d) * FinalMult(%d)",
			result.Total, result.FinalChips, result.FinalMult)
	}
}
