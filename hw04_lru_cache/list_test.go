package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		middle := l.Back().Next // 20
		l.Remove(middle)        // [10, 30]
		require.Equal(t, l.Len(), 2)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 30, 10, 40, 60, 80, 70}, elems)

		l.Remove(l.Back()) // [70, 80, 60, 40, 10, 30]
		require.Equal(t, l.Back().Value, 30)
		require.Equal(t, l.Back().Next.Value, 10)
		l.Remove(l.Front()) // [80, 60, 40, 10, 30]
		require.Equal(t, l.Front().Value, 80)
		require.Equal(t, l.Front().Prev.Value, 60)
		l.Remove(l.Front().Prev) // [80, 40, 10, 30]
		l.Remove(l.Back().Next)  // [80, 40, 30]
		require.Equal(t, l.Len(), 3)
		require.Equal(t, l.Back().Value, 30)
		require.Equal(t, l.Back().Next.Value, 40)
		require.Equal(t, l.Front().Value, 80)
		require.Equal(t, l.Front().Prev.Value, 40)

		lt := l.Len()
		for i := 0; i < lt; i++ {
			l.Remove(l.Front())
		}
		require.Equal(t, l.Len(), 0)

		l.Remove(l.Back())
		require.Equal(t, l.Len(), 0)
	})
}
