package emul_test

import (
	"fmt"

	"github.com/KusoKaihatsuSha/tray_helper/internal/emul"
)

func Example_parseCommands() {
	tmp := emul.ParseCommands("TARGET@Notepad|GEN@-16|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms|PRESS@ENTER|PRINT@||SLEEP@300ms|EXECSTD@ping google.com|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms|CLIP@test|cmd \\k C:\\Program Files\\totalcmd\\TOTALCMD64.EXE")

	for _, val := range tmp {
		fmt.Printf("| action: %s, data: %s |\n", val.Name, val.Data)
	}

	// Output:
	// | action: TARGET, data: Notepad |
	// | action: GEN, data: -16 |
	// | action: SLEEP, data: 300ms |
	// | action: PRESS, data: CTRL+V |
	// | action: SLEEP, data: 300ms |
	// | action: PRESS, data: ENTER |
	// | action: PRINT, data:  |
	// | action: SLEEP, data: 300ms |
	// | action: EXECSTD, data: ping google.com |
	// | action: SLEEP, data: 300ms |
	// | action: PRESS, data: CTRL+V |
	// | action: SLEEP, data: 300ms |
	// | action: CLIP, data: test |
	// | action: cmd \k C:\Program Files\totalcmd\TOTALCMD64.EXE, data:  |
}

func Example_parse() {

	text01 := `c:\Users\user\Downloads\tempApp\some.exe`
	text02 := `"c:\Users\user\Downloads\temp App\some.exe"`
	text03 := `echo -e "Hello world\n\tfrom golang"`
	text04 := `ping -n 10 ya.ru`
	text05 := `cmd.exe /k "c:\Users\user\Downloads\temp App\some.exe"`
	text06 := `ssh -i ~/.ssh/dops test@127.0.0.1`

	for _, val := range emul.ParseCommandCmd(text01) {
		fmt.Printf("%s\n", val)
	}
	for _, val := range emul.ParseCommandCmd(text02) {
		fmt.Printf("%s\n", val)
	}
	for _, val := range emul.ParseCommandCmd(text03) {
		fmt.Printf("%s\n", val)
	}
	for _, val := range emul.ParseCommandCmd(text04) {
		fmt.Printf("%s\n", val)
	}
	for _, val := range emul.ParseCommandCmd(text05) {
		fmt.Printf("%s\n", val)
	}
	for _, val := range emul.ParseCommandCmd(text06) {
		fmt.Printf("%s\n", val)
	}

	// Output:
	// c:\Users\user\Downloads\tempApp\some.exe
	// c:\Users\user\Downloads\temp App\some.exe
	// echo
	// -e
	// Hello world\n\tfrom golang
	// ping
	// -n
	// 10
	// ya.ru
	// cmd.exe
	// /k
	// c:\Users\user\Downloads\temp App\some.exe
	// ssh
	// -i
	// ~/.ssh/dops
	// test@127.0.0.1
}
