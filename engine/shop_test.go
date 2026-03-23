package engine

import (
	"math/rand/v2"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 为 engine 包测试注册测试用 Joker（engine 不能导入 engine/jokers，避免循环）
	for _, id := range []string{"test-joker-a", "test-joker-b", "test-joker-c"} {
		id := id
		DefaultRegistry.Register(&JokerDef{
			ID:   id,
			Name: id,
			Cost: 4,
		})
	}
	os.Exit(m.Run())
}

func TestGenerateShopReturnsTwo(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 0))
	items := GenerateShop(rng)
	if len(items) != 2 {
		t.Errorf("GenerateShop returned %d items, want 2", len(items))
	}
}

func TestGenerateShopSameSeed(t *testing.T) {
	rng1 := rand.New(rand.NewPCG(99, 0))
	items1 := GenerateShop(rng1)
	rng2 := rand.New(rand.NewPCG(99, 0))
	items2 := GenerateShop(rng2)
	for i := range items1 {
		if items1[i].Def.ID != items2[i].Def.ID {
			t.Errorf("item[%d]: %s vs %s", i, items1[i].Def.ID, items2[i].Def.ID)
		}
	}
}

func TestGenerateShopPricePositive(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 0))
	items := GenerateShop(rng)
	for i, item := range items {
		if item.Price <= 0 {
			t.Errorf("item[%d] price = %d, want > 0", i, item.Price)
		}
	}
}
