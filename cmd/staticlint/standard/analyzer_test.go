package standard

import (
	"testing"
)

func TestAnalyzers(t *testing.T) {
	if len(Analyzers) == 0 {
		t.Fatal("Analyzers is empty")
	}

	names := make(map[string]bool, len(Analyzers))
	for _, a := range Analyzers {
		names[a.Name] = true
	}

	required := []string{"printf", "shadow", "structtag", "copylocks", "unmarshal", "loopclosure"}
	for _, name := range required {
		if !names[name] {
			t.Errorf("missing required analyzer: %s", name)
		}
	}
}
