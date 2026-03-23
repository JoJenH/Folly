package engine

import "sort"

// HandType 手牌类型，数值越大优先级越高
type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

func (h HandType) String() string {
	names := []string{
		"高牌", "一对", "两对", "三条",
		"顺子", "同花", "葫芦", "四条", "同花顺",
	}
	if int(h) < len(names) {
		return names[h]
	}
	return "未知"
}

// BaseChips 返回手牌类型的基础筹码
func BaseChips(h HandType) int {
	return []int{5, 10, 20, 30, 30, 35, 40, 60, 100}[h]
}

// BaseMult 返回手牌类型的基础倍率
func BaseMult(h HandType) int {
	return []int{1, 2, 2, 3, 4, 4, 6, 7, 8}[h]
}

// EvaluateResult 手牌评估结果
type EvaluateResult struct {
	Type         HandType
	ScoringCards []Card // 参与计分的牌
	KickerCards  []Card // 踢牌（不参与计分）
	BaseChips    int
	BaseMult     int
}

// Evaluate 识别手牌类型，区分 ScoringCards 与 KickerCards
func Evaluate(cards []Card) EvaluateResult {
	if len(cards) == 0 {
		return EvaluateResult{Type: HighCard, BaseChips: BaseChips(HighCard), BaseMult: BaseMult(HighCard)}
	}

	// 按点数降序排列
	sorted := make([]Card, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rank > sorted[j].Rank
	})

	flush := isFlush(sorted)
	straight, straightCards := isStraight(sorted)

	if flush && straight {
		return result(StraightFlush, straightCards, nil)
	}

	counts := rankCounts(sorted)

	if h, scoring, kickers := findOfAKind(sorted, counts, 4); h {
		return result(FourOfAKind, scoring, kickers)
	}

	if h, scoring, kickers := findFullHouse(sorted, counts); h {
		return result(FullHouse, scoring, kickers)
	}

	if flush {
		return result(Flush, sorted, nil)
	}

	if straight {
		return result(Straight, straightCards, nil)
	}

	if h, scoring, kickers := findOfAKind(sorted, counts, 3); h {
		return result(ThreeOfAKind, scoring, kickers)
	}

	if h, scoring, kickers := findTwoPair(sorted, counts); h {
		return result(TwoPair, scoring, kickers)
	}

	if h, scoring, kickers := findOfAKind(sorted, counts, 2); h {
		return result(OnePair, scoring, kickers)
	}

	// 高牌：只有最高的一张参与计分
	return result(HighCard, sorted[:1], sorted[1:])
}

func result(t HandType, scoring, kickers []Card) EvaluateResult {
	return EvaluateResult{
		Type:         t,
		ScoringCards: scoring,
		KickerCards:  kickers,
		BaseChips:    BaseChips(t),
		BaseMult:     BaseMult(t),
	}
}

// rankCounts 统计每种点数的数量
func rankCounts(cards []Card) map[Rank]int {
	counts := make(map[Rank]int)
	for _, c := range cards {
		counts[c.Rank]++
	}
	return counts
}

// isFlush 是否同花
func isFlush(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	suit := cards[0].Suit
	for _, c := range cards[1:] {
		if c.Suit != suit {
			return false
		}
	}
	return true
}

// isStraight 是否顺子，返回参与计分的牌（支持 A 高低）
func isStraight(cards []Card) (bool, []Card) {
	if len(cards) < 5 {
		return false, nil
	}
	// 正常顺子
	if isConsecutive(cards) {
		return true, cards
	}
	// A 作为低牌：A-2-3-4-5
	// 把 A(14) 视为 1，检查最低的 4 张是否为 2-3-4-5
	if cards[0].Rank == RankAce {
		rest := cards[1:] // 已降序排列，前 4 张应为 5,4,3,2
		expected := []Rank{Rank5, Rank4, Rank3, Rank2}
		if len(rest) >= 4 {
			ok := true
			for i, r := range expected {
				if rest[i].Rank != r {
					ok = false
					break
				}
			}
			if ok {
				// A-low straight：A 排最后
				return true, append(rest[:4:4], cards[0])
			}
		}
	}
	return false, nil
}

func isConsecutive(cards []Card) bool {
	for i := 1; i < len(cards); i++ {
		if int(cards[i-1].Rank)-int(cards[i].Rank) != 1 {
			return false
		}
	}
	return true
}

// findOfAKind 查找 n 条（2/3/4 条），返回触发牌和踢牌
func findOfAKind(sorted []Card, counts map[Rank]int, n int) (bool, []Card, []Card) {
	var scoring, kickers []Card
	found := false
	for _, c := range sorted {
		if counts[c.Rank] == n && !found {
			scoring = append(scoring, c)
			if len(scoring) == n {
				found = true
			}
		} else {
			kickers = append(kickers, c)
		}
	}
	if !found {
		return false, nil, nil
	}
	return true, scoring, kickers
}

// findFullHouse 三条+一对
func findFullHouse(sorted []Card, counts map[Rank]int) (bool, []Card, []Card) {
	var threeRank, pairRank Rank
	hasThree, hasPair := false, false
	for rank, cnt := range counts {
		if cnt == 3 {
			threeRank = rank
			hasThree = true
		} else if cnt == 2 {
			pairRank = rank
			hasPair = true
		}
	}
	if !hasThree || !hasPair {
		return false, nil, nil
	}
	var scoring []Card
	for _, c := range sorted {
		if c.Rank == threeRank || c.Rank == pairRank {
			scoring = append(scoring, c)
		}
	}
	return true, scoring, nil
}

// findTwoPair 两对
func findTwoPair(sorted []Card, counts map[Rank]int) (bool, []Card, []Card) {
	var pairRanks []Rank
	for rank, cnt := range counts {
		if cnt == 2 {
			pairRanks = append(pairRanks, rank)
		}
	}
	if len(pairRanks) < 2 {
		return false, nil, nil
	}
	// 取最高的两对
	sort.Slice(pairRanks, func(i, j int) bool { return pairRanks[i] > pairRanks[j] })
	topTwo := map[Rank]bool{pairRanks[0]: true, pairRanks[1]: true}
	var scoring, kickers []Card
	for _, c := range sorted {
		if topTwo[c.Rank] {
			scoring = append(scoring, c)
		} else {
			kickers = append(kickers, c)
		}
	}
	return true, scoring, kickers
}
