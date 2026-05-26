// Package exitcheck реализует анализатор, запрещающий прямой вызов os.Exit
// в функции main пакета main.
package exitcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer запрещает прямой вызов os.Exit в функции main пакета main.
// Рекомендуемая альтернатива: log.Fatal, log.Fatalf или возврат ошибки.
var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "prohibit direct os.Exit call in the main function of package main; use log.Fatal or return an error instead",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		if isGenerated(file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch call := node.(type) {
			case *ast.FuncDecl:
				return call.Name.Name == "main"
			case *ast.CallExpr:
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkg, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}
				if pkg.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(), "direct call to os.Exit in main function is not allowed; use log.Fatal or return an error instead")
				}
			}

			return true
		})
	}

	return nil, nil
}

// isGenerated проверяет, является ли файл автоматически сгенерированным.
func isGenerated(f *ast.File) bool {
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, "Code generated") && strings.Contains(c.Text, "DO NOT EDIT") {
				return true
			}
		}
	}

	return false
}
