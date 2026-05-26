// Package static содержит анализаторы из пакета staticcheck.io
// (honnef.co/go/tools).
//
// Включены все анализаторы класса SA (staticcheck)
//
// Из остальных классов включены:
//   - QF1001
//   - S1000
//   - ST1000
package static

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/quickfix/qf1001"
	"honnef.co/go/tools/simple/s1000"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck/st1000"
)

func Analyzers() []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	checks = append(checks, qf1001.Analyzer)
	checks = append(checks, s1000.Analyzer)
	checks = append(checks, st1000.Analyzer)

	return checks
}
