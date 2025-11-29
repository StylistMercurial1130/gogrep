package statemachine

import "gogrep/primitives"

const MATCH_STATE = 257
const SPLIT_STATE = 256

type State struct {
	c        int
	outOne   *State
	outTwo   *State
	lastList int
}

type StateList []*State

func newStateList(s *State) StateList {
	return []*State{s}
}

func newStateListFromStates(s ...*State) StateList {
	var stateList []*State
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
	for _, state := range *stateList {
		state.outOne = s
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
				e1 := fragmentStack.Pop().(Frag)
				e2 := fragmentStack.Pop().(Frag)
				e1.stateList.patch(e2.state)
				fragmentStack.Push(e1)

				break
			}
		case '|':
			{
				e1 := fragmentStack.Pop().(Frag)
				e2 := fragmentStack.Pop().(Frag)
				splitState := createSplitState(e1.state, e2.state)
				fragmentStack.Push(
					createFragment(
						splitState,
						newStateListFromStates(splitState.outOne, splitState.outTwo),
					),
				)

				break
			}
		case '?':
			{
				e1 := fragmentStack.Pop().(Frag)
				s := createSplitState(e1.state, nil)
				fragmentStack.Push(
					createFragment(
						s,
						newStateListFromStates(e1.state.outOne, s.outTwo),
					),
				)

				break
			}
		case '*':
			{
				e1 := fragmentStack.Pop().(Frag)
				s := createSplitState(e1.state, nil)
				e1.stateList.patch(s)
				fragmentStack.Push(
					createFragment(
						s, newStateList(s.outTwo),
					),
				)
			}
		case '+':
			{
				e1 := fragmentStack.Pop().(Frag)
				s := createSplitState(e1.state, nil)
				e1.stateList.patch(s)
				fragmentStack.Push(
					createFragment(
						s, newStateList(s.outTwo),
					),
				)

				break
			}
		default:
			{
				s := createState(int(character))
				fragmentStack.Push(createFragment(s, newStateList(s.outOne)))

				break
			}
		}
	}

	e := fragmentStack.Pop().(Frag)
	e.stateList.patch(createMatchState())

	return e.state
}
