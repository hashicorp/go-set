package set

import (
	"testing"

	"github.com/shoenig/test/must"
)

func Test_set(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := makeStack(0)
		must.False(t, s.empty())
		must.Zero(t, s.pop())
		must.True(t, s.empty())
	})

	t.Run("push pop", func(t *testing.T) {
		s := makeStack("a")
		s.push("b")
		s.push("c")
		s.push("d")
		s.push("e")
		must.Eq(t, "e", s.pop())
		must.Eq(t, "d", s.pop())
		must.Eq(t, "c", s.pop())
		must.Eq(t, "b", s.pop())
		must.Eq(t, "a", s.pop())
		must.True(t, s.empty())
	})
}
