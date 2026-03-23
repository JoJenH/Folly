package jokers

import "balatro-cli/engine"

func init() {
	engine.DefaultRegistry.Register(&engine.JokerDef{
		ID:          "half",
		Name:        "半张",
		Description: "手牌张数 ≤3 时 +20 Mult",
		Cost:        5,
		Hooks: map[engine.GameEvent]func(ctx *engine.ScoreContext) engine.HookResult{
			engine.OnHandScored: func(ctx *engine.ScoreContext) engine.HookResult {
				if len(ctx.PlayedCards) <= 3 {
					return engine.HookResult{MultDelta: 20}
				}
				return engine.HookResult{}
			},
		},
	})
}
