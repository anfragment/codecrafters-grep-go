package regexp

import (
	"fmt"
)

// Parse and create a state tree from a regular expression
func NewRegexp(pattern string) (Regexp, error) {
	p := []rune(pattern)
	st, err := compile(p)
	if err != nil {
		return Regexp{}, err
	}

	return Regexp{Root: st}, nil
}

func compile(pattern []rune) (st *StateTree, err error) {
	st = &StateTree{}
	switch pattern[0] {
	case '^':
		st.State = StateStart{}
	case '$':
		st.State = StateEnd{}
	case '.':
		st.State = StateChar{Type: StateTypeAny}
	case '\\':
		pattern = pattern[1:]
		if len(pattern) == 0 {
			return nil, fmt.Errorf("escape sequence at end of pattern")
		}
		switch pattern[0] {
		case 'd', 'D':
			st.State = StateChar{Type: StateTypeDigit}
		case 'w', 'W':
			st.State = StateChar{Type: StateTypeAlpha}
		case '\\':
			st.State = StateChar{Type: StateTypeChar, Char: pattern[0]}
		default:
			return nil, fmt.Errorf("unknown sequence %c", pattern[0])
		}
	case '[':
		pattern = pattern[1:]
		var state StateGroup
		if pattern[0] == '^' {
			state = StateGroup{Type: StateTypeGroupNegative}
			pattern = pattern[1:]
		} else {
			state = StateGroup{Type: StateTypeGroupPositive}
		}
		for {
			if len(pattern) == 0 {
				return nil, fmt.Errorf("unterminated character class")
			}
			if pattern[0] == ']' {
				break
			}
			state.Chars = append(state.Chars, pattern[0])
			pattern = pattern[1:]
		}
		st.State = state
	case '(':
		pattern = pattern[1:]
		alt := StateAlternation{}
		for i := 0; ; i++ {
			if i == len(pattern) {
				return nil, fmt.Errorf("unterminated alternation")
			}
			if pattern[i] == ')' {
				child, err := compile(pattern[:i])
				if err != nil {
					return nil, err
				}
				alt.Children = append(alt.Children, child)
				pattern = pattern[i:]
				break
			}
			if pattern[i] == '|' {
				child, err := compile(pattern[:i])
				if err != nil {
					return nil, err
				}
				alt.Children = append(alt.Children, child)
				pattern = pattern[i+1:]
				i = 0
			}
		}
		st.State = alt
	default:
		st.State = StateChar{Type: StateTypeChar, Char: pattern[0]}
	}

	if len(pattern) > 1 {
		switch pattern[1] {
		case '*', '?':
			st.Quantifier = StateTreeZeroOrMore
			pattern = pattern[1:]
		case '+':
			st.Quantifier = StateTreeOneOrMore
			pattern = pattern[1:]
		default:
			st.Quantifier = StateTreeOne
		}
	}

	if len(pattern) > 1 {
		child, err := compile(pattern[1:])
		if err != nil {
			return nil, err
		}
		st.Child = child
	}

	return st, nil
}
