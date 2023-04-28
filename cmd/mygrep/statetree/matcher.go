package statetree

import (
	"unicode"
)

func (st State) Match(r rune) bool {
	switch st.Type {
	case StateTypeAny:
		return true
	case StateTypeChar:
		return r == st.Char
	case StateTypeAlpha:
		return unicode.IsLetter(r)
	case StateTypeDigit:
		return unicode.IsDigit(r)
	default:
		return false
	}
}

func (st StateTree) Match(line []byte) bool {
	for i := 0; i < len(line); i++ {
		if st.State.Match(rune(line[i])) {
			if len(st.Children) == 0 {
				return true
			} else {
				// TODO: multiple children
				st = *st.Children[0]
			}
		}
	}
	return false
}
