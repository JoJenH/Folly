package views

import (
	"math/rand/v2"
	"strings"
	"testing"

	"balatro-cli/engine"
	_ "balatro-cli/engine/jokers"
)

func TestRenderShopShowsItems(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 0))
	items := engine.GenerateShop(rng)
	if len(items) == 0 {
		t.Fatal("expected shop items")
	}
	out := RenderShop(items, 10, 0)
	for _, item := range items {
		if !strings.Contains(out, item.Def.Name) {
			t.Errorf("item name %q not found in shop output", item.Def.Name)
		}
		if !strings.Contains(out, item.Def.Description) {
			t.Errorf("item description %q not found in shop output", item.Def.Description)
		}
	}
}

func TestRenderShopShowsPrice(t *testing.T) {
	items := []engine.ShopItem{
		{Def: &engine.JokerDef{ID: "test", Name: "测试小丑", Description: "测试", Cost: 5}, Price: 5},
	}
	out := RenderShop(items, 10, 0)
	if !strings.Contains(out, "5") {
		t.Error("price not found in shop output")
	}
}

func TestRenderShopInsufficientGold(t *testing.T) {
	items := []engine.ShopItem{
		{Def: &engine.JokerDef{ID: "a", Name: "贵小丑", Description: "很贵", Cost: 10}, Price: 10},
	}
	// gold=3, price=10 → insufficient, should have visual distinction (ANSI)
	outPoor := RenderShop(items, 3, 0)
	outRich := RenderShop(items, 20, 0)
	if outPoor == outRich {
		t.Error("insufficient gold should produce different rendering")
	}
}
