package main

import (
	"strings"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/MWT-proger/shortener/internal/staticlint/exitcheckanalyzer"
)

func getStaticcheckAnalyzersSA(checks []*analysis.Analyzer) []*analysis.Analyzer {

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		}
	}
	return checks
}

func getStylecheckAnalyzers(checks []*analysis.Analyzer, data map[string]bool) []*analysis.Analyzer {

	for _, v := range stylecheck.Analyzers {

		if data[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	return checks
}

func main() {
	var checks []*analysis.Analyzer

	checks = getStaticcheckAnalyzersSA(checks)
	checks = getStylecheckAnalyzers(checks, map[string]bool{"ST1001": true})

	checks = append(checks, printf.Analyzer)
	checks = append(checks, shadow.Analyzer)
	checks = append(checks, structtag.Analyzer)

	checks = append(checks, ineffassign.Analyzer)
	checks = append(checks, bodyclose.Analyzer)

	checks = append(checks, exitcheckanalyzer.ExitAnalyzer)

	multichecker.Main(
		checks...,
	)
}
