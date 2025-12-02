package primitives

type Stack []interface{}

func (s *Stack) Push(value interface{}) {
	*s = append(*s, value)
}

func (s *Stack) Pop() interface{} {
	value := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return value
}

func (s *Stack) Peek() interface{} {
	return (*s)[len(*s)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// PopOk pops the top element and returns it along with a boolean that
// indicates whether a value was returned. It does not panic on empty stack.
func (s *Stack) PopOk() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, true
}

// PeekOk returns the top element without removing it and a boolean which
// indicates whether a value was returned. It does not panic on empty stack.
func (s *Stack) PeekOk() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	return (*s)[len(*s)-1], true
}
