package set

type stack[T any] struct {
	top *object[T]
}

type object[T any] struct {
	item T
	next *object[T]
}

func makeStack[T any]() *stack[T] {
	return &stack[T]{top: nil}
}

func (s *stack[T]) push(item T) {
	obj := &object[T]{
		item: item,
		next: s.top,
	}
	s.top = obj
}

func (s *stack[T]) pop() T {
	item := s.top.item
	s.top = s.top.next
	return item
}

func (s *stack[T]) empty() bool {
	return s.top == nil
}