package abnf

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

type Match struct {
	Rule   string
	Result []byte
	Child  []Match
}

type Rule func(input []byte) *Match

func Byte(name string, target byte) Rule {
	return func(input []byte) *Match {
		if input[0] == target {
			return &Match{
				Rule:   name,
				Result: input[0:1],
				Child:  nil,
			}
		}
		return nil
	}
}

func ByteRange(name string, start, stop byte) Rule {
	return func(input []byte) *Match {
		if input[0] >= start && input[0] <= stop {
			return &Match{
				Rule:   name,
				Result: input[0:1],
				Child:  nil,
			}
		}
		return nil
	}
}

func String(name, target string, caseSensitive bool) Rule {
	if caseSensitive {
		return func(input []byte) *Match {
			if len(target) > len(input) {
				return nil
			}
			if bytes.Equal([]byte(target), input[:len(target)]) {
				return &Match{
					Rule:   name,
					Result: input[:len(target)],
					Child:  nil,
				}
			}
			return nil
		}
	}
	return func(input []byte) *Match {
		if len(target) < len(input) {
			return nil
		}
		if strings.ToUpper(target) == strings.ToUpper(string(input[:len(target)])) {
			return &Match{
				Rule:   name,
				Result: input[:len(target)],
				Child:  nil,
			}
		}
		return nil
	}
}

func Sequence(name string, rules ...Rule) Rule {
	return func(input []byte) *Match {
		matches := make([]Match, len(rules))
		temp := input
		matched := 0
		for i, rule := range rules {
			match := rule(temp)
			if match == nil {
				return nil
			}
			matches[i] = *match
			temp = temp[len(match.Result):]
			matched += len(match.Result)
		}
		return &Match{
			Rule:   name,
			Result: input[:matched],
			Child:  matches,
		}
	}
}

func Or(name string, rules ...Rule) Rule {
	return func(input []byte) *Match {
		for _, rule := range rules {
			match := rule(input)
			if match != nil {
				return &Match{
					Rule:   name,
					Result: match.Result,
					Child:  []Match{*match},
				}
			}
		}
		return nil
	}
}

func Repeat(name string, min, max int, rule Rule) Rule {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		// TODO: Probably a better way to deal with 'inifinity' here. Is there
		// some sane way to cap based on memory or something?
		max = math.MaxInt64
	}
	return func(input []byte) *Match {
		var matches []Match
		temp := input
		matched := 0
		for {
			match := rule(temp)
			if match == nil {
				break
			}

			matches = append(matches, *match)
			temp = temp[len(match.Result):]
			matched += len(match.Result)
		}

		if len(matches) < min || len(matches) > max {
			return nil
		}

		return &Match{
			Rule:   name,
			Result: input[:matched],
			Child:  matches,
		}
	}
}

func Option(name string, rule Rule) Rule {
	return func(input []byte) *Match {
		match := rule(input)
		if match == nil {
			return &Match{
				Rule:   name,
				Result: nil,
				Child:  nil,
			}
		}
		return &Match{
			Rule:   name,
			Result: match.Result,
			Child:  []Match{*match},
		}
	}
}
