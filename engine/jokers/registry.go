// Package jokers 包含所有内置 Joker 实现，每个 Joker 通过 init() 自注册到 DefaultRegistry。
package jokers

import "balatro-cli/engine"

func init() {
	// 各 Joker 文件（greedy.go、half.go、retrigger.go 等）各自通过 init() 自注册。
	// 本文件仅作为包入口，确保 import 本包时所有 init() 得到执行。
	//
	// 为使注册表非空，此处注册一个最基础的 no-op 占位 Joker（测试锚点）。
	engine.DefaultRegistry.Register(&engine.JokerDef{
		ID:          "noop",
		Name:        "占位",
		Description: "无效果，仅用于确保注册表非空",
		Cost:        0,
	})
}
