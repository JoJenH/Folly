package views

import (
	"fmt"
	"strings"

	"balatro-cli/game"
)

// RenderGameOver 渲染游戏结束界面，返回纯字符串
func RenderGameOver(state game.GameState, victory bool) string {
	var sb strings.Builder

	if victory {
		sb.WriteString(scoreStyle.Render("★ 胜利！Victory！★") + "\n")
		sb.WriteString("恭喜通关 Ante 8！\n")
	} else {
		sb.WriteString(warningStyle.Render("✗ 失败 Game Over") + "\n")
		sb.WriteString("挑战结束。\n")
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Ante: %d\n", state.Run.Ante))
	sb.WriteString(fmt.Sprintf("种子 Seed: %d\n", state.Seed))
	sb.WriteString("\n按任意键退出\n")

	return sb.String()
}
