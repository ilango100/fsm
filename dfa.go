package fsm

// Definition of DFA lives here. It allows DFA definition in the form of
//  - Transition Table
// 	- Transition Function
//	- Transition Graph

// DFA defines Deterministic Finite Automata
type DFA interface {
	CurrentState() string
	StartState() string
	AcceptingStates() []string
	Reset() string
	SetState(string) error
	Next(rune) (string, error)
}
