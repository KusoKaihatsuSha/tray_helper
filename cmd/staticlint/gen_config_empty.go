//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
)

// generate example 'settings.data'
func main() {
	type MenuItem struct {
		EmulateAction string `json:"actions"`
		Timer         string `json:"timer"`
		Repeat        int    `json:"repeat"`
		Silent        bool   `json:"silent"`
	}

	tmpEl01 := new(MenuItem)
	tmpEl01.Repeat = 1
	tmpEl01.Timer = "30s"
	tmpEl01.Silent = false
	tmpEl01.EmulateAction = "GEN@-16"

	tmpEl02 := new(MenuItem)
	tmpEl02.Repeat = 1
	tmpEl02.Timer = ""
	tmpEl01.Silent = false
	tmpEl02.EmulateAction = "EXEC@notepad|EXECSTD@ping google.com|SLEEP@300ms|TARGET@Notepad|SLEEP@300ms|PRESS@CTRL+V|SLEEP@300ms"

	tmp := make(map[string]*MenuItem, 2)
	tmp["Generate 16 len PASS"] = tmpEl01
	tmp["Ping google.com and paste into open notepad"] = tmpEl02

	output, err := json.MarshalIndent(tmp, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}
