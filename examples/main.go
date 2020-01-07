package main

import (
	"encoding/json"
	"fmt"

	"github.com/nsmith5/abnf"
)

func main() {
	m := abnf.RuleList([]byte("this = %d31\r\n"))
	s, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(s))
}
