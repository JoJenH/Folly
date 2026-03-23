package views

import (
	"fmt"
	"strings"

	"balatro-cli/engine"
	"balatro-cli/game"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	redSuitStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444"))
	goldCardStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	jokerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC44FF"))
	scoreStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#44FF88"))
	warningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444"))
)

func renderCard(c engine.Card) string {
	s := c.String()
	if c.IsGold() {
		return goldCardStyle.Render(s)
	}
	if c.IsRed() {
		return redSuitStyle.Render(s)
	}
	return s
}

// RenderHand 渲染手牌界面，返回纯字符串。
// cursor 是未选中牌列表中的光标位置（0-based）。
func RenderHand(state game.GameState, cursor int, selected []int) string {
	r := state.Run
	round := r.Round

	selSet := make(map[int]bool, len(selected))
	for _, i := range selected {
		selSet[i] = true
	}

	// 拆分已选中与未选中
	var selCards []engine.Card
	var unselCards []engine.Card
	for i, c := range round.Hand {
		if selSet[i] {
			selCards = append(selCards, c)
		} else {
			unselCards = append(unselCards, c)
		}
	}

	var sb strings.Builder

	// 状态栏
	sb.WriteString(fmt.Sprintf("Ante %d  盲注 %d  $%d  出牌 %d  弃牌 %d  目标 %d  得分 %d\n",
		r.Ante, r.BlindIndex+1, r.Gold,
		round.HandsLeft, round.DiscardsLeft,
		round.Target, round.Score))

	// 小丑栏
	if len(r.Jokers) > 0 {
		sb.WriteString("小丑: ")
		for i, oj := range r.Jokers {
			if i > 0 {
				sb.WriteString(" | ")
			}
			name := oj.DefID
			if def, ok := engine.DefaultRegistry.ByID(oj.DefID); ok {
				name = def.Name
			}
			sb.WriteString(jokerStyle.Render(name))
		}
		sb.WriteString("\n")
	}

	// 已选中区
	sb.WriteString(fmt.Sprintf("\n已选中 (%d/5):\n", len(selCards)))
	if len(selCards) == 0 {
		sb.WriteString(" （无）\n")
	} else {
		sb.WriteString(" ")
		for i, c := range selCards {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(selectedStyle.Render("[") + renderCard(c) + selectedStyle.Render("]"))
		}
		sb.WriteString("\n")
	}

	// 手牌区（未选中）
	sb.WriteString(fmt.Sprintf("\n手牌 (%d张):\n", len(unselCards)))
	if len(unselCards) == 0 {
		sb.WriteString(" （空）\n")
	} else {
		sb.WriteString(" ")
		for i, c := range unselCards {
			if i > 0 {
				sb.WriteString(" ")
			}
			if i == cursor {
				sb.WriteString(cursorStyle.Render(">" + renderCard(c)))
			} else {
				sb.WriteString(renderCard(c))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n[← →] 移动  [Space] 选择  [Backspace] 取消  [Enter] 出牌  [d] 弃牌  [?] 帮助\n")
	return sb.String()
}
