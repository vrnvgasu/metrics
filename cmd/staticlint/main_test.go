package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/vrnvgasu/metrics/cmd/staticlint/exitcheck"
)

func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), exitcheck.Analyzer, "./...")
}
