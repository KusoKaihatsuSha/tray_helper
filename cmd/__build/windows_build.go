//go:build ignore
// +build ignore

package windows_build

// Example of use. It must be in the CI, not here. Or at last run 'go generate' in CI.
// input
//go:generate -command runscript01 cmd /k echo FOR /F %%A IN ('git rev-parse HEAD') DO set "GOOCOMMIT=%%~A"
//go:generate -command runscript02 cmd /k echo FOR /F "tokens=1* delims=" %%A IN ('echo %date% [%time:~0,2%-%time:~3,2%-%time:~6,2%]') DO set "GOODATE=%%~A"
//go:generate -command runscript03 cmd /k echo FOR /F %%A IN ('echo 1.0.0') DO set "GOOVERSION=%%~A"
//go:generate -command runscript04ico cmd /k echo go run github.com/akavel/rsrc@latest -ico=../../internal/config/icons/icon.ico -o=../tray_helper/rsrc.syso
//go:generate -command runscript04 cmd /k echo go run github.com/mitchellh/gox@latest -os windows -ldflags "-s -w -H=windowsgui -X main.buildVersion=%GOOVERSION% -X 'main.buildDate=%GOODATE%' -X main.buildCommit=%GOOCOMMIT%" -output bin_windows/{{.Dir}}_{{.OS}}_{{.Arch}} ../...
//go:generate -command rmdir cmd /k rmdir /s /q
//go:generate -command mkdir cmd /k mkdir
//go:generate -command rm cmd /k DEL /S
//go:generate -command run cmd /k
//go:generate -command htmlexport cmd /k go tool cover -html=wcover.out -o ./bin_windows/coverage.html
//go:generate -command runtest cmd /k "go test ../../... -race -cover -coverpkg=../../internal/... -coverprofile ./wcover.out && go tool cover -func ./wcover.out | findstr /c:total >> ./bin_windows/coverage_output.txt"
//go:generate -command runlinter cmd /k "go run ../staticlint/gen_config.go >> ./bin_windows/config.json"
//go:generate -command runexample cmd /k "go run ../staticlint/gen_config_empty.go >> ./bin_windows/settings.data"
//go:generate -command savecover cmd /k "go vet -vettool=bin_windows/staticlint_windows_amd64.exe ../../cmd/... ../../internal/... 2>> ./bin_windows/linter_output.txt"
// run
//go:generate rm runscript.bat
//go:generate rm wcover.out
//go:generate runscript01 >> runscript.bat
//go:generate runscript02 >> runscript.bat
//go:generate runscript03 >> runscript.bat
//go:generate runscript04ico >> runscript.bat
//go:generate runscript04 >> runscript.bat
//go:generate rmdir bin_windows
//go:generate mkdir bin_windows
//go:generate run runscript.bat
//go:generate runtest
//go:generate runlinter
//go:generate runexample
//go:generate savecover
//go:generate htmlexport
//go:generate rm runscript.bat
//go:generate rm wcover.out


go run github.com/akavel/rsrc@latest -ico="icon.ico" -o="rsrc.syso"