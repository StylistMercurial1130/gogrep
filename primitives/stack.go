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
	return (*s)[:len(*s)-1]
}
