// Package main implements a custom multichecker for static code analysis.
//
// The multichecker combines:
//
//   - standard Go analyzers from golang.org/x/tools/go/analysis/passes;
//   - all SA-class analyzers from staticcheck;
//   - style analyzers from stylecheck;
//   - third-party analyzers;
//   - a custom analyzer that forbids direct os.Exit calls
//     inside the main function of package main.
//
// Run the multichecker:
//
//	go run ./cmd/staticlint ./...
//
// Or build and run:
//
//	go build -o staticlint ./cmd/staticlint
//	./staticlint ./...
//
// Included standard analyzers:
//
//	assign          detects useless assignments;
//	atomic          checks for common mistakes using sync/atomic;
//	bools           detects common boolean expression mistakes;
//	buildtag        validates build tags;
//	cgocall         detects invalid cgo pointer passing;
//	composite       checks for unkeyed composite literals;
//	copylock        detects locks passed by value;
//	errorsas        validates errors.As usage;
//	httpresponse    checks for mistakes using HTTP responses;
//	loopclosure     detects references to loop variables from closures;
//	lostcancel      checks for missing context cancellation calls;
//	nilfunc         detects useless comparisons with nil functions;
//	printf          validates Printf-style formatting;
//	shadow          detects variable shadowing;
//	shift           checks for invalid bit shifts;
//	stdmethods      checks signatures of well-known interface methods;
//	stringintconv   detects suspicious string-to-int conversions;
//	structtag       validates struct field tags;
//	tests           checks incorrect usages of tests and examples;
//	unmarshal       detects invalid unmarshaling targets;
//	unreachable     detects unreachable code;
//	unsafeptr       checks unsafe.Pointer usage;
//	unusedresult    detects unused function results.
//
// Included staticcheck analyzers:
//
//	SA* analyzers detect correctness issues, bugs,
//	suspicious constructs, and dangerous patterns.
//
// Included stylecheck analyzers:
//
//	ST* analyzers enforce style and naming conventions.
//
// Included third-party analyzers:
//
//	nilerr      detects returning nil after non-nil errors;
//	bodyclose   checks whether HTTP response bodies are closed.
//
// Included custom analyzer:
//
//	exitcheck   forbids direct os.Exit calls inside
//	            the main function of package main.
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	// standard passes
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	// third-party analyzers
	"github.com/gostaticanalysis/nilerr"
	"github.com/timakin/bodyclose/passes/bodyclose"

	// staticcheck
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	// exit check
	"github.com/FeshLig/metrcollector/cmd/staticlint/exitcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	// standard analyzers
	analyzers = append(analyzers,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	// all SA analyzers
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	// at least one non-SA analyzer
	for _, v := range stylecheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	// public analyzers
	analyzers = append(analyzers,
		nilerr.Analyzer,
		bodyclose.Analyzer,
	)

	// exit check
	analyzers = append(analyzers,
		exitcheck.Analyzer,
	)

	multichecker.Main(analyzers...)
}
