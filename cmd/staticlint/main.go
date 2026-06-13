// staticlint — multichecker для статического анализа Go-кода.
//
// # Запуск
//
//	go build -o staticlint ./cmd/staticlint
//	./staticlint ./...
//
// Для проверки конкретного пакета:
//
//	./staticlint ./internal/server/...
//
// Список доступных анализаторов и их описания:
//
//	./staticlint help
//	./staticlint help <name>
//
// # Состав анализаторов
//
// Стандартные анализаторы (golang.org/x/tools/go/analysis/passes):
// анализаторы из официального пакета passes.
//
// Анализаторы staticcheck.io:
//   - SA (staticcheck) — все анализаторы.
//   - QF (quickfix) — QF1001.
//   - S1 (simple) — S1000.
//   - ST1 (stylecheck) — ST1000.
//
// Публичные анализаторы:
//   - bodyclose (github.com/timakin/bodyclose): проверяет, что тело HTTP-ответа
//     закрывается после использования, предотвращая утечку соединений.
//   - nilerr (github.com/gostaticanalysis/nilerr): обнаруживает возврат nil-ошибки
//     вместо реальной ненулевой ошибки.
//
// Собственный анализатор:
//   - exitcheck: запрещает прямой вызов os.Exit в функции main пакета main.
//     Подробнее: ./staticlint help exitcheck
package main

import (
	"github.com/gostaticanalysis/nilerr"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/vrnvgasu/metrics/cmd/staticlint/exitcheck"
	"github.com/vrnvgasu/metrics/cmd/staticlint/standard"
	"github.com/vrnvgasu/metrics/cmd/staticlint/static"
)

func main() {
	checks := make([]*analysis.Analyzer, 0)
	checks = append(checks, standard.Analyzers...)
	checks = append(checks, static.Analyzers()...)

	// Публичные анализаторы
	checks = append(checks, bodyclose.Analyzer)
	checks = append(checks, nilerr.Analyzer)

	// Своя проверка на exit в main
	checks = append(checks, exitcheck.Analyzer)

	multichecker.Main(checks...)
}
