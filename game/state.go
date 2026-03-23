package game

import (
	"balatro-cli/engine"
)

// 类型别名，直接使用 engine 包中的定义
type OwnedJoker = engine.OwnedJoker
type RoundState = engine.RoundState
type RunState = engine.RunState
type GameState = engine.GameState

// NewGame 创建新游戏
func NewGame(seed int64) GameState {
	deck := engine.NewDeck(seed)
	deck.Shuffle()
	hand := deck.Deal(8)
	sortHand(hand)
	remaining := deck.Deal(44)

	target := blindTarget(1, 0)

	return GameState{
		Seed: seed,
		Run: RunState{
			Ante:       1,
			BlindIndex: 0,
			Gold:       4,
			Jokers:     []OwnedJoker{},
			JokerSlots: 5,
			Round: RoundState{
				Deck:         remaining,
				Hand:         hand,
				HandSize:     8,
				HandsLeft:    4,
				DiscardsLeft: 4,
				Score:        0,
				Target:       target,
			},
		},
	}
}

// IsGameOver 出牌次数耗尽且未达目标分数
func IsGameOver(state GameState) bool {
	r := state.Run.Round
	return r.HandsLeft == 0 && r.Score < r.Target
}

// IsVictory 通过 Ante 8 Boss Blind（Ante 超过 8）
func IsVictory(state GameState) bool {
	return state.Run.Ante > 8
}

// CurrentTarget 返回当前 Blind 的目标分数
func CurrentTarget(state GameState) int {
	return blindTarget(state.Run.Ante, state.Run.BlindIndex)
}

// blindTarget 根据 Ante 和 BlindIndex 返回目标分数
func blindTarget(ante, blindIndex int) int {
	type blindRow struct{ small, big, boss int }
	table := []blindRow{
		{300, 450, 600},
		{800, 1200, 1600},
		{2000, 3000, 4000},
		{5000, 7500, 10000},
		{11000, 16500, 22000},
		{20000, 30000, 40000},
		{35000, 52500, 70000},
		{50000, 75000, 100000},
	}
	if ante < 1 || ante > 8 {
		return 0
	}
	row := table[ante-1]
	switch blindIndex {
	case 0:
		return row.small
	case 1:
		return row.big
	case 2:
		return row.boss
	}
	return 0
}
