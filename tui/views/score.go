package views

import (
	"fmt"
	"strings"

	"balatro-cli/engine"
	"github.com/charmbracelet/lipgloss"
)

var (
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	flashOnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#44FF88")).Bold(true)
	currentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#44FF88")).Bold(true)
)

// RenderScore 渲染计分过程（静态，供测试使用）。
// accumulated 为可选参数，若提供则用累计分数判断是否达标，否则用 result.Total。
func RenderScore(result engine.ScoreResult, target int, accumulated ...int) string {
	var sb strings.Builder

	sb.WriteString(scoreStyle.Render(result.HandType.String()) + "\n")
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for _, step := range result.Steps {
		sb.WriteString(fmt.Sprintf("  %s  (Chips:%d × Mult:%d)\n",
			step.Description, step.ChipsAfter, step.MultAfter))
	}

	sb.WriteString(strings.Repeat("-", 40) + "\n")
	sb.WriteString(scoreStyle.Render(
		fmt.Sprintf("总分 Total: %d × %d = %d", result.FinalChips, result.FinalMult, result.Total),
	) + "\n")

	scoreToCheck := result.Total
	if len(accumulated) > 0 {
		scoreToCheck = accumulated[0]
	}
	if scoreToCheck >= target {
		sb.WriteString(scoreStyle.Render(fmt.Sprintf("✓ 达成目标（目标: %d）", target)) + "\n")
	} else {
		sb.WriteString(warningStyle.Render(fmt.Sprintf("✗ 未达目标（目标: %d）", target)) + "\n")
	}

	return sb.String()
}

// RenderScoreAnimated 动画渲染计分，逐步展示每个计分项。
// animStep: 当前高亮的步骤索引（0-based）；>= len(steps) 表示动画结束。
// flash: 当前步骤是否处于高亮闪烁状态。
func RenderScoreAnimated(result engine.ScoreResult, accumulatedScore int, target int, animStep int, flash bool) string {
	var sb strings.Builder

	// 手牌类型标题
	sb.WriteString("\n" + scoreStyle.Render("【"+result.HandType.String()+"】") + "\n")
	sb.WriteString(strings.Repeat("─", 44) + "\n")

	// 计分步骤
	for i, step := range result.Steps {
		if i > animStep {
			break
		}
		line := fmt.Sprintf("  %-38s  %d × %d", step.Description, step.ChipsAfter, step.MultAfter)
		if i == animStep {
			if flash {
				sb.WriteString(flashOnStyle.Render(line) + "\n")
			} else {
				sb.WriteString(currentStyle.Render(line) + "\n")
			}
		} else {
			sb.WriteString(dimStyle.Render(line) + "\n")
		}
	}

	sb.WriteString(strings.Repeat("─", 44) + "\n")

	// 实时计分显示
	var chips, mult int
	if animStep >= 0 && len(result.Steps) > 0 {
		step := result.Steps[min(animStep, len(result.Steps)-1)]
		chips = step.ChipsAfter
		mult = step.MultAfter
	}
	totalNow := chips * mult
	sb.WriteString(scoreStyle.Render(fmt.Sprintf("  Chips %-5d × Mult %-3d = %d", chips, mult, totalNow)) + "\n")

	// 动画结束后显示最终结果
	if animStep >= len(result.Steps) {
		sb.WriteString("\n")
		if accumulatedScore >= target {
			sb.WriteString(scoreStyle.Render(
				fmt.Sprintf("✓ 达成目标！（目标: %d，累计: %d）", target, accumulatedScore),
			) + "\n")
		} else {
			sb.WriteString(warningStyle.Render(
				fmt.Sprintf("✗ 未达目标（目标: %d，累计: %d）", target, accumulatedScore),
			) + "\n")
		}
		sb.WriteString("\n按任意键继续\n")
	} else {
		sb.WriteString("\n（按任意键跳过）\n")
	}

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
