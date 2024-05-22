package compute

import (
	errors "memory-db/internal/errors/compute"
	"strings"
)

type state uint8

const (
	initialState state = iota
	symbolState
	whitespaceState

	// must be the last
	statesNumber
)

type event uint8

const (
	foundSymbolEvent event = iota
	foundWhitespaceEvent

	// must be the last
	eventsNumber
)

type transition struct {
	nextState state
	action    func(r rune)
}

type fsm struct {
	sb          strings.Builder
	state       state
	transitions [statesNumber][eventsNumber]transition
	tokens      []string
}

func newFSM() *fsm {
	machine := &fsm{}
	machine.transitions = [statesNumber][eventsNumber]transition{
		initialState: {
			foundWhitespaceEvent: {
				nextState: whitespaceState,
			},
			foundSymbolEvent: {
				nextState: symbolState,
				action:    machine.handleSymbol,
			},
		},
		symbolState: {
			foundSymbolEvent: {
				nextState: symbolState,
				action:    machine.handleSymbol,
			},
			foundWhitespaceEvent: {
				nextState: whitespaceState,
				action:    machine.handleWhitespace,
			},
		},
		whitespaceState: {
			foundWhitespaceEvent: {
				nextState: whitespaceState,
			},
			foundSymbolEvent: {
				nextState: symbolState,
				action:    machine.handleSymbol,
			},
		},
	}

	return machine
}

func (f *fsm) parse(query string) ([]string, error) {
	for _, r := range query {
		if isValidSymbol(r) {
			f.handleEvent(foundSymbolEvent, r)
		} else if isWhitespace(r) {
			f.handleEvent(foundWhitespaceEvent, r)
		} else {
			return nil, errors.InvalidSymbolError
		}
	}

	f.handleEvent(foundWhitespaceEvent, ' ') // hack to write last symbols
	return f.tokens, nil
}

func (f *fsm) handleEvent(event event, r rune) {
	transit := f.transitions[f.state][event]
	f.state = transit.nextState

	if transit.action != nil {
		transit.action(r)
	}
}

func (f *fsm) handleSymbol(r rune) {
	f.sb.WriteRune(r)
}

func (f *fsm) handleWhitespace(rune) {
	if f.sb.Len() > 0 {
		f.tokens = append(f.tokens, f.sb.String())
		f.sb.Reset()
	}
}
