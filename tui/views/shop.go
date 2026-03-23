package views

import (
	"fmt"
	"strings"

	"balatro-cli/engine"
)

// RenderShop 渲染商店界面，返回纯字符串
func RenderShop(items []engine.ShopItem, gold int, cursor int) string {
	var sb strings.Builder

	sb.WriteString("=== 商店 ===\n")
	sb.WriteString(fmt.Sprintf("金币: $%d\n\n", gold))

	for i, item := range items {
		prefix := "  "
		if i == cursor {
			prefix = cursorStyle.Render("> ")
		}

		nameStr := jokerStyle.Render(item.Def.Name)
		var priceStr string
		if gold < item.Price {
			priceStr = warningStyle.Render(fmt.Sprintf("$%d", item.Price))
		} else {
			priceStr = goldCardStyle.Render(fmt.Sprintf("$%d", item.Price))
		}

		sb.WriteString(fmt.Sprintf("%s%s  %s\n", prefix, nameStr, priceStr))
		sb.WriteString(fmt.Sprintf("   %s\n", item.Def.Description))
		sb.WriteString("\n")
	}

	sb.WriteString("[↑ ↓] 移动  [b] 购买  [n] 下一关\n")
	return sb.String()
}
