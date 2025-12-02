package statemachine

import (
	"gogrep/primitives"
	"strings"
)

const MATCH_STATE = 257
const SPLIT_STATE = 256

type State struct {
	c        int
	outOne   *State
	outTwo   *State
	lastList int
}
type StateList []**State

func newStateList(s **State) StateList {
	return []**State{s}
}

func newStateListFromStates(s ...**State) StateList {
	var stateList []**State
	stateList = append(stateList, s...)
	return stateList
}

func createMatchState() *State {
	return &State{
		c: MATCH_STATE,
	}
}

func createState(c int) *State {
	return &State{
		c: c,
	}
}

func createSplitState(outOne, outTwo *State) *State {
	return &State{
		c:      SPLIT_STATE,
		outOne: outOne,
		outTwo: outTwo,
	}
}

func (stateList *StateList) concat(other *StateList) {
	*stateList = append(*stateList, *other...)
}

func (stateList *StateList) patch(s *State) {
	for _, statePtr := range *stateList {
		// statePtr is **State (pointer to a *State field). Assign the
		// provided *State into that field so the dangling pointer is fixed.
		*statePtr = s
	}
}

type Frag struct {
	state     *State
	stateList StateList
}

func createFragment(state *State, stateList StateList) Frag {
	return Frag{
		state: state, stateList: stateList,
	}
}

func postfix2Nfa(postfix string) *State {
	var fragmentStack primitives.Stack

	for _, character := range postfix {
		switch character {
		case '.':
			{
				v, ok := fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e2 for concat")
				}
				e2 := v.(Frag)
				v, ok = fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e1 for concat")
				}
				e1 := v.(Frag)
				e1.stateList.patch(e2.state)
				fragmentStack.Push(createFragment(e1.state, e2.stateList))

				break
			}
		case '|':
			{
				v, ok := fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e1 for union")
				}
				e1 := v.(Frag)
				v, ok = fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e2 for union")
				}
				e2 := v.(Frag)
				splitState := createSplitState(e1.state, e2.state)
				fragmentStack.Push(
					createFragment(
						splitState,
						newStateListFromStates(&splitState.outOne, &splitState.outTwo),
					),
				)

				break
			}
		case '?':
			{
				v, ok := fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e1 for question")
				}
				e1 := v.(Frag)
				s := createSplitState(e1.state, nil)
				fragmentStack.Push(
					createFragment(
						s,
						newStateListFromStates(&e1.state.outOne, &s.outTwo),
					),
				)

				break
			}
		case '*':
			{
				v, ok := fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e1 for star")
				}
				e1 := v.(Frag)
				s := createSplitState(e1.state, nil)
				e1.stateList.patch(s)
				fragmentStack.Push(
					createFragment(
						s, newStateList(&s.outTwo),
					),
				)

				break
			}
		case '+':
			{
				v, ok := fragmentStack.PopOk()
				if !ok {
					panic("postfix2Nfa: stack underflow while popping e1 for plus")
				}
				e1 := v.(Frag)
				s := createSplitState(e1.state, nil)
				e1.stateList.patch(s)
				fragmentStack.Push(
					createFragment(
						s, newStateList(&s.outTwo),
					),
				)

				break
			}
		default:
			{
				s := createState(int(character))
				fragmentStack.Push(createFragment(s, newStateList(&s.outOne)))

				break
			}
		}
	}

	v, ok := fragmentStack.PopOk()
	if !ok {
		panic("postfix2Nfa: stack underflow while popping final fragment")
	}
	e := v.(Frag)
	matchState := createMatchState()
	e.stateList.patch(matchState)

	return e.state
}

func isOperator(c rune) bool {
	return c == '+' || c == '*' || c == '|' || c == '.' || c == '?'
}

func isLiteral(c rune) bool {
	return !(isOperator(c) || c == '(' || c == ')')
}

func getPrec(c rune) int {
	switch c {
	case '*', '+', '?':
		return 3
	case '.':
		return 2
	case '|':
		return 1
	default:
		return -1
	}
}

func insertConcatonation(str string) string {
	var result strings.Builder

	for i := 0; i < len(str); i++ {
		result.WriteRune(rune(str[i]))

		if i+1 < len(str) {
			if (isLiteral(rune(str[i])) || str[i] == '+' || str[i] == '*' || str[i] == '?' || str[i] == ')') &&
				(isLiteral(rune(str[i+1])) || str[i+1] == '(') {
				result.WriteRune('.')
			}
		}
	}

	return result.String()
}

func ExprToPostFix(str string) string {
	var opStack primitives.Stack
	var strBuilder strings.Builder

	s := insertConcatonation(str)

	for _, c := range s {
		if isLiteral(c) {
			strBuilder.WriteRune(c)
		} else if c == '(' {
			opStack.Push(c)
		} else if c == ')' {
			for {
				peek, ok := opStack.PeekOk()
				if !ok || peek == '(' {
					break
				}

				if character, ok := opStack.PopOk(); ok {
					if ch, ok2 := character.(rune); ok2 {
						strBuilder.WriteRune(ch)
					}
				} else {
					break
				}
			}

			// pop the '('
			opStack.PopOk()
		} else if isOperator(c) {
			// Postfix operators should be pushed immediately without precedence checks
			if c == '*' || c == '+' || c == '?' {
				opStack.Push(c)
			} else {
				// Infix operators: apply precedence rules
				for {
					pv, ok := opStack.PeekOk()
					if !ok {
						break
					}
					peek, ok2 := pv.(rune)
					if !ok2 || peek == '(' {
						break
					}

					if getPrec(peek) >= getPrec(c) {
						if charv, ok := opStack.PopOk(); ok {
							if ch, ok2 := charv.(rune); ok2 {
								strBuilder.WriteRune(ch)
							}
						} else {
							break
						}
					} else {
						break
					}
				}
				opStack.Push(c)
			}
		}
	}

	for {
		if character, ok := opStack.PopOk(); ok {
			if ch, ok2 := character.(rune); ok2 {
				strBuilder.WriteRune(ch)
			}
		} else {
			break
		}
	}

	return strBuilder.String()
}
