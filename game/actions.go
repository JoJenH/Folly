package game

import (
	"errors"
	"sort"

	"balatro-cli/engine"
)

// PlayHand 打出选中的牌，返回计分结果
func PlayHand(state *GameState, selectedIdx []int) (engine.ScoreResult, error) {
	if len(selectedIdx) == 0 {
		return engine.ScoreResult{}, errors.New("必须至少选择 1 张牌")
	}
	if len(selectedIdx) > 5 {
		return engine.ScoreResult{}, errors.New("最多选择 5 张牌")
	}
	if state.Run.Round.HandsLeft <= 0 {
		return engine.ScoreResult{}, errors.New("出牌次数已用完")
	}

	// 验证索引合法
	hand := state.Run.Round.Hand
	for _, idx := range selectedIdx {
		if idx < 0 || idx >= len(hand) {
			return engine.ScoreResult{}, errors.New("无效的牌索引")
		}
	}

	// 收集选中的牌
	played := make([]engine.Card, len(selectedIdx))
	for i, idx := range selectedIdx {
		played[i] = hand[idx]
	}

	// Gold Card 奖励金币
	for _, c := range played {
		if c.IsGold() {
			state.Run.Gold++
		}
	}

	// 构建 Joker Hook 列表
	jokers := make([]engine.JokerHook, 0, len(state.Run.Jokers))
	for _, oj := range state.Run.Jokers {
		if def, ok := engine.DefaultRegistry.ByID(oj.DefID); ok {
			jokers = append(jokers, def)
		}
	}

	// 计分
	result := engine.ScoreHand(played, jokers)
	state.Run.Round.Score += result.Total
	state.Run.Round.HandsLeft--

	// 从手牌移除打出的牌（从高到低移除避免索引偏移）
	state.Run.Round.Hand = removeIndices(hand, selectedIdx)

	// 补牌至 HandSize
	need := state.Run.Round.HandSize - len(state.Run.Round.Hand)
	if need > 0 && len(state.Run.Round.Deck) > 0 {
		draw := need
		if draw > len(state.Run.Round.Deck) {
			draw = len(state.Run.Round.Deck)
		}
		state.Run.Round.Hand = append(state.Run.Round.Hand, state.Run.Round.Deck[:draw]...)
		state.Run.Round.Deck = state.Run.Round.Deck[draw:]
	}
	sortHand(state.Run.Round.Hand)

	return result, nil
}

// Discard 弃掉选中的牌并补牌
func Discard(state *GameState, selectedIdx []int) error {
	if len(selectedIdx) == 0 {
		return errors.New("必须至少选择 1 张牌")
	}
	if len(selectedIdx) > 5 {
		return errors.New("最多弃 5 张牌")
	}
	if state.Run.Round.DiscardsLeft <= 0 {
		return errors.New("弃牌次数已用完")
	}

	hand := state.Run.Round.Hand
	for _, idx := range selectedIdx {
		if idx < 0 || idx >= len(hand) {
			return errors.New("无效的牌索引")
		}
	}

	state.Run.Round.DiscardsLeft--
	state.Run.Round.Hand = removeIndices(hand, selectedIdx)

	// 补牌
	need := state.Run.Round.HandSize - len(state.Run.Round.Hand)
	if need > 0 && len(state.Run.Round.Deck) > 0 {
		draw := need
		if draw > len(state.Run.Round.Deck) {
			draw = len(state.Run.Round.Deck)
		}
		state.Run.Round.Hand = append(state.Run.Round.Hand, state.Run.Round.Deck[:draw]...)
		state.Run.Round.Deck = state.Run.Round.Deck[draw:]
	}
	sortHand(state.Run.Round.Hand)

	return nil
}

// IsBlindComplete 累计得分是否达到目标
func IsBlindComplete(state GameState) bool {
	return state.Run.Round.Score >= state.Run.Round.Target
}

// NextBlind 推进到下一个 Blind，发放金币奖励并重置回合状态
func NextBlind(state *GameState) error {
	if !IsBlindComplete(*state) {
		return errors.New("尚未完成当前 Blind")
	}

	// 发放 Blind 奖励
	rewards := []int{3, 4, 5} // Small=3, Big=4, Boss=5
	state.Run.Gold += rewards[state.Run.BlindIndex]

	// 推进 Blind/Ante
	state.Run.BlindIndex++
	if state.Run.BlindIndex > 2 {
		state.Run.BlindIndex = 0
		state.Run.Ante++
	}

	if IsVictory(*state) {
		return nil
	}

	// 重置回合状态
	newTarget := blindTarget(state.Run.Ante, state.Run.BlindIndex)
	state.Run.Round.Score = 0
	state.Run.Round.Target = newTarget
	state.Run.Round.HandsLeft = 4
	state.Run.Round.DiscardsLeft = 4

	// 刷新牌堆：用派生 seed 生成新牌堆并洗牌
	newDeck := engine.NewDeck(state.Seed*100 + int64(state.Run.Ante)*10 + int64(state.Run.BlindIndex))
	newDeck.Shuffle()
	state.Run.Round.Hand = newDeck.Deal(state.Run.Round.HandSize)
	state.Run.Round.Deck = newDeck.Deal(52 - state.Run.Round.HandSize)
	sortHand(state.Run.Round.Hand)

	return nil
}

// BuyJoker 购买 Joker
func BuyJoker(state *GameState, jokerID string) error {
	if len(state.Run.Jokers) >= state.Run.JokerSlots {
		return errors.New("Joker 槽位已满")
	}
	j, ok := engine.DefaultRegistry.ByID(jokerID)
	if !ok {
		return errors.New("未知的 Joker ID")
	}
	if state.Run.Gold < j.Cost {
		return errors.New("金币不足")
	}
	state.Run.Gold -= j.Cost
	state.Run.Jokers = append(state.Run.Jokers, OwnedJoker{DefID: jokerID})
	return nil
}

// removeIndices 从切片中移除指定索引的元素（原地操作，不修改原切片）
func removeIndices(cards []engine.Card, indices []int) []engine.Card {
	remove := make(map[int]bool, len(indices))
	for _, i := range indices {
		remove[i] = true
	}
	result := make([]engine.Card, 0, len(cards)-len(indices))
	for i, c := range cards {
		if !remove[i] {
			result = append(result, c)
		}
	}
	return result
}

// sortHand 按点数降序、花色升序排列手牌
func sortHand(hand []engine.Card) {
	sort.SliceStable(hand, func(i, j int) bool {
		if hand[i].Rank != hand[j].Rank {
			return hand[i].Rank > hand[j].Rank
		}
		return hand[i].Suit < hand[j].Suit
	})
}
