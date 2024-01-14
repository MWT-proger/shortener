package exitcheckanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for direct os.Exit calls in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		for _, decl := range file.Decls {

			fn, ok := decl.(*ast.FuncDecl)

			if !ok {
				continue
			}

			if fn.Name.Name == "main" {
				ast.Inspect(fn.Body, func(node ast.Node) bool {

					switch x := node.(type) {

					case *ast.CallExpr:
						if se, ok := x.Fun.(*ast.SelectorExpr); ok {
							if se.Sel.Name == "Exit" {
								pass.Reportf(x.Fun.(*ast.SelectorExpr).Pos(), "direct call to os.Exit in main")
							}
						}

					}

					return true
				})
			}
		}
	}

	return nil, nil
}
