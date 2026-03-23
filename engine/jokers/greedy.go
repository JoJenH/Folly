package jokers

import "balatro-cli/engine"

func init() {
	engine.DefaultRegistry.Register(&engine.JokerDef{
		ID:          "greedy",
		Name:        "贪心小丑",
		Description: "每张打出的♦牌计分时 +3 Mult",
		Cost:        5,
		Hooks: map[engine.GameEvent]func(ctx *engine.ScoreContext) engine.HookResult{
			engine.OnCardScored: func(ctx *engine.ScoreContext) engine.HookResult {
				if ctx.CurrentCard != nil && ctx.CurrentCard.Suit == engine.SuitDiamond {
					return engine.HookResult{MultDelta: 3}
				}
				return engine.HookResult{}
			},
		},
	})
}
