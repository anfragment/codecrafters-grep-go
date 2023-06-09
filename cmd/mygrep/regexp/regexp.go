package regexp

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

type StateAlternation struct {
	Children []*StateTree
}

type State interface {
	Match(*[]rune, int) (bool, int)
}

const (
	StateTreeOne = iota
	StateTreeOneOrMore
	StateTreeZeroOrMore
)

type StateTree struct {
	State      State
	Quantifier int
	Child      *StateTree
}

type Regexp struct {
	Root *StateTree
}
