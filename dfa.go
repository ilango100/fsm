package fsm

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Definition of DFA lives here. It allows DFA definition in the form of transition table.

// DFA defines information about DFA
type DFA struct {
	States          []string
	Alphabet        []rune
	InitialState    string
	AcceptingStates []string
	CurrentState    string
	Transition      map[string]map[rune]string
}

func (dfa DFA) String() string {
	str := fmt.Sprintf("\nStates:\t\t\t%v\n", dfa.States)
	str += "Alphabet:\t\t["
	for _, a := range dfa.Alphabet {
		str += fmt.Sprintf("%c ", a)
	}
	str = str[:len(str)-1]
	str += fmt.Sprintf("]\nInitialState:\t\t%s\nAcceptingStates:\t%v\n", dfa.InitialState, dfa.AcceptingStates)
	str = strings.ReplaceAll(str, "[]string", "")
	str = strings.ReplaceAll(str, "[]rune", "")
	return str
}

// FromCSV builds a DFA from Transcition Table specified in a csv file.
func FromCSV(file string) (*DFA, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	c := csv.NewReader(f)
	header, err := c.Read()
	if err != nil {
		return nil, err
	}
	if header[0] != "DFA" {
		return nil, fmt.Errorf("First record in header line must be 'DFA'")
	}

	// Parse the table to get all alphabets and states
	alphabet := make([]rune, 0)
	for i := range header[1:] {
		if utf8.RuneCountInString(header[i+1]) > 1 {
			return nil, fmt.Errorf("The header line records must contain single rune")
		}
		alpha := []rune(header[i+1])[0]
		for i := range alphabet {
			if alpha == alphabet[i] {
				return nil, fmt.Errorf("Duplicate character in header row")
			}
		}
		alphabet = append(alphabet, alpha)
	}

	states := make([]string, 0)
	initialState := ""
	acceptingStates := make([]string, 0)
	transitionTable := make(map[string]map[rune]string)

	for record, err := c.Read(); err == nil; record, err = c.Read() {
		state := record[0]

		if state[0] == '>' {
			state = state[1:]
			initialState = state
		}
		if state[0] == '*' {
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

		if _, p := transitionTable[state]; !p {
			transitionTable[state] = make(map[rune]string)
		}

		for i := range alphabet {
			transitionTable[state][alphabet[i]] = record[i+1]
		}

		states = append(states, state)

	}

	if initialState == "" {
		return nil, fmt.Errorf("No initial state specified")
	}

	return &DFA{
		States:          states,
		Alphabet:        alphabet,
		InitialState:    initialState,
		AcceptingStates: acceptingStates,
		CurrentState:    initialState,
		Transition:      transitionTable,
	}, nil

}

// FromTable creates a DFA from Transition Table
func FromTable(transitionTable map[string]map[rune]string) (*DFA, error) {
	// Parse the table to get all alphabets and states
	alphabet := make([]rune, 0)
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

	return &DFA{
		States:          states,
		Alphabet:        alphabet,
		InitialState:    initialState,
		AcceptingStates: acceptingStates,
		CurrentState:    initialState,
		Transition:      transitionTable,
	}, nil
}

// Reset resets the DFA to start state
func (dfa *DFA) Reset() string {
	dfa.CurrentState = dfa.InitialState
	return dfa.CurrentState
}

// IsAccepted checks whether the input so far is accepted
func (dfa *DFA) IsAccepted() bool {
	if dfa.CurrentState == "DEAD" {
		return false
	}

	for i := range dfa.AcceptingStates {
		if dfa.CurrentState == dfa.AcceptingStates[i] {
			return true
		}
	}
	return false
}

// Next proceeds to the next state of DFA according to transition table
func (dfa *DFA) Next(r rune) string {
	if dfa.CurrentState == "DEAD" {
		return "DEAD"
	}

	row, present := dfa.Transition[dfa.CurrentState]
	if !present {
		dfa.CurrentState = "DEAD"
		return dfa.CurrentState
	}

	nxt, present := row[r]
	if !present || nxt == "" {
		dfa.CurrentState = "DEAD"
		return dfa.CurrentState
	}

	dfa.CurrentState = nxt
	return nxt
}

// ToCSV writes the DFA definition to a csv file
func (dfa DFA) ToCSV(file string) error {
	fl, err := os.Create(file)
	if err != nil {
		return err
	}

	fl.WriteString("DFA")
	for i := range dfa.Alphabet {
		fmt.Fprintf(fl, ",%c", dfa.Alphabet[i])
	}
	fmt.Fprintln(fl)

	for state := range dfa.Transition {
		if state == dfa.InitialState {
			fmt.Fprint(fl, ">")
		}
		for i := range dfa.AcceptingStates {
			if state == dfa.AcceptingStates[i] {
				fmt.Fprint(fl, "*")
				break
			}
		}
		r := dfa.Transition[state]
		fmt.Fprint(fl, state)
		for _, j := range dfa.Alphabet {
			if nxt, p := r[j]; p {
				fmt.Fprintf(fl, ",%s", nxt)
			} else {
				fmt.Fprint(fl, ",DEAD")
			}
		}
		fmt.Fprintln(fl)
	}

	return nil
}
