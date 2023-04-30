package statetree

import (
	"fmt"
)

// Parse and create a state tree from a regular expression
func NewStateTree(re string) (*StateTree, error) {
	reRune := []rune(re)
	st, err := compile(reRune, 0)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func compile(re []rune, i int) (st *StateTree, err error) {
	if i >= len(re) {
		return nil, fmt.Errorf("end of pattern")
	}
	st = &StateTree{}
	switch re[i] {
	case '^':
		st.State = StateStart{}
	case '$':
		st.State = StateEnd{}
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
		case '\\':
			st.State = StateChar{Type: StateTypeChar, Char: re[i]}
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
		switch re[i+1] {
		case '*':
			st.Quantifier = StateTreeZeroOrMore
			i++
		case '+':
			st.Quantifier = StateTreeOneOrMore
			i++
		default:
			st.Quantifier = StateTreeOne
		}
	}

	if i+1 < len(re) {
		child, err := compile(re, i+1)
		if err != nil {
			return nil, err
		}
		st.Children = append(st.Children, child)
	}

	return st, nil
}
