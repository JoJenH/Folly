package views

import (
	"strings"
	"testing"
)

func TestRenderHelpContainsKeys(t *testing.T) {
	out := RenderHelp()
	keys := []string{"←", "→", "Space", "p", "d", "?", "Ctrl+C"}
	for _, key := range keys {
		if !strings.Contains(out, key) {
			t.Errorf("help screen missing key: %q", key)
		}
	}
}
