package statetree

// StateChar types
const (
	StateTypeAny = iota
	StateTypeChar
	StateTypeAlpha
	StateTypeDigit
)

type StateChar struct {
	Type int
	Char rune
}

// StateGroup types
const (
	StateTypeGroupPositive = iota
	StateTypeGroupNegative
)

type StateGroup struct {
	Type  int
	Chars []rune
}

type StateStart struct{}

type StateEnd struct{}

type State interface {
	Match(*[]rune, int) (bool, int)
}

type StateTree struct {
	State    State
	Children []*StateTree
}
