# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供此仓库的开发指引。

> 开发过程中，不应无故修改测试文件，如有必要修改应先行列出修改原因由用户确认，永运测试先行
> 国际化说明：代码标识符（变量名、函数名、类型名）使用英文；用户可见文本（TUI 显示内容、错误提示）应同时支持中文与英文，默认显示中文。文档（`.md` 文件）使用中文撰写。
> 文档中未提及的内容，比如规则细节，你可以参考官方小丑牌现行规则，并与用户确认

## 项目概述

基于终端的小丑牌（Balatro）核心玩法实现，使用 Go + bubbletea TUI。目标是可玩的核心循环，且架构支持内容扩展而无需重构引擎。

## 技术栈

- **语言**：Go 1.22+
- **TUI**：[bubbletea](https://github.com/charmbracelet/bubbletea) + [lipgloss](https://github.com/charmbracelet/lipgloss)
- **存档格式**：JSON（`encoding/json`）
- **随机数**：`math/rand/v2` with seed

## 常用命令

```bash
# 构建
go build ./...

# 运行
go run . [--seed <值>]

# 全量测试
go test ./...

# 单包测试
go test ./engine/...
go test ./game/...
go test ./tui/...

# 单个测试用例
go test ./engine/ -run TestEvaluate

# 静态检查
go vet ./...
```

存档路径：`~/.config/balatro-cli/save.json`

## 架构

三层严格单向依赖：

```
tui/ → game/ → engine/
```

### `engine/` — 纯游戏逻辑，无 I/O，无 TUI 依赖

- `card.go`：`Card`、`Suit`、`Rank`、`Enhancement` 类型。`Chips()` 返回基础筹码值（2–9=面值，10/J/Q/K=10，A=11）。
- `deck.go`：`Deck`，带 seed 的 Fisher-Yates 洗牌，`Deal(n)` 发牌。
- `hand.go`：`Evaluate(cards []Card) EvaluateResult` — 识别手牌类型，区分 `ScoringCards`（参与计分）与 `KickerCards`。
- `score.go`：`ScoreHand(cards []Card, run *RunState) ScoreResult` — 完整计分流水线，产出 `[]ScoreStep` 供 TUI 逐步展示。
- `joker.go`：Hook 系统 — `JokerDef` 持有 `map[GameEvent]func(ctx *ScoreContext) HookResult`，全局 `DefaultRegistry`。
- `jokers/`：每个 Joker 一个文件，通过 `init()` 自注册。新增 Joker 只需添加文件，无需修改引擎。
- `blind.go`：`BlindDef`（名称、目标分数、奖励），`BlindAt(ante, blindIndex)`，`BossBlindEffect` 接口（占位）。
- `shop.go`：`GenerateShop(rng) []ShopItem` — 随机生成 2 个 Joker 商品。
- `save.go`：`SaveGame` / `LoadGame` — 存档损坏或缺失时返回错误，不 panic。

### `game/` — 状态机与动作

- `state.go`：`GameState`（可序列化根状态）、`RunState`（ante/blind/金币/Joker）、`RoundState`（牌堆/手牌/计数器/得分）。`NewGame(seed)`、`IsGameOver`、`IsVictory`。
- `actions.go`：`PlayHand`、`Discard`、`BuyJoker`、`NextBlind`、`IsBlindComplete`。这些是唯一修改 `GameState` 的函数。

### `tui/` — 纯展示层，只依赖 `game/`

- `model.go`：bubbletea `Model`，`ViewType` 路由（`ViewHand`、`ViewScore`、`ViewShop`、`ViewGameOver`、`ViewHelp`）。持有 `cursor`、`selected`、`pendingQuit`、`lastWarning`。
- `views/`：各界面的纯函数 `(state, ...) → string`。
- `styles.go`：所有 lipgloss 颜色常量（青色=光标，黄色=选中/金币，红色=红心/方块，紫色=Joker，绿色=得分，暗红色=警告）。

## 核心设计模式

### Hook / 事件系统

Joker 效果通过 `GameEvent` hook 在 `ScoreHand` 中触发：
- `OnCardScored`：每张计分牌触发一次（`Retrigger > 0` 时额外触发）
- `OnHandScored`：所有牌计分后触发一次
- `OnDiscard`：每张弃牌触发一次
- `OnRoundEnd`：通过一个 Blind 时触发

`ScoreContext` 在 hook 链中传递可变的 `Chips`/`Mult`。每个 hook 返回 `HookResult{ChipsDelta, MultDelta, Retrigger}`。

### 计分步骤录制

`ScoreHand` 在产出最终总分的同时记录 `[]ScoreStep`，每步包含 `Description`、`ChipsAfter`、`MultAfter`。TUI 据此逐步回放动画——引擎不控制动画时序。

### 任务规范（见 TODO.md）

奇数任务写测试，偶数任务写实现，每个任务只改一个文件。

## 游戏规则参考

| 手牌类型 | 基础筹码 | 基础倍率 |
|---------|---------|--------|
| 同花顺 Straight Flush | 100 | 8 |
| 四条 Four of a Kind | 60 | 7 |
| 葫芦 Full House | 40 | 6 |
| 同花 Flush | 35 | 4 |
| 顺子 Straight | 30 | 4 |
| 三条 Three of a Kind | 30 | 3 |
| 两对 Two Pair | 20 | 2 |
| 一对 One Pair | 10 | 2 |
| 高牌 High Card | 5 | 1 |

Blind 奖励：小盲=$3，大盲=$4，Boss盲=$5。初始状态：$4 金币，4 次出牌，4 次弃牌，手牌上限 8 张，从 Ante 1 小盲开始。胜利条件：通过 Ante 8 Boss 盲。
