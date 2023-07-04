//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
)

// generate linter config file
func main() {
	type Config struct {
		Staticcheck []string `json:"staticcheck"`
		Default     []string `json:"default"`
		Critic      bool     `json:"critic"`
		Bodyclose   bool     `json:"bodyclose"`
		Custom      bool     `json:"custom"`
	}

	tmp := new(Config)
	tmp.Staticcheck = append(tmp.Staticcheck,
		"SA*",
		"S1006",
		"S1008",
		"S1010",
		"S1017",
		"S1018",
		"S1025",
		"S1028",
		"S1030",
		"S1036",
		"S1038",
		"ST1005",
		"ST1008",
		"ST1012",
		"ST1015",
		"ST1017",
		"ST1018",
		"ST1020",
		"ST1021",
		"ST1022",
		"QF1004",
		"QF1006",
		"QF1009",
		"QF1012",
	)
	tmp.Default = append(tmp.Default,
		"*",
	)

	tmp.Critic = true
	tmp.Bodyclose = true
	tmp.Custom = true
	output, err := json.MarshalIndent(tmp, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}
