# TODO.md — 原子任务列表

规则：奇数任务写测试，偶数任务写实现。每个任务只改一个文件。

---

## Phase 1 — 可运行骨架

### T01 🧪 `engine/card_test.go`
测试 Card 的 String() 输出格式：
- `[A♠]` `[10♥]` `[2♦]` `[K♣]`
- 红心/方块的花色符号颜色标记（验证 IsRed()）
- Gold Card 的 String() 包含 `$` 标记

### T02 ✏️ `engine/card.go`
实现：
- `Suit`、`Rank`、`Enhancement` 枚举及常量
- `Card` 结构体，`Chips()` 返回牌面基础筹码（2-9=面值，10/J/Q/K=10，A=11）
- `String()` 返回 `[A♠]` 格式
- `IsRed()`、`IsGold()`

### T03 🧪 `engine/deck_test.go`
测试 Deck：
- 新建 Deck 包含且仅包含 52 张不重复的牌
- 相同 seed 洗牌结果完全一致
- 不同 seed 洗牌结果不同
- 发牌后牌堆数量相应减少
- 牌堆空时发牌返回空切片不 panic

### T04 ✏️ `engine/deck.go`
实现：
- `Deck` 结构体（包含 `[]Card` 和 `*rand.Rand`）
- `NewDeck(seed int64) Deck`：52 张标准牌
- `Shuffle()`：Fisher-Yates，使用 seed 初始化的 rng
- `Deal(n int) []Card`：从顶部取 n 张

### T05 🧪 `engine/hand_test.go`
测试 Evaluate()，覆盖全部 9 种手牌类型：
- 5 张同花顺：`[A♠][K♠][Q♠][J♠][10♠]` → StraightFlush
- 4 张四条：`[A♠][A♥][A♦][A♣][2♠]` → FourOfAKind
- 葫芦：`[K♠][K♥][K♦][A♣][A♠]` → FullHouse
- 同花（非顺）：`[2♠][5♠][7♠][9♠][J♠]` → Flush
- 顺子（非同花）：`[A♠][2♥][3♦][4♣][5♠]` → Straight（A低）
- 顺子：`[10♠][J♥][Q♦][K♣][A♠]` → Straight（A高）
- 三条：`[7♠][7♥][7♦][2♣][5♠]` → ThreeOfAKind
- 两对：`[J♠][J♥][9♦][9♣][3♠]` → TwoPair
- 一对：`[4♠][4♥][K♦][7♣][2♠]` → OnePair
- 高牌：`[A♠][K♥][Q♦][J♣][9♠]` → HighCard
- 打出 2 张一对：`[A♠][A♥]` → OnePair
- 打出 1 张：`[A♠]` → HighCard
- ScoringCards 正确区分触发牌和 Kicker

### T06 ✏️ `engine/hand.go`
实现：
- `HandType` 枚举（9 种，含基础 Chips/Mult 表）
- `EvaluateResult` 结构体
- `Evaluate(cards []Card) EvaluateResult`
- `BaseChips(h HandType) int`、`BaseMult(h HandType) int`
- A 同时作为最高/最低牌的顺子识别

### T07 🧪 `engine/score_test.go`
测试无 Joker 时的计分（ScoreHand）：
- 一对 A：基础 Chips=10，基础 Mult=2，两张 A 各贡献 11 chips → Total=(10+11+11)*2=64
- 同花顺 5 张：验证所有 ScoringCards 都贡献筹码
- ScoreResult.Steps 数量正确（每张触发牌一步 + 手牌类型一步）
- ScoreResult.Total == FinalChips * FinalMult

### T08 ✏️ `engine/score.go`
实现：
- `ScoreContext` 结构体
- `ScoreStep` 结构体
- `ScoreResult` 结构体
- `ScoreHand(cards []Card, run *RunState) ScoreResult`
  - 从 EvaluateResult 取基础值
  - 遍历 ScoringCards，累加 Chips，触发 OnCardScored Hook
  - 触发 OnHandScored Hook
  - 每步追加到 Steps

### T09 🧪 `engine/joker_test.go`
测试 Hook 系统框架：
- 空 Joker 列表时 ScoreHand 结果与无 Joker 相同
- 注册一个 +10 Mult 的 mock Joker，验证计分结果 Mult 增加 10
- 注册两个 Joker，验证按顺序触发
- Retrigger=1 时当前牌计分触发两次，Chips 翻倍

### T10 ✏️ `engine/joker.go`
实现：
- `GameEvent` 枚举
- `HookResult` 结构体
- `JokerDef` 结构体
- `JokerRegistry`：`Register()`、`All()`、`ByID()`
- 全局 `DefaultRegistry`
- 修改 `score.go` 的 ScoreHand 调用 Hook 链（此处允许跨文件，但只改 joker.go 接口，score.go 在 T08 中已预留 Hook 调用位）

### T11 🧪 `engine/jokers/registry_test.go`
测试 Joker 注册表：
- `All()` 返回非空列表
- `ByID()` 能找到已注册的 Joker
- `ByID()` 对未知 ID 返回 false

### T12 ✏️ `engine/jokers/registry.go`
实现：
- `init()` 注册所有内置 Joker 到 DefaultRegistry

### T13 🧪 `engine/jokers/greedy_test.go`
测试贪心小丑（OnCardScored，打出♦时 +3 Mult）：
- 打出 `[A♦]` 时 Mult +3
- 打出 `[A♠]` 时 Mult 不变
- 打出 3 张♦时 Mult +9

### T14 ✏️ `engine/jokers/greedy.go`
实现贪心小丑 JokerDef，注册到 DefaultRegistry。

### T15 🧪 `engine/jokers/half_test.go`
测试半张（OnHandScored，手牌≤3张时 +20 Mult）：
- 打出 3 张牌时 Mult +20
- 打出 5 张牌时 Mult 不变
- 打出 1 张牌时 Mult +20

### T16 ✏️ `engine/jokers/half.go`
实现半张 JokerDef，注册到 DefaultRegistry。

### T17 🧪 `engine/jokers/retrigger_test.go`
测试一张覆盖 Retrigger 的 Joker（如：重复触发最后一张计分牌）：
- 有 Retrigger Joker 时，最后一张 ScoringCard 的 Chips 贡献出现两次
- ScoreResult.Steps 包含对应的重复触发记录

### T18 ✏️ `engine/jokers/retrigger.go`
实现一个 Retrigger Joker（如"袜子和半身裙"：重复触发最后打出的牌）。

### T19 🧪 `game/state_test.go`
测试 GameState 初始化：
- `NewGame(seed)` 返回合法初始状态
  - Ante=1，BlindIndex=0，Gold=4
  - HandsLeft=4，DiscardsLeft=4，HandSize=8
  - Hand 包含 8 张牌，Deck 包含 44 张
  - JokerSlots=5，Jokers 为空
- 相同 seed 两次 NewGame 手牌相同

### T20 ✏️ `game/state.go`
实现：
- `GameState`、`RunState`、`RoundState`、`OwnedJoker` 结构体
- `NewGame(seed int64) GameState`
- `IsGameOver(state GameState) bool`
- `IsVictory(state GameState) bool`（通过 Ante 8 Boss Blind）
- `CurrentTarget(state GameState) int`（基于 Ante/BlindIndex 返回目标分数，使用原版数值）

### T21 🧪 `game/actions_test.go`
测试 PlayHand：
- 选中 5 张牌打出，HandsLeft -1，手牌补充至 8 张
- Score 增加对应计分结果的 Total
- Score ≥ Target 后 IsBlindComplete() 返回 true
- HandsLeft=0 且 Score < Target 时 IsGameOver 返回 true
- 选中超过 5 张返回 error
- 选中 0 张返回 error

测试 Discard：
- 选中 3 张弃牌，DiscardsLeft -1，手牌补充至 8 张
- DiscardsLeft=0 时返回 error
- 选中 0 张返回 error

测试 Gold Card：
- 手牌中有 Gold Card，打出时 Gold +1

### T22 ✏️ `game/actions.go`
实现：
- `PlayHand(state *GameState, selectedIdx []int) (engine.ScoreResult, error)`
- `Discard(state *GameState, selectedIdx []int) error`
- `IsBlindComplete(state GameState) bool`
- `NextBlind(state *GameState) error`（推进 BlindIndex/Ante，发放金币奖励）
- `BuyJoker(state *GameState, jokerID string) error`

### T23 🧪 `engine/save_test.go`
测试存读档：
- SaveGame 生成合法 JSON 文件
- LoadGame 读回后与原始 GameState 深度相等
- LoadGame 读取损坏 JSON 返回 error，不 panic
- LoadGame 文件不存在返回 error，不 panic

### T24 ✏️ `engine/save.go`
实现：
- `SaveGame(state GameState, path string) error`（原子写：先写临时文件再 rename）
- `LoadGame(path string) (GameState, error)`
- `SavePath() string`（返回 `~/.config/balatro-cli/save.json`）

### T25 🧪 `engine/blind_test.go`
测试 Blind 目标分数：
- Ante 1 Small Blind = 300
- Ante 1 Big Blind = 450
- Ante 8 Boss Blind = 原版对应数值
- BlindReward() 返回正确金币奖励

### T26 ✏️ `engine/blind.go`
实现：
- `BlindType` 枚举（Small/Big/Boss）
- `BlindDef` 结构体（Name、Target、Reward）
- `BlindAt(ante int, blindIndex int) BlindDef`（查原版数值表）
- `BossBlindEffect` 接口（预留，当前实现为空 no-op）

### T27 🧪 `engine/shop_test.go`
测试商店：
- `GenerateShop(rng)` 返回 2 个 Joker
- 相同 rng seed 生成相同商品
- 商品价格 > 0

### T28 ✏️ `engine/shop.go`
实现：
- `ShopItem` 结构体（JokerDef + Price）
- `GenerateShop(rng *rand.Rand) []ShopItem`

### T29 🧪 `main_test.go`
集成测试：
- `--seed 12345` 启动，模拟出牌序列，验证最终得分可复现
- 存档保存后重新加载，游戏状态一致

### T30 ✏️ `main.go`
实现命令行入口（Phase 1 为 stdin 交互版）：
- 解析 `--seed` 参数
- 检测存档，提示继续/新游戏
- 主循环：显示手牌 → 读取输入 → 执行动作 → 显示计分步骤
- Ctrl+C 二次确认退出并保存

---

## Phase 2 — TUI 层

### T31 🧪 `tui/styles_test.go`
测试样式常量：
- 验证颜色常量非空（lipgloss.Color 合法）
- CursorStyle、SelectedStyle、RedSuitStyle 可区分

### T32 ✏️ `tui/styles.go`
实现 lipgloss 样式常量：
- `CursorStyle`（青色）
- `SelectedStyle`（黄色）
- `RedSuitStyle`（红色，红心/方块）
- `JokerStyle`（紫色）
- `ScoreStyle`（绿色）
- `GoldStyle`（黄色）
- `WarningStyle`（暗红）

### T33 🧪 `tui/views/hand_test.go`
测试手牌 View 渲染：
- 8 张手牌全部出现在输出字符串中
- 光标位置的牌包含 CursorStyle ANSI 序列
- 已选中牌包含 SelectedStyle ANSI 序列
- 状态栏显示正确的 Ante/Gold/HandsLeft/DiscardsLeft

### T34 ✏️ `tui/views/hand.go`
实现手牌界面 View 函数（纯函数，接收 GameState 返回 string）。

### T35 🧪 `tui/views/score_test.go`
测试计分 View：
- 每个 ScoreStep 占一行
- 最终 Total 行存在
- ✓/✗ 标记根据是否达标正确显示

### T36 ✏️ `tui/views/score.go`
实现计分过程 View，将 ScoreResult.Steps 渲染为逐行文本。

### T37 🧪 `tui/views/shop_test.go`
测试商店 View：
- 两个商品各自显示名称、描述、价格
- 金币不足的商品有视觉区分（WarningStyle）

### T38 ✏️ `tui/views/shop.go`
实现商店界面 View。

### T39 🧪 `tui/views/gameover_test.go`
测试结束 View：
- 胜利/失败文案正确
- 显示坚持到的 Ante 和 Seed

### T40 ✏️ `tui/views/gameover.go`
实现游戏结束界面 View。

### T41 🧪 `tui/views/help_test.go`
测试帮助 View：
- 包含所有按键说明（←→ Space p d ? Ctrl+C）

### T42 ✏️ `tui/views/help.go`
实现帮助页面 View。

### T43 🧪 `tui/model_test.go`
测试 bubbletea Model：
- 初始 ViewType 为 ViewHand
- 按 `?` 切换到 ViewHelp，再按任意键返回 ViewHand
- 按 `p`（无选中牌）状态不变，显示警告
- 第一次 Ctrl+C 设置 pendingQuit=true
- 第二次 Ctrl+C 在 pendingQuit=true 时返回 tea.Quit

### T44 ✏️ `tui/model.go`
实现 bubbletea Model：
- `Model` 结构体（GameState、ViewType、Cursor、Selected、pendingQuit、lastWarning）
- `Init()`、`Update()`、`View()` 实现
- 键盘事件路由到对应 View 的 Update 逻辑
- `tui.Start(state game.GameState) error`

---

## Phase 3 — 完善

### T45 🧪 `engine/blind_test.go`（追加）
测试 BossBlindEffect 接口：
- 当前 no-op 实现不修改任何状态
- 接口方法签名固定，未来实现可替换

### T46 ✏️ `engine/blind.go`（追加）
完善 BossBlindEffect 接口定义，no-op 默认实现。

### T47 🧪 `engine/jokers/discard_test.go`
测试弃牌触发型 Joker（OnDiscard Hook）：
- 弃牌时 Hook 被调用
- 弃 3 张牌触发 3 次

### T48 ✏️ `engine/jokers/discard.go`
实现一个 OnDiscard Joker（如"残缺的旅行家"：每次弃牌 +$1）。

### T49 🧪 `engine/jokers/roundend_test.go`
测试回合结束触发型 Joker（OnRoundEnd Hook）：
- 通过 Blind 时 Hook 被调用一次
- 效果正确应用到 RunState

### T50 ✏️ `engine/jokers/roundend.go`
实现一个 OnRoundEnd Joker（如"蓝筹股"：每回合结束 +$1）。
