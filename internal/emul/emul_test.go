package emul_test

import (
	"fmt"

	"github.com/KusoKaihatsuSha/tray_helper/internal/emul"
)

func Example_parseCommands() {
	tmp := emul.ParseCommands("TARGET@Notepad|GEN@-16|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms|PRESS@ENTER|PRINT@||SLEEP@300ms|EXECSTD@ping google.com|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms|CLIP@test")

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
}
