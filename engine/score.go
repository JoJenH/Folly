package engine

import "fmt"

// GameEvent 游戏事件类型
type GameEvent int

const (
	OnCardScored GameEvent = iota
	OnHandScored
	OnDiscard
	OnRoundEnd
)

// HookResult Hook 返回值
type HookResult struct {
	ChipsDelta int
	MultDelta  int
	Retrigger  int
}

// JokerHook Joker Hook 接口，供 ScoreHand 调用
type JokerHook interface {
	ApplyHook(event GameEvent, ctx *ScoreContext) HookResult
}

// ScoreContext 计分上下文，在 Hook 链中传递
type ScoreContext struct {
	Chips        int
	Mult         int
	PlayedCards  []Card
	ScoringCards []Card
	CurrentCard  *Card // 当前正在计分的牌（OnCardScored 时有值）
}

// ScoreStep 单步计分记录
type ScoreStep struct {
	Description string
	ChipsAfter  int
	MultAfter   int
}

// ScoreResult 完整计分结果
type ScoreResult struct {
	HandType   HandType
	FinalChips int
	FinalMult  int
	Total      int
	Steps      []ScoreStep
}

// ScoreHand 完整计分流水线
func ScoreHand(cards []Card, jokers []JokerHook) ScoreResult {
	eval := Evaluate(cards)

	ctx := &ScoreContext{
		Chips:        eval.BaseChips,
		Mult:         eval.BaseMult,
		PlayedCards:  cards,
		ScoringCards: eval.ScoringCards,
	}

	var steps []ScoreStep

	// 步骤 0：手牌类型
	steps = append(steps, ScoreStep{
		Description: fmt.Sprintf("%s（基础 Chips+%d, Mult×%d）", eval.Type, eval.BaseChips, eval.BaseMult),
		ChipsAfter:  ctx.Chips,
		MultAfter:   ctx.Mult,
	})

	// 遍历 ScoringCards，累加 Chips
	for i := range eval.ScoringCards {
		c := eval.ScoringCards[i]
		ctx.CurrentCard = &c
		chipDelta := c.Chips()
		ctx.Chips += chipDelta

	// OnCardScored Hook 预留位
		for _, j := range jokers {
			hr := j.ApplyHook(OnCardScored, ctx)
			ctx.Chips += hr.ChipsDelta
			ctx.Mult += hr.MultDelta
			for r := 0; r < hr.Retrigger; r++ {
				ctx.Chips += chipDelta
				steps = append(steps, ScoreStep{
					Description: fmt.Sprintf("%s 重复触发 +%d chips", c.String(), chipDelta),
					ChipsAfter:  ctx.Chips,
					MultAfter:   ctx.Mult,
				})
			}
		}

		steps = append(steps, ScoreStep{
			Description: fmt.Sprintf("%s +%d chips", c.String(), chipDelta),
			ChipsAfter:  ctx.Chips,
			MultAfter:   ctx.Mult,
		})
	}

	// OnHandScored Hook 预留位
	ctx.CurrentCard = nil
	for _, j := range jokers {
		hr := j.ApplyHook(OnHandScored, ctx)
		ctx.Chips += hr.ChipsDelta
		ctx.Mult += hr.MultDelta
		if hr.ChipsDelta != 0 || hr.MultDelta != 0 {
			steps = append(steps, ScoreStep{
				Description: fmt.Sprintf("Joker OnHandScored chips+%d mult+%d", hr.ChipsDelta, hr.MultDelta),
				ChipsAfter:  ctx.Chips,
				MultAfter:   ctx.Mult,
			})
		}
	}

	return ScoreResult{
		HandType:   eval.Type,
		FinalChips: ctx.Chips,
		FinalMult:  ctx.Mult,
		Total:      ctx.Chips * ctx.Mult,
		Steps:      steps,
	}
}
