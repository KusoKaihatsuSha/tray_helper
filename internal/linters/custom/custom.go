package custom

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var NoMainExit = &analysis.Analyzer{
	Name: "notosexit",
	Doc:  "check exist os.Exit function on main(). Incorrect call in main()",
	Run:  run001,
}

func run001(pass *analysis.Pass) (interface{}, error) {
	exprExit := func(x *ast.ExprStmt) {
		call, ok := x.X.(*ast.CallExpr)
		if !ok {
			return
		}
		f, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		parent, ok := f.X.(*ast.Ident)
		if !ok {
			return
		}
		if strings.ToLower(parent.Name) != "os" {
			return
		}
		if function := call.Fun.(*ast.SelectorExpr).Sel; strings.ToLower(function.Name) == "exit" {
			pass.Reportf(function.NamePos, "function 'Exit' not permit in 'main' file and 'main' function")
		}
	}
	check := func(s string, vals []string) bool {
		for _, val := range vals {
			if strings.EqualFold(val, s) {
				return true
			}
		}
		return false
	}
	checks := []string{"main", "pkg1"}
	for _, file := range pass.Files {
		findParentFunction := false
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.ExprStmt:
				if findParentFunction {
					exprExit(x)
				}
			case *ast.FuncDecl:
				if check(file.Name.Name, checks) && check(x.Name.Name, checks) {
					findParentFunction = true
				}
			}
			return true
		})
	}
	return nil, nil
}
