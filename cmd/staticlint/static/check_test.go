package static

import (
	"testing"
)

func TestAnalyzers(t *testing.T) {
	analyzers := Analyzers()
	if len(analyzers) == 0 {
		t.Fatal("Analyzers() returned empty list")
	}

	names := make(map[string]bool, len(analyzers))
	for _, a := range analyzers {
		names[a.Name] = true
	}

	required := []string{"SA1000", "SA4006", "SA9004", "QF1001", "S1000", "ST1000"}
	for _, name := range required {
		if !names[name] {
			t.Errorf("missing required analyzer: %s", name)
		}
	}
}
