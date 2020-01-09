package main

import (
	"encoding/json"
	"fmt"

	"github.com/nsmith5/abnf"
)

const grammar = `
this = %d31
`

func main() {
	m := abnf.RuleList([]byte(grammar))
	s, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(s))
}
