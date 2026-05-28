// Package exitcheck provides a custom static analyzer
// that forbids direct calls to os.Exit inside the main
// function of package main.
package exitcheck

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports direct os.Exit calls
// inside the main function of package main.
var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for direct os.Exit calls inside main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		if strings.Contains(filename, "/go-build/") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if fn.Name.Name != "main" {
				return false
			}

			if fn.Body == nil {
				return false
			}

			ast.Inspect(fn.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}

				selector, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkgIdent, ok := selector.X.(*ast.Ident)
				if !ok {
					return true
				}

				if pkgIdent.Name == "os" && selector.Sel.Name == "Exit" {
					pass.Reportf(
						call.Pos(),
						"direct call to os.Exit inside main function is forbidden",
					)
				}

				return true
			})

			return false
		})
	}

	return nil, nil
}
