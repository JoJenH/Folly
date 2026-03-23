package tui

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"balatro-cli/engine"
	_ "balatro-cli/engine/jokers"
	"balatro-cli/game"
	"balatro-cli/tui/views"

	tea "github.com/charmbracelet/bubbletea"
)

// savePath 存档路径（避免循环导入，直接调用 engine）
var savePath = engine.SavePath()

type scoreAdvanceMsg struct{}
type scoreFlashMsg struct{}

func scoreAdvanceCmd() tea.Cmd {
	return tea.Tick(750*time.Millisecond, func(time.Time) tea.Msg { return scoreAdvanceMsg{} })
}
func scoreFlashCmd() tea.Cmd {
	return tea.Tick(330*time.Millisecond, func(time.Time) tea.Msg { return scoreFlashMsg{} })
}

// ViewType 当前显示的界面类型
type ViewType int

const (
	ViewHand ViewType = iota
	ViewScore
	ViewShop
	ViewGameOver
	ViewHelp
)

// Model bubbletea 模型
type Model struct {
	state         game.GameState
	view          ViewType
	cursor        int
	selected      []int
	pendingQuit   bool
	lastWarning   string
	scoreResult   engine.ScoreResult
	shopItems     []engine.ShopItem
	scoreAnimStep int
	scoreFlash    bool
	scoreAnimDone bool
}

// NewModel 创建新模型
func NewModel(state game.GameState) Model {
	return Model{
		state:    state,
		view:     ViewHand,
		cursor:   0,
		selected: []int{},
	}
}

// Init 实现 tea.Model 接口
func (m Model) Init() tea.Cmd {
	return nil
}

// Update 实现 tea.Model 接口
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case scoreAdvanceMsg:
		if m.view == ViewScore && !m.scoreAnimDone {
			m.scoreAnimStep++
			if m.scoreAnimStep >= len(m.scoreResult.Steps) {
				m.scoreAnimDone = true
				m.scoreFlash = false
				return m, nil
			}
			return m, scoreAdvanceCmd()
		}
	case scoreFlashMsg:
		if m.view == ViewScore && !m.scoreAnimDone {
			m.scoreFlash = !m.scoreFlash
			return m, scoreFlashCmd()
		}
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Ctrl+C 双击确认退出
	if msg.Type == tea.KeyCtrlC {
		if m.pendingQuit {
			engine.SaveGame(m.state, savePath)
			return m, tea.Quit
		}
		m.pendingQuit = true
		m.lastWarning = "再按一次 Ctrl+C 退出 Press Ctrl+C again to quit"
		return m, nil
	}
	m.pendingQuit = false

	switch m.view {
	case ViewHelp:
		m.view = ViewHand
		return m, nil

	case ViewScore:
		if !m.scoreAnimDone {
			// 任意键跳过动画
			m.scoreAnimStep = len(m.scoreResult.Steps)
			m.scoreAnimDone = true
			m.scoreFlash = false
			return m, nil
		}
		// 动画结束后任意键继续
		if game.IsBlindComplete(m.state) {
			return m.enterShop()
		}
		if game.IsGameOver(m.state) {
			m.view = ViewGameOver
			return m, nil
		}
		m.view = ViewHand
		return m, nil

	case ViewShop:
		return m.handleShopKey(msg)

	case ViewGameOver:
		os.Remove(savePath)
		return m, tea.Quit

	case ViewHand:
		return m.handleHandKey(msg)
	}
	return m, nil
}

func (m Model) handleHandKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	hand := m.state.Run.Round.Hand

	// 计算未选中牌的实际手牌索引列表
	selSet := make(map[int]bool, len(m.selected))
	for _, i := range m.selected {
		selSet[i] = true
	}
	unselIdx := make([]int, 0, len(hand))
	for i := range hand {
		if !selSet[i] {
			unselIdx = append(unselIdx, i)
		}
	}

	switch msg.Type {
	case tea.KeyLeft:
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case tea.KeyRight:
		if m.cursor < len(unselIdx)-1 {
			m.cursor++
		}
		return m, nil
	case tea.KeySpace:
		if m.cursor < len(unselIdx) {
			actualIdx := unselIdx[m.cursor]
			m.selected = append(m.selected, actualIdx)
			// 光标不超出新的未选中列表范围
			if m.cursor >= len(unselIdx)-1 && m.cursor > 0 {
				m.cursor--
			}
		}
		return m, nil
	case tea.KeyBackspace, tea.KeyDelete:
		// 取消最后一张选中的牌
		if len(m.selected) > 0 {
			m.selected = m.selected[:len(m.selected)-1]
		}
		return m, nil
	}

	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
		switch msg.Runes[0] {
		case '?':
			m.view = ViewHelp
			return m, nil
		case 'p':
			// 保留 p 键向后兼容，同 Enter
			fallthrough
		case '\r':
			if len(m.selected) == 0 {
				m.lastWarning = "请先选择牌 Please select cards first"
				return m, nil
			}
			result, err := game.PlayHand(&m.state, m.selected)
			if err != nil {
				m.lastWarning = err.Error()
				return m, nil
			}
			m.scoreResult = result
			m.selected = []int{}
			m.cursor = 0
			m.lastWarning = ""
			m.view = ViewScore
			m.scoreAnimStep = 0
			m.scoreAnimDone = false
			m.scoreFlash = true
			return m, tea.Batch(scoreAdvanceCmd(), scoreFlashCmd())
		case 'd':
			if len(m.selected) == 0 {
				m.lastWarning = "请先选择牌 Please select cards first"
				return m, nil
			}
			err := game.Discard(&m.state, m.selected)
			if err != nil {
				m.lastWarning = err.Error()
				return m, nil
			}
			m.selected = []int{}
			m.cursor = 0
			m.lastWarning = ""
			return m, nil
		}
	}

	// Enter 键出牌
	if msg.Type == tea.KeyEnter {
		if len(m.selected) == 0 {
			m.lastWarning = "请先选择牌 Please select cards first"
			return m, nil
		}
		result, err := game.PlayHand(&m.state, m.selected)
		if err != nil {
			m.lastWarning = err.Error()
			return m, nil
		}
		m.scoreResult = result
		m.selected = []int{}
		m.cursor = 0
		m.lastWarning = ""
		m.view = ViewScore
		m.scoreAnimStep = 0
		m.scoreAnimDone = false
		m.scoreFlash = true
		return m, tea.Batch(scoreAdvanceCmd(), scoreFlashCmd())
	}

	return m, nil
}

func (m Model) handleShopKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.cursor < len(m.shopItems)-1 {
			m.cursor++
		}
		return m, nil
	}

	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
		switch msg.Runes[0] {
		case 'b':
			if m.cursor < len(m.shopItems) {
				err := game.BuyJoker(&m.state, m.shopItems[m.cursor].Def.ID)
				if err != nil {
					m.lastWarning = err.Error()
				} else {
					m.lastWarning = ""
					// 移除已购买商品
					if m.cursor < len(m.shopItems) {
						m.shopItems = append(m.shopItems[:m.cursor], m.shopItems[m.cursor+1:]...)
						if m.cursor >= len(m.shopItems) && m.cursor > 0 {
							m.cursor--
						}
					}
				}
			}
			return m, nil
		case 'n':
			err := game.NextBlind(&m.state)
			if err != nil {
				m.lastWarning = err.Error()
				return m, nil
			}
			if game.IsVictory(m.state) {
				m.view = ViewGameOver
				return m, nil
			}
			m.cursor = 0
			m.selected = []int{}
			m.lastWarning = ""
			m.view = ViewHand
			return m, nil
		}
	}
	return m, nil
}

func (m Model) enterShop() (tea.Model, tea.Cmd) {
	// 生成商店商品
	deckRng := rand.New(rand.NewPCG(uint64(m.state.Seed), uint64(m.state.Run.Ante)))
	m.shopItems = engine.GenerateShop(deckRng)
	m.cursor = 0
	m.view = ViewShop
	return m, nil
}

// View 实现 tea.Model 接口
func (m Model) View() string {
	var content string
	switch m.view {
	case ViewHand:
		content = views.RenderHand(m.state, m.cursor, m.selected)
	case ViewScore:
		content = views.RenderScoreAnimated(
			m.scoreResult,
			m.state.Run.Round.Score,
			m.state.Run.Round.Target,
			m.scoreAnimStep,
			m.scoreFlash,
		)
	case ViewShop:
		content = views.RenderShop(m.shopItems, m.state.Run.Gold, m.cursor)
	case ViewGameOver:
		content = views.RenderGameOver(m.state, game.IsVictory(m.state))
	case ViewHelp:
		content = views.RenderHelp()
	}

	if m.lastWarning != "" {
		content += WarningStyle.Render(fmt.Sprintf("⚠ %s", m.lastWarning)) + "\n"
	}
	return content
}

// Start 启动 bubbletea TUI
func Start(state game.GameState) error {
	p := tea.NewProgram(NewModel(state))
	_, err := p.Run()
	return err
}
