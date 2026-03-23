package engine

import "testing"

// mockJoker 实现 JokerHook 接口用于测试
type mockJoker struct {
	event     GameEvent
	chipsDelta int
	multDelta  int
	retrigger  int
	called     int
}

func (m *mockJoker) ApplyHook(event GameEvent, ctx *ScoreContext) HookResult {
	if event == m.event {
		m.called++
		return HookResult{ChipsDelta: m.chipsDelta, MultDelta: m.multDelta, Retrigger: m.retrigger}
	}
	return HookResult{}
}

func TestScoreNoJokers(t *testing.T) {
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankAce),
	}
	without := ScoreHand(cards, nil)
	withEmpty := ScoreHand(cards, []JokerHook{})
	if without.Total != withEmpty.Total {
		t.Errorf("empty joker list should match nil: %d vs %d", without.Total, withEmpty.Total)
	}
}

func TestScoreMockJokerPlusMult(t *testing.T) {
	// 一对 A 无 Joker: (10+11+11)*2 = 64
	// 加 +10 Mult OnHandScored Joker: (10+11+11)*(2+10) = 384
	cards := []Card{
		card(SuitSpade, RankAce), card(SuitHeart, RankAce),
		card(SuitDiamond, RankKing), card(SuitClub, Rank7), card(SuitSpade, Rank2),
	}
	mj := &mockJoker{event: OnHandScored, multDelta: 10}
	result := ScoreHand(cards, []JokerHook{mj})
	if result.FinalMult != 12 {
		t.Errorf("FinalMult = %d, want 12", result.FinalMult)
	}
	if result.Total != 32*12 {
		t.Errorf("Total = %d, want %d", result.Total, 32*12)
	}
}

func TestScoreJokerOrder(t *testing.T) {
	// 两个 Joker 按顺序触发，都是 OnHandScored +5 Mult
	cards := []Card{card(SuitSpade, RankAce)}
	mj1 := &mockJoker{event: OnHandScored, multDelta: 5}
	mj2 := &mockJoker{event: OnHandScored, multDelta: 5}
	result := ScoreHand(cards, []JokerHook{mj1, mj2})
	if mj1.called != 1 {
		t.Errorf("mj1 called %d times, want 1", mj1.called)
	}
	if mj2.called != 1 {
		t.Errorf("mj2 called %d times, want 1", mj2.called)
	}
	// HighCard A: (5+11)*(1+5+5) = 16*11 = 176
	if result.Total != 16*11 {
		t.Errorf("Total = %d, want %d", result.Total, 16*11)
	}
}

func TestScoreRetrigger(t *testing.T) {
	// 单张 A，Retrigger=1 时该牌计分触发两次，Chips = 5 + 11 + 11 = 27
	cards := []Card{card(SuitSpade, RankAce)}
	mj := &mockJoker{event: OnCardScored, retrigger: 1}
	result := ScoreHand(cards, []JokerHook{mj})
	if result.FinalChips != 5+11+11 {
		t.Errorf("Retrigger FinalChips = %d, want %d", result.FinalChips, 5+11+11)
	}
	// Steps: 1 hand type + 1 retrigger step + 1 normal card step = 3
	if len(result.Steps) != 3 {
		t.Errorf("Retrigger Steps = %d, want 3", len(result.Steps))
	}
}
