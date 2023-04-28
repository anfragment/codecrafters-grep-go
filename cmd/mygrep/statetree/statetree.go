package statetree

// state types
const (
	StateTypeAny = iota
	StateTypeChar
	StateTypeAlpha
	StateTypeDigit
)

type State struct {
	Type int
	Char rune
}

type StateTree struct {
	State    State
	Children []*StateTree
}
