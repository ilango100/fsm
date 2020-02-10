package fsm

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Definition of NFA lives here. It allows NFA definition in the form of transition table.

// NFA defines Non Deterministic Finite Automaton
type NFA struct {
	States          []string
	Alphabet        []rune
	InitialState    string
	AcceptingStates []string
	CurrentStates   []string
	Transition      map[string]map[rune][]string
}

func (nfa NFA) String() string {
	str := fmt.Sprintf("\nStates:\t\t\t%v\n", nfa.States)
	str += "Alphabet:\t\t["
	for _, a := range nfa.Alphabet {
		str += fmt.Sprintf("%c ", a)
	}
	str = str[:len(str)-1]
	str += fmt.Sprintf("]\nInitialState:\t\t%s\nAcceptingStates:\t%v\n", nfa.InitialState, nfa.AcceptingStates)
	str = strings.ReplaceAll(str, "[]string", "")
	str = strings.ReplaceAll(str, "[]rune", "")
	return str
}

// NFAfromCSV builds an NFA from Transcition Table specified in a csv file.
func NFAfromCSV(file string) (*NFA, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	c := csv.NewReader(f)
	header, err := c.Read()
	if err != nil {
		return nil, err
	}
	if header[0] != "NFA" {
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
	transitionTable := make(map[string]map[rune][]string)

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

		for i, s := range states {
			if s == state {
				return nil, fmt.Errorf("Duplicate rows for state: %s, first found in row %d", s, i+1)
			}
		}

		if _, p := transitionTable[state]; !p {
			transitionTable[state] = make(map[rune][]string)
		}

		for i := range alphabet {
			stts := strings.Split(record[i+1], " ")
			stts[0] = stts[0][1:]
			stts[len(stts)-1] = stts[len(stts)-1][:len(stts[len(stts)-1])-1]
			for s := range stts {
				transitionTable[state][alphabet[i]] = append(transitionTable[state][alphabet[i]], stts[s])
			}
		}

		states = append(states, state)

	}

	if initialState == "" {
		return nil, fmt.Errorf("No initial state specified")
	}

	return &NFA{
		States:          states,
		Alphabet:        alphabet,
		InitialState:    initialState,
		AcceptingStates: acceptingStates,
		CurrentStates:   []string{initialState},
		// Transition:      transitionTable,
	}, nil

}

// NFAfromTable creates a DFA from Transition Table
func NFAfromTable(transitionTable map[string]map[rune][]string) (*NFA, error) {
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

	return &NFA{
		States:          states,
		Alphabet:        alphabet,
		InitialState:    initialState,
		AcceptingStates: acceptingStates,
		CurrentStates:   []string{initialState},
		Transition:      transitionTable,
	}, nil
}

// Reset resets the DFA to start state
func (nfa *NFA) Reset() []string {
	nfa.CurrentStates = []string{nfa.InitialState}
	return nfa.CurrentStates
}

// IsAccepted checks whether the input so far is accepted
func (nfa *NFA) IsAccepted() bool {
	for i := range nfa.AcceptingStates {
		for j := range nfa.CurrentStates {
			if nfa.AcceptingStates[i] == nfa.CurrentStates[j] {
				return true
			}
		}
	}
	return false
}

// Next proceeds to the next state of DFA according to transition table
func (nfa *NFA) Next(r rune) []string {

	nxts := states(make([]string, 0))
	for _, cstate := range nfa.CurrentStates {

		row, present := nfa.Transition[cstate]
		if !present {
			continue
		}

		nxt, present := row[r]
		if !present {
			continue
		}

		nxts.addmany(nxt)
	}

	// E-Close
	for i := range nxts {
		row, present := nfa.Transition[nxts[i]]
		if !present {
			continue
		}

		nxt, present := row[0]
		if !present {
			continue
		}
		nxts.addmany(nxt)
	}
	// Run twice because of sorting in states
	// TODO: Fix
	for i := range nxts {
		row, present := nfa.Transition[nxts[i]]
		if !present {
			continue
		}

		nxt, present := row[0]
		if !present {
			continue
		}
		nxts.addmany(nxt)
	}

	nfa.CurrentStates = nxts
	return nxts
}

// ToCSV writes the DFA definition to a csv file
func (nfa NFA) ToCSV(file string) error {
	fl, err := os.Create(file)
	if err != nil {
		return err
	}

	fl.WriteString("NFA")
	for i := range nfa.Alphabet {
		fmt.Fprintf(fl, ",%c", nfa.Alphabet[i])
	}
	fmt.Fprintln(fl)

	for state := range nfa.Transition {
		if state == nfa.InitialState {
			fmt.Fprint(fl, ">")
		}
		for i := range nfa.AcceptingStates {
			if state == nfa.AcceptingStates[i] {
				fmt.Fprint(fl, "*")
				break
			}
		}
		r := nfa.Transition[state]
		fmt.Fprint(fl, state)
		for _, j := range nfa.Alphabet {
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
