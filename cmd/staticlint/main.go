// Binary linker.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	critic "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	qf "honnef.co/go/tools/quickfix"
	s "honnef.co/go/tools/simple"
	sa "honnef.co/go/tools/staticcheck"
	st "honnef.co/go/tools/stylecheck"

	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"

	"github.com/KusoKaihatsuSha/tray_helper/internal/linters/custom"
)

const Config = `config.json`

// ConfigData consist data from file config
type ConfigData struct {
	Staticcheck []string `json:"staticcheck"`
	Default     []string `json:"default"`
	Critic      bool     `json:"critic"`
	Bodyclose   bool     `json:"bodyclose"`
	Custom      bool     `json:"custom"`
}

func main() {
	defaultAnalysis := map[string][]*analysis.Analyzer{
		asmdecl.Analyzer.Name:             {asmdecl.Analyzer},
		assign.Analyzer.Name:              {assign.Analyzer},
		atomic.Analyzer.Name:              {atomic.Analyzer},
		atomicalign.Analyzer.Name:         {atomicalign.Analyzer},
		bools.Analyzer.Name:               {bools.Analyzer},
		buildssa.Analyzer.Name:            {buildssa.Analyzer},
		buildtag.Analyzer.Name:            {buildtag.Analyzer},
		cgocall.Analyzer.Name:             {cgocall.Analyzer},
		composite.Analyzer.Name:           {composite.Analyzer},
		copylock.Analyzer.Name:            {copylock.Analyzer},
		ctrlflow.Analyzer.Name:            {ctrlflow.Analyzer},
		deepequalerrors.Analyzer.Name:     {deepequalerrors.Analyzer},
		errorsas.Analyzer.Name:            {errorsas.Analyzer},
		fieldalignment.Analyzer.Name:      {fieldalignment.Analyzer},
		findcall.Analyzer.Name:            {findcall.Analyzer},
		ifaceassert.Analyzer.Name:         {ifaceassert.Analyzer},
		lostcancel.Analyzer.Name:          {lostcancel.Analyzer},
		nilness.Analyzer.Name:             {nilness.Analyzer},
		shadow.Analyzer.Name:              {shadow.Analyzer},
		stringintconv.Analyzer.Name:       {stringintconv.Analyzer},
		unmarshal.Analyzer.Name:           {unmarshal.Analyzer},
		framepointer.Analyzer.Name:        {framepointer.Analyzer},
		httpresponse.Analyzer.Name:        {httpresponse.Analyzer},
		inspect.Analyzer.Name:             {inspect.Analyzer},
		loopclosure.Analyzer.Name:         {loopclosure.Analyzer},
		nilfunc.Analyzer.Name:             {nilfunc.Analyzer},
		pkgfact.Analyzer.Name:             {pkgfact.Analyzer},
		printf.Analyzer.Name:              {printf.Analyzer},
		reflectvaluecompare.Analyzer.Name: {reflectvaluecompare.Analyzer},
		shift.Analyzer.Name:               {shift.Analyzer},
		sigchanyzer.Analyzer.Name:         {sigchanyzer.Analyzer},
		sortslice.Analyzer.Name:           {sortslice.Analyzer},
		stdmethods.Analyzer.Name:          {stdmethods.Analyzer},
		structtag.Analyzer.Name:           {structtag.Analyzer},
		testinggoroutine.Analyzer.Name:    {testinggoroutine.Analyzer},
		tests.Analyzer.Name:               {tests.Analyzer},
		unreachable.Analyzer.Name:         {unreachable.Analyzer},
		unsafeptr.Analyzer.Name:           {unsafeptr.Analyzer},
		unusedresult.Analyzer.Name:        {unusedresult.Analyzer},
		unusedwrite.Analyzer.Name:         {unusedwrite.Analyzer},
		usesgenerics.Analyzer.Name:        {usesgenerics.Analyzer},
	}

	// add analysis from std vet
	var tmp []*analysis.Analyzer
	for _, v := range defaultAnalysis {
		tmp = append(tmp, v...)
	}
	defaultAnalysis["*"] = tmp

	allStaticcheck := helpers.Concatenate(sa.Analyzers, s.Analyzers, st.Analyzers, qf.Analyzers)
	defaultAnalysisStaticcheck := map[string][]*analysis.Analyzer{}
	for _, v := range allStaticcheck {
		defaultAnalysisStaticcheck[v.Analyzer.Name] = []*analysis.Analyzer{v.Analyzer}
		defaultAnalysisStaticcheck["*"] = append(defaultAnalysisStaticcheck["*"], v.Analyzer)
		prefix, _ := helpers.SplitPrefix(v.Analyzer.Name)
		defaultAnalysisStaticcheck[prefix+"*"] = append(defaultAnalysisStaticcheck[prefix+"*"], v.Analyzer)
	}

	appfile, err := os.Executable()

	if err != nil {
		panic(err)
	}

	filebody, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	err = json.Unmarshal(filebody, &cfg)
	if err != nil {
		panic(err)
	}

	mychecks := []*analysis.Analyzer{}

	for _, v := range cfg.Default {
		mychecks = append(mychecks, defaultAnalysis[v]...)
	}

	for _, v := range cfg.Staticcheck {
		mychecks = append(mychecks, defaultAnalysisStaticcheck[v]...)
	}

	if cfg.Critic {
		// Add all flags and then disable redundant
		critic.Analyzer.Flags.Set("enable-all", "true")
		// -disable=appendAssign
		// -disable=#diagnostic,#opinionated,#security,#style,#performance,#experimental"
		critic.Analyzer.Flags.Set("disable", "#experimental,#performance,#opinionated")
		mychecks = append(mychecks, critic.Analyzer)
	}
	if cfg.Bodyclose {
		mychecks = append(mychecks, bodyclose.Analyzer)
	}
	if cfg.Custom {
		mychecks = append(mychecks, custom.NoMainExit)
	}
	multichecker.Main(
		mychecks...,
	)
}
