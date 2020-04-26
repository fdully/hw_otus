package hw06_pipeline_execution //nolint:golint,stylecheck

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(name string, f func(v I) I) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v I) I { return v }),
		g("Multiplier (* 2)", func(v I) I { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v I) I { return v.(int) + 100 }),
		g("Stringifier", func(v I) I { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, result, []string{"102", "104", "106", "108", "110"})
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("nothing to execute", func(t *testing.T) {
		in := make(Bi)
		require.Nil(t, ExecutePipeline(in, nil, []Stage{}...))
	})

	t.Run("capitalize words", func(t *testing.T) {
		in := make(Bi)
		data := []string{"Hi", " ", "my", " ", "друг", "!"}

		stages := []Stage{
			g("capitalize", func(v I) I {
				str, ok := v.(string)
				if !ok {
					return nil
				}
				return strings.ToUpper(str)
			}),
			g("dummy", func(v I) I {
				return v
			}),
			g("remove spaces", func(v I) I {
				str, ok := v.(string)
				if !ok {
					return nil
				}
				if str == " " {
					return nil
				}
				return v
			}),
		}
		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 3)
		for s := range ExecutePipeline(in, nil, stages...) {
			str, ok := s.(string)
			if ok {
				result = append(result, str)
			}
		}
		require.Equal(t, result, []string{"HI", "MY", "ДРУГ", "!"})
	})
}
