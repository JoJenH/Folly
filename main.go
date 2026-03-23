package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"balatro-cli/engine"
	"balatro-cli/game"
	"balatro-cli/tui"
)

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", 0, "游戏随机种子（0=随机生成）")
	flag.Parse()

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	var state game.GameState
	savePath := engine.SavePath()

	// 检测存档
	if existing, err := engine.LoadGame(savePath); err == nil {
		fmt.Printf("发现存档（Ante %d，种子 %d）。继续游戏？[y/n]: ",
			existing.Run.Ante, existing.Seed)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(strings.TrimSpace(answer)) == "y" {
			state = existing
		} else {
			os.Remove(savePath)
			state = game.NewGame(seed)
		}
	} else {
		state = game.NewGame(seed)
	}

	if err := tui.Start(state); err != nil {
		fmt.Fprintf(os.Stderr, "TUI 错误: %v\n", err)
		os.Exit(1)
	}
}
