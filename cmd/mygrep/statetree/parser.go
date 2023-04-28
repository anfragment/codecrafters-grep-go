package statetree

import (
	"fmt"
)

// Parse and create a state tree from a regular expression
func NewStateTree(re string) (*StateTree, error) {
	reRune := []rune(re)
	st, err := parse(reRune, 0)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func parse(re []rune, i int) (st *StateTree, err error) {
	if i >= len(re) {
		return nil, fmt.Errorf("end of pattern")
	}
	st = &StateTree{}
	switch re[i] {
	case '\\':
		i++
		if i >= len(re) {
			return nil, fmt.Errorf("escape sequence at position %d", i)
		}
		switch re[i] {
		case 'd', 'D':
			st.State = StateChar{Type: StateTypeDigit}
		case 'w', 'W':
			st.State = StateChar{Type: StateTypeAlpha}
		default:
			return nil, fmt.Errorf("unknown sequence at position %d", i)
		}
	case '[':
		i++
		var state StateGroup
		if re[i] == '^' {
			state = StateGroup{Type: StateTypeGroupNegative}
			i++
		} else {
			state = StateGroup{Type: StateTypeGroupPositive}
		}
		for {
			if i >= len(re) {
				return nil, fmt.Errorf("cannot find closing bracket at position %d", i)
			}
			if re[i] == ']' {
				break
			}
			state.Chars = append(state.Chars, re[i])
			i++
		}
		st.State = state
	default:
		st.State = StateChar{Type: StateTypeChar, Char: re[i]}
	}

	if i+1 < len(re) {
		child, err := parse(re, i+1)
		if err != nil {
			return nil, err
		}
		st.Children = append(st.Children, child)
	}

	return st, nil
}
