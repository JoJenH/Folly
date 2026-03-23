package jokers

import "balatro-cli/engine"

func init() {
	engine.DefaultRegistry.Register(&engine.JokerDef{
		ID:          "retrigger",
		Name:        "袜子和半身裙",
		Description: "重复触发最后打出的计分牌一次",
		Cost:        6,
		Hooks: map[engine.GameEvent]func(ctx *engine.ScoreContext) engine.HookResult{
			engine.OnCardScored: func(ctx *engine.ScoreContext) engine.HookResult {
				if ctx.CurrentCard == nil || len(ctx.ScoringCards) == 0 {
					return engine.HookResult{}
				}
				// 检查当前牌是否为最后一张 ScoringCard
				last := ctx.ScoringCards[len(ctx.ScoringCards)-1]
				if ctx.CurrentCard.Suit == last.Suit && ctx.CurrentCard.Rank == last.Rank {
					return engine.HookResult{Retrigger: 1}
				}
				return engine.HookResult{}
			},
		},
	})
}
