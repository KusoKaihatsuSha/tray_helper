/*
Package multichecker consists couple of linters/SAS

1. Fill 'config.json' file

it consists 4 sections:
  - staticcheck - external public checks. array type '[]' with analysis type(use * or SA* for check group analysis)
    see list: https://staticcheck.io/docs/checks/
    also see the info inside main_test.go
  - default - default vet checks. array type '[]' with analysis type(use * or SA* for check group analysis)
    see list: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes
    also see the info inside main_test.go
  - critic - https://github.com/go-critic/go-critic
    see list: https://go-critic.com/overview
  - bodyclose - https://github.com/timakin/bodyclose
    check only exist Body.Close()
  - custom - internal/linters/custom
    check os.Exit() in main.go/main

example:

		//  {
		//  "staticcheck": [
		//      "SA*",
		//      "QF1012"
		//  ],
		//  "default": [
		//      "*"
		//  ],
		//  "critic": true,
		//  "bodyclose": true,
		//  "custom": true
		//  }

	 2. Build file
	    // go build -o ./cmd/staticlint/staticlint.exe ./cmd/staticlint/main.go

	 3. Run
	    // go vet -vettool=cmd/staticlint/staticlint.exe ./...
*/
package main
