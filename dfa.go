package fsm

import "fmt"

import "strings"

// Definition of DFA lives here. It allows DFA definition in the form of
//  - Transition Table
// 	- Transition Function
//	- Transition Graph

// DFAInfo defines information about DFA
type DFAInfo struct {
	States          []string
	Alphabet        []byte
	InitialState    string
	AcceptingStates []string
}

func (di DFAInfo) String() string {
	str := fmt.Sprintf("\nStates:\t\t\t%v\n", di.States)
	str += "Alphabet:\t\t["
	for _, a := range di.Alphabet {
		str += fmt.Sprintf("%c ", a)
	}
	str = str[:len(str)-1]
	str += fmt.Sprintf("]\nInitialState:\t\t%s\nAcceptingStates:\t%v\n", di.InitialState, di.AcceptingStates)
	str = strings.ReplaceAll(str, "[]string", "")
	str = strings.ReplaceAll(str, "[]byte", "")
	return str
}

// DFA defines Deterministic Finite Automata
type DFA interface {
	Reset() string
	Next(byte) string
	IsAccepted() bool
	Info() DFAInfo

	// CurrentState() string // It is obtained on each call to next
	// StartState() string	// Obtained from a call to Reset()
	// SetState(string) error	// May allow abuse
}

// DFATable creates a DFA from Transition Table
func DFATable(transitionTable map[string]map[byte]string) (DFA, error) {
	// Parse the table to get all alphabets and states
	alphabet := make([]byte, 0)
	states := make([]string, 0)
	initialState := ""
	acceptingStates := make([]string, 0)

	for state := range transitionTable {
		if state[0] == '>' {
			transitionTable[state[1:]] = transitionTable[state]
			delete(transitionTable, state)
			state = state[1:]
			initialState = state
		}
		if state[0] == '*' {
			transitionTable[state[1:]] = transitionTable[state]
			delete(transitionTable, state)
			if initialState == state {
				initialState = state[1:]
			}
			state = state[1:]
			acceptingStates = append(acceptingStates, state)
		}

		for _, s := range states {
			if s == state {
				return nil, fmt.Errorf("Duplicate rows for state: %s", s)
			}
		}

		states = append(states, state)
	}

	for _, row := range transitionTable {
		for r := range row {
			got := false
			for _, rr := range alphabet {
				if r == rr {
					got = true
					break
				}
			}
			if !got {
				alphabet = append(alphabet, r)
			}
		}
	}

	return &tableDFA{
		currentState: initialState,
		transition:   transitionTable,
		info: DFAInfo{
			States:          states,
			Alphabet:        alphabet,
			InitialState:    initialState,
			AcceptingStates: acceptingStates,
		},
	}, nil
}

type tableDFA struct {
	currentState string
	transition   map[string]map[byte]string
	info         DFAInfo
}

func (t *tableDFA) Info() DFAInfo {
	return t.info
}

func (t *tableDFA) Reset() string {
	t.currentState = t.info.InitialState
	return t.currentState
}

func (t *tableDFA) IsAccepted() bool {
	if t.currentState == "NULL" {
		return false
	}

	for i := range t.info.AcceptingStates {
		if t.currentState == t.info.AcceptingStates[i] {
			return true
		}
	}
	return false
}

func (t *tableDFA) Next(r byte) string {
	if t.currentState == "NULL" {
		return "NULL"
	}

	row, present := t.transition[t.currentState]
	if !present {
		t.currentState = "NULL"
		return t.currentState
	}

	nxt, present := row[r]
	if !present {
		t.currentState = "NULL"
		return t.currentState
	}

	t.currentState = nxt
	return nxt
}

// 	{
// 	nxt, err := t.transition(t.currentState, r)
// 	if err != nil {
// 		return d.currentState, err
// 	}
// 	d.currentState = nxt
// 	return nxt, err
// }
