package views

// RenderHelp 渲染帮助页面，返回纯字符串
func RenderHelp() string {
	return `=== 帮助 Help ===

按键说明 Key Bindings:
  ←  /  →     移动光标 Move cursor
  Space        选择牌 Select card
  Backspace    取消最后选择 Deselect last
  Enter / p    出牌 Play selected cards
  d            弃牌 Discard selected cards
  ?            显示/隐藏帮助 Toggle help
  Ctrl+C       退出（按两次确认）Quit (press twice to confirm)

商店 Shop:
  ↑ / ↓        移动光标 Move cursor
  b            购买 Buy
  n            下一关 Next blind

按任意键返回 Press any key to return
`
}
