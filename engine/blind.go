package engine

// BlindType Blind 类型
type BlindType int

const (
	BlindSmall BlindType = iota
	BlindBig
	BlindBoss
)

// BlindDef Blind 定义
type BlindDef struct {
	Name     string
	Type     BlindType
	Target   int
	Reward   int
	Effect   BossBlindEffect
}

// BossBlindEffect Boss Blind 特殊效果接口（预留，当前为 no-op）
type BossBlindEffect interface {
	Apply(state *RunState)
}

// noopBossEffect 默认 no-op 实现
type noopBossEffect struct{}

func (n noopBossEffect) Apply(state *RunState) {}

// BlindAt 根据 ante（1-8）和 blindIndex（0=Small,1=Big,2=Boss）返回 BlindDef
func BlindAt(ante int, blindIndex int) BlindDef {
	type row struct{ small, big, boss int }
	targets := []row{
		{300, 450, 600},
		{800, 1200, 1600},
		{2000, 3000, 4000},
		{5000, 7500, 10000},
		{11000, 16500, 22000},
		{20000, 30000, 40000},
		{35000, 52500, 70000},
		{50000, 75000, 100000},
	}

	names := [][]string{
		{"小盲注", "大盲注", "首领盲注"},
	}

	if ante < 1 || ante > 8 {
		return BlindDef{}
	}

	r := targets[ante-1]
	var target, reward int
	var btype BlindType
	var name string

	switch blindIndex {
	case 0:
		target, reward, btype, name = r.small, 3, BlindSmall, "小盲注"
	case 1:
		target, reward, btype, name = r.big, 4, BlindBig, "大盲注"
	case 2:
		target, reward, btype, name = r.boss, 5, BlindBoss, "首领盲注"
	default:
		return BlindDef{}
	}

	_ = names // 使用 switch 替代，避免未使用警告

	return BlindDef{
		Name:   name,
		Type:   btype,
		Target: target,
		Reward: reward,
		Effect: noopBossEffect{},
	}
}
