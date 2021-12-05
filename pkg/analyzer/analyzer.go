package analyzer

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "mylinter",
		Doc:   "finds bad words about teamlead",
		Run:   run,
		Flags: flagSet,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		// parse file
		ast.Inspect(f, func(node ast.Node) bool {
			f, ok := node.(*ast.Comment)
			if !ok {
				return true
			}
			pass.Reportf(node.Pos(), "calculated cyclomatic complexity for function %s is %d, max is %d", f.Name.Name, comp, maxComplexity)

		})
	}

	return nil, nil
}
