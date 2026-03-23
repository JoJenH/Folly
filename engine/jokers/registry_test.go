package jokers

import (
	"balatro-cli/engine"
	"testing"
)

func TestRegistryAllNonEmpty(t *testing.T) {
	all := engine.DefaultRegistry.All()
	if len(all) == 0 {
		t.Error("DefaultRegistry.All() should return non-empty list after init")
	}
}

func TestRegistryByIDFound(t *testing.T) {
	all := engine.DefaultRegistry.All()
	if len(all) == 0 {
		t.Skip("no jokers registered")
	}
	first := all[0]
	j, ok := engine.DefaultRegistry.ByID(first.ID)
	if !ok {
		t.Errorf("ByID(%q) not found", first.ID)
	}
	if j.ID != first.ID {
		t.Errorf("ByID returned wrong joker: got %q, want %q", j.ID, first.ID)
	}
}

func TestRegistryByIDNotFound(t *testing.T) {
	_, ok := engine.DefaultRegistry.ByID("nonexistent-joker-id")
	if ok {
		t.Error("ByID with unknown ID should return false")
	}
}
