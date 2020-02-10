package fsm

import (
	"fmt"
	"strings"
)

// NFAtoDFA converts NFA to DFA
func NFAtoDFA(nfa NFA) (DFA, error) {

	// Need to store a group of group of states. Set of set of states
	newstates := []states{
		states([]string{nfa.InitialState}),
	}

	transitionTable := make(map[string]map[rune]string)
	acceptingStates := make(states, 0)

	for i := 0; i < len(newstates); i++ {
		newstate := strings.Join(newstates[i], "")
		fmt.Println("Running for", newstates[i], newstate)

		for _, alpha := range nfa.Alphabet {

			// Set of states for this alphabet
			nfa.CurrentStates = newstates[i]
			nxt := nfa.Next(alpha)

			already := false
			for _, newst := range newstates {
				if newst.equal(nxt) {
					already = true
					break
				}
			}
			if !already {
				newstates = append(newstates, nxt)
			}

			if _, present := transitionTable[newstate]; !present {
				transitionTable[newstate] = make(map[rune]string)
			}
			transitionTable[newstate][alpha] = strings.Join(nxt, "")
		}

		for _, state := range newstates[i] {
			for _, st := range nfa.AcceptingStates {
				if state == st {
					acceptingStates.add(newstate)
					break
				}
			}
		}
	}

	newstatenames := make([]string, 0)
	for newst := range transitionTable {
		newstatenames = append(newstatenames, newst)
	}

	return DFA{
		AcceptingStates: acceptingStates,
		Alphabet:        nfa.Alphabet,
		CurrentState:    nfa.InitialState,
		InitialState:    nfa.InitialState,
		States:          newstatenames,
		Transition:      transitionTable,
	}, nil

}
