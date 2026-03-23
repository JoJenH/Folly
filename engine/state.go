package engine

// OwnedJoker 玩家持有的 Joker 实例
type OwnedJoker struct {
	DefID string // 对应 JokerDef.ID
}

// RoundState 单个 Blind 回合状态
type RoundState struct {
	Deck         []Card
	Hand         []Card
	HandSize     int
	HandsLeft    int
	DiscardsLeft int
	Score        int
	Target       int
}

// RunState 一局游戏（Ante 1-8）的持久状态
type RunState struct {
	Ante       int
	BlindIndex int // 0=Small 1=Big 2=Boss
	Gold       int
	Jokers     []OwnedJoker
	JokerSlots int
	Round      RoundState
}

// GameState 完整游戏状态（可序列化）
type GameState struct {
	Seed int64
	Run  RunState
}
