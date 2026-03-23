package engine

// JokerDef 单个 Joker 定义
type JokerDef struct {
	ID          string
	Name        string
	Description string
	Cost        int
	Hooks       map[GameEvent]func(ctx *ScoreContext) HookResult
}

// ApplyHook 实现 JokerHook 接口
func (j *JokerDef) ApplyHook(event GameEvent, ctx *ScoreContext) HookResult {
	if j.Hooks == nil {
		return HookResult{}
	}
	if fn, ok := j.Hooks[event]; ok {
		return fn(ctx)
	}
	return HookResult{}
}

// JokerRegistry Joker 注册表
type JokerRegistry struct {
	jokers []*JokerDef
	byID   map[string]*JokerDef
}

// Register 注册一个 Joker
func (r *JokerRegistry) Register(j *JokerDef) {
	if r.byID == nil {
		r.byID = make(map[string]*JokerDef)
	}
	r.jokers = append(r.jokers, j)
	r.byID[j.ID] = j
}

// All 返回所有已注册的 Joker
func (r *JokerRegistry) All() []*JokerDef {
	return r.jokers
}

// ByID 按 ID 查找 Joker
func (r *JokerRegistry) ByID(id string) (*JokerDef, bool) {
	j, ok := r.byID[id]
	return j, ok
}

// DefaultRegistry 全局默认注册表
var DefaultRegistry JokerRegistry
