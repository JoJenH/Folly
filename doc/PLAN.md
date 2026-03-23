# PLAN.md — 技术方案

## 技术栈

- **语言**：Go 1.22+
- **TUI 框架**：[bubbletea](https://github.com/charmbracelet/bubbletea) + [lipgloss](https://github.com/charmbracelet/lipgloss)
- **存档格式**：JSON（encoding/json 标准库）
- **随机**：`math/rand/v2` with seed

---

## 目录结构

```
balatro-cli/
├── main.go                  # 入口，解析 --seed 参数，启动 TUI
├── go.mod
├── go.sum
│
├── engine/                  # 纯游戏逻辑，无副作用，无 TUI 依赖
│   ├── card.go              # Card、Suit、Rank、Enhancement 定义
│   ├── deck.go              # Deck：洗牌、发牌、补牌
│   ├── hand.go              # 手牌类型判定（HandType、Evaluate）
│   ├── score.go             # 计分引擎：ScoreContext、ScoreStep、ScoreHand
│   ├── joker.go             # Joker 定义：GameEvent、Hook、JokerDef
│   ├── jokers/              # 每个 Joker 一个文件，注册到全局 Registry
│   │   ├── registry.go      # JokerRegistry，All()、ByID()
│   │   ├── greedy.go        # 贪心小丑
│   │   ├── half.go          # 半张
│   │   └── ...
│   ├── blind.go             # Blind 定义：SmallBlind/BigBlind/BossBlind，目标分数
│   ├── shop.go              # 商店逻辑：生成商品、购买
│   └── save.go              # 序列化/反序列化 GameState
│
├── game/                    # 游戏状态机
│   ├── state.go             # GameState、RunState、RoundState
│   └── actions.go           # PlayHand、Discard、BuyJoker、NextBlind 等动作
│
├── tui/                     # TUI 层，只依赖 game 包
│   ├── model.go             # bubbletea Model，顶层 Update/View 路由
│   ├── views/
│   │   ├── hand.go          # 手牌界面 View
│   │   ├── score.go         # 计分过程 View（逐步展示）
│   │   ├── shop.go          # 商店界面 View
│   │   ├── gameover.go      # 游戏结束界面
│   │   └── help.go          # 帮助页面
│   └── styles.go            # lipgloss 颜色、样式常量
│
└── testdata/                # 测试用固定牌局数据
    └── hands.json
```

---

## 核心数据模型

### Card

```go
type Suit int
const (SuitSpade Suit = iota; SuitHeart; SuitDiamond; SuitClub)

type Rank int
const (Rank2 Rank = 2; ...; RankAce = 14)

type Enhancement int
const (EnhancementNone Enhancement = iota; EnhancementGold)

type Card struct {
    Suit        Suit
    Rank        Rank
    Enhancement Enhancement
    // 牌的基础筹码值
    Chips       int
}

func (c Card) String() string  // "[A♠]"
func (c Card) IsGold() bool
```

### HandType & Evaluate

```go
type HandType int
const (
    HighCard HandType = iota
    OnePair; TwoPair; ThreeOfAKind
    Straight; Flush; FullHouse
    FourOfAKind; StraightFlush
)

type EvaluateResult struct {
    Type          HandType
    // 参与计分的牌（触发牌）
    ScoringCards  []Card
    // 未触发的牌
    KickerCards   []Card
    BaseChips     int
    BaseMult      int
}

func Evaluate(cards []Card) EvaluateResult
```

### ScoreContext & ScoreStep

```go
// 在 Hook 链中流动的上下文，所有 Joker 读写此结构
type ScoreContext struct {
    Chips       int
    Mult        int
    Card        *Card      // 当前触发的牌（OnCardScored 时非空）
    ScoringCards []Card
    HandType    HandType
    Round       *RoundState
    Run         *RunState
}

// 计分过程中的单条记录，用于 TUI 逐步展示
type ScoreStep struct {
    Description string   // "[K♥] +10 chips" / "[贪心小丑] Mult +4"
    ChipsAfter  int
    MultAfter   int
}

type ScoreResult struct {
    Steps       []ScoreStep
    FinalChips  int
    FinalMult   int
    Total       int          // FinalChips * FinalMult
}

func ScoreHand(cards []Card, run *RunState) ScoreResult
```

### Joker

```go
type GameEvent int
const (
    OnCardScored GameEvent = iota
    OnHandScored
    OnDiscard
    OnRoundEnd
)

type HookResult struct {
    ChipsDelta  int
    MultDelta   int
    Retrigger   int    // 重复触发当前牌几次
}

type JokerDef struct {
    ID          string
    Name        string
    Description string
    Cost        int
    Hooks       map[GameEvent]func(ctx *ScoreContext) HookResult
}
```

### GameState

```go
// 完整游戏状态，可序列化为 JSON
type GameState struct {
    Seed        int64
    Run         RunState
}

type RunState struct {
    Ante        int           // 1-8
    BlindIndex  int           // 0=Small 1=Big 2=Boss
    Gold        int
    Jokers      []OwnedJoker
    JokerSlots  int           // 默认 5
    Round       RoundState
}

type RoundState struct {
    Deck        []Card        // 剩余牌堆
    Hand        []Card        // 当前手牌
    HandSize    int           // 默认 8
    HandsLeft   int
    DiscardsLeft int
    Score       int           // 本 Blind 累计得分
    Target      int           // 本 Blind 目标分数
}

type OwnedJoker struct {
    DefID   string    // 对应 JokerDef.ID
    // 未来扩展：Joker 自身状态（如计数器）
}
```

---

## 接口定义

### engine 包对外暴露

```go
// 判定手牌类型
func Evaluate(cards []Card) EvaluateResult

// 执行完整计分，返回逐步记录
func ScoreHand(cards []Card, run *RunState) ScoreResult

// Joker 注册表
func (r *JokerRegistry) All() []JokerDef
func (r *JokerRegistry) ByID(id string) (JokerDef, bool)

// 存档
func SaveGame(state GameState, path string) error
func LoadGame(path string) (GameState, error)
```

### game 包对外暴露（动作层）

```go
// 打出选中的牌，返回计分结果
func PlayHand(state *GameState, selectedIdx []int) (ScoreResult, error)

// 弃牌
func Discard(state *GameState, selectedIdx []int) error

// 购买 Joker
func BuyJoker(state *GameState, jokerID string) error

// 离开商店进入下一个 Blind
func NextBlind(state *GameState) error

// 初始化新游戏
func NewGame(seed int64) GameState

// 判断游戏是否结束
func IsGameOver(state GameState) bool
func IsVictory(state GameState) bool
```

### tui 包

```go
// bubbletea 入口
func Start(state game.GameState) error

// 内部 ViewType 枚举
type ViewType int
const (
    ViewHand ViewType = iota
    ViewScore
    ViewShop
    ViewGameOver
    ViewHelp
)
```

---

## 颜色方案（lipgloss）

| 元素 | 颜色 |
|------|------|
| 光标位置牌 | 青色 Cyan |
| 已选中牌 | 黄色 Yellow |
| 黑桃/梅花 | 默认前景色 |
| 红心/方块 | 红色 |
| Joker 名称 | 紫色 |
| 得分数字 | 绿色 |
| 警告/不可操作提示 | 暗红色 |
| 金币 | 黄色 |

---

## 实施阶段

### Phase 1 — 可运行骨架（无 TUI）

目标：命令行交互版，能完整跑通一局游戏的核心循环。

- `engine/card.go`：Card、Suit、Rank 定义及 String()
- `engine/deck.go`：标准 52 张牌，洗牌（seed），发牌
- `engine/hand.go`：9 种手牌类型判定
- `engine/score.go`：基础计分（无 Joker）
- `engine/joker.go`：Hook 系统框架
- `engine/jokers/`：3-5 张覆盖全部 Hook 类型的 Joker
- `game/state.go`：GameState 定义
- `game/actions.go`：PlayHand、Discard、NewGame
- `engine/save.go`：JSON 存读档
- `main.go`：stdin 交互，数字选牌，验证逻辑正确

### Phase 2 — TUI 层

目标：bubbletea 接管交互，视觉呈现完整。

- `tui/styles.go`：lipgloss 样式定义
- `tui/model.go`：顶层 Model，ViewType 路由
- `tui/views/hand.go`：手牌界面
- `tui/views/score.go`：逐步计分展示
- `tui/views/shop.go`：商店界面
- `tui/views/gameover.go`：结束画面
- `tui/views/help.go`：帮助页面

### Phase 3 — 完善与扩展

目标：补齐剩余规则，为内容扩展做准备。

- `engine/blind.go`：Boss Blind 接口（行为占位，接口完整）
- `engine/jokers/`：补充更多 Joker
- 商店商品扩展接口（塔罗牌、星球牌预留）
- `--seed` 参数支持
- Ctrl+C 二次确认退出
