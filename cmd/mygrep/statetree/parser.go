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
			st.State = State{Type: StateTypeDigit}
		case 'w', 'W':
			st.State = State{Type: StateTypeAlpha}
		default:
			return nil, fmt.Errorf("escape sequence at position %d", i)
		}
	default:
		st.State = State{Type: StateTypeChar, Char: re[i]}
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
