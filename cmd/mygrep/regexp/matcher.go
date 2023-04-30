package regexp

import (
	"unicode"
)

func (st StateChar) Match(line *[]rune, i int) (good bool, skip int) {
	if i >= len(*line) {
		return false, 1
	}
	r := (*line)[i]
	switch st.Type {
	case StateTypeAny:
		return true, 1
	case StateTypeChar:
		return r == st.Char, 1
	case StateTypeAlpha:
		return unicode.IsLetter(r), 1
	case StateTypeDigit:
		return unicode.IsDigit(r), 1
	default:
		return false, 1
	}
}

func (st StateGroup) Match(line *[]rune, i int) (good bool, skip int) {
	if i >= len(*line) {
		return false, 1
	}
	r := (*line)[i]
	switch st.Type {
	case StateTypeGroupPositive:
		for _, c := range st.Chars {
			if r == c {
				return true, 1
			}
		}
		return false, 1
	case StateTypeGroupNegative:
		for _, c := range st.Chars {
			if r == c {
				return false, 1
			}
		}
		return true, 1
	default:
		return false, 1
	}
}

func (st StateStart) Match(line *[]rune, i int) (bool, int) {
	return i == 0, 0
}

func (st StateEnd) Match(line *[]rune, i int) (bool, int) {
	return i == len(*line), 0
}

func (st StateAlternation) Match(line *[]rune, i int) (bool, int) {
	for _, child := range st.Children {
		good, skip := child.Match(line, i)
		if good {
			return true, skip
		}
	}
	return false, -1
}

func (st StateTree) Match(line *[]rune, i int) (bool, int) {
	for ; i <= len(*line); i++ {
		good := match(&st, line, i)
		if good {
			return true, i
		}
	}
	return false, -1
}

func match(st *StateTree, line *[]rune, i int) bool {
	good, skip := st.State.Match(line, i)
	if !good && st.Quantifier != StateTreeZeroOrMore {
		return false
	}
	if st.Child == nil {
		return true
	}

	switch st.Quantifier {
	case StateTreeOne:
		return match(st.Child, line, i+skip)
	case StateTreeOneOrMore:
		st := st
		st.Quantifier = StateTreeZeroOrMore
		return match(st.Child, line, i+skip) || match(st, line, i+skip)
	case StateTreeZeroOrMore:
		if !good {
			return match(st.Child, line, i)
		}
		return match(st.Child, line, i+skip) || match(st, line, i+skip)
	default:
		return false
	}
}

func (r Regexp) Match(line *[]rune) bool {
	match, _ := r.Root.Match(line, 0)
	return match
}
