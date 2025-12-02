package statemachine

type stateBuffer []*State

func (sb *stateBuffer) add(s *State) {
	*sb = append((*sb), s)
}

func (sb *stateBuffer) clear() {
	(*sb) = (*sb)[:0]
}

func (sb *stateBuffer) containsMatch() bool {
	for _, state := range *sb {
		if state.c == MATCH_STATE {
			return true
		}
	}

	return false
}

type Evaluator interface {
	Load(nfa *State)
	IsMatch(string) bool
}

type nfaEvaluator struct {
	sb     *stateBuffer
	swapSb *stateBuffer
	listId int
}

func populateStateBuffer(sb *stateBuffer, nfa *State, listId int) {
	if nfa == nil || nfa.lastList == listId {
		return
	}

	nfa.lastList = listId

	if nfa.c == SPLIT_STATE {
		populateStateBuffer(sb, nfa.outOne, listId)
		populateStateBuffer(sb, nfa.outTwo, listId)

		return
	}

	sb.add(nfa)
}

func (evaluator *nfaEvaluator) populateStateBufferFromNfa(nfa *State) {
	populateStateBuffer(evaluator.sb, nfa, evaluator.listId)
}

func (evaluator *nfaEvaluator) populateSwapStateBufferFromNfa(nfa *State) {
	populateStateBuffer(evaluator.swapSb, nfa, evaluator.listId)
}

func (evaluator *nfaEvaluator) step(c int) {
	evaluator.listId++
	evaluator.swapSb.clear()

	for _, state := range *(evaluator.sb) {
		if state.c == c {
			evaluator.populateSwapStateBufferFromNfa(state.outOne)
		}
	}
}

func (evaluator *nfaEvaluator) IsMatch(expr string) bool {
	for _, character := range expr {
		evaluator.step(int(character))
		evaluator.swapBuffers()
	}

	return evaluator.sb.containsMatch()
}

func (evaluator *nfaEvaluator) swapBuffers() {
	temp := evaluator.sb
	evaluator.sb, evaluator.swapSb = evaluator.swapSb, temp
}

func (evaluator *nfaEvaluator) Load(nfa *State) {
	evaluator.listId++
	evaluator.populateStateBufferFromNfa(nfa)
}

func NewEvaluator(nfa *State) Evaluator {
	return &nfaEvaluator{
		sb:     &stateBuffer{},
		swapSb: &stateBuffer{},
		listId: 0,
	}
}
