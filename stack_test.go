package set

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestStack_simple(t *testing.T) {
	s := makeStack[int]()
	must.True(t, s.empty())

	s.push(1)
	must.False(t, s.empty())

	value := s.pop()
	must.Eq(t, 1, value)
	must.True(t, s.empty())
}

func TestStack_complex(t *testing.T) {
	s := makeStack[byte]()

	s.push('a')
	s.push('b')
	s.push('c')
	s.push('d')
	s.push('e')
	s.push('f')

	must.Eq(t, 'f', s.pop())
	must.Eq(t, 'e', s.pop())
	must.Eq(t, 'd', s.pop())

	s.push('x')
	s.push('y')

	must.Eq(t, 'y', s.pop())

	s.push('z')

	must.Eq(t, 'z', s.pop())
	must.Eq(t, 'x', s.pop())
	must.Eq(t, 'c', s.pop())
	must.Eq(t, 'b', s.pop())
	must.Eq(t, 'a', s.pop())
	must.True(t, s.empty())
}
