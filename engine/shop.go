package engine

import "math/rand/v2"

// ShopItem 商店商品
type ShopItem struct {
	Def   *JokerDef
	Price int
}

// GenerateShop 随机生成 2 个 Joker 商品
func GenerateShop(rng *rand.Rand) []ShopItem {
	all := DefaultRegistry.All()
	if len(all) == 0 {
		return []ShopItem{}
	}

	items := make([]ShopItem, 0, 2)
	used := make(map[string]bool)

	for len(items) < 2 {
		idx := int(rng.Int64N(int64(len(all))))
		def := all[idx]
		if used[def.ID] {
			continue
		}
		used[def.ID] = true
		price := def.Cost
		if price <= 0 {
			price = 4 // 默认价格
		}
		items = append(items, ShopItem{Def: def, Price: price})
		if len(all) <= len(items) {
			break // 注册的 Joker 不足 2 个时提前退出
		}
	}

	return items
}
