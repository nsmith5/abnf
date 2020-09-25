package abnf

import (
	"fmt"

	"github.com/nsmith5/abnf/rules"
)

func Parse(r *rules.Rule, input []byte) {
	var stack []*rules.Rule
	stack = append(stack, r)

	cursor := 0 // Position in input

	for len(stack) != 0 {
		// Pop top of stack
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		children := (*current).Children()
		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, children[i])
		}

		if len(children) == 0 {
			// This node is terminal
			end := len(input) - 1
			if cursor > end {
				break
			}
			cursor++
			fmt.Println("Terminal")
		}
		fmt.Printf("Def: %s\n", (*current).Definition())
	}
}
