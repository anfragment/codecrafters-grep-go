package statetree

import (
	"unicode"
)

func (st StateChar) Match(line *[]rune, i int) (match bool, skip int) {
	skip = 1
	r := (*line)[i]
	switch st.Type {
	case StateTypeAny:
		return true, skip
	case StateTypeChar:
		return r == st.Char, skip
	case StateTypeAlpha:
		return unicode.IsLetter(r), skip
	case StateTypeDigit:
		return unicode.IsDigit(r), skip
	default:
		return false, skip
	}
}

func (st StateGroup) Match(line *[]rune, i int) (match bool, skip int) {
	skip = 1
	r := (*line)[i]
	switch st.Type {
	case StateTypeGroupPositive:
		for _, c := range st.Chars {
			if r == c {
				return true, skip
			}
		}
		return false, 1
	case StateTypeGroupNegative:
		for _, c := range st.Chars {
			if r == c {
				return false, skip
			}
		}
		return true, skip
	default:
		return false, skip
	}
}

func (st StateTree) Match(line *[]rune) bool {
	for i := 0; i < len(*line); i++ {
		for j := i; j < len(*line); {
			st := st
			match, skip := st.State.Match(line, j)
			if !match {
				break
			}
			j += skip
			if len(st.Children) == 0 {
				return true
			} else {
				st = *st.Children[0]
			}
		}
	}
	return false
}
