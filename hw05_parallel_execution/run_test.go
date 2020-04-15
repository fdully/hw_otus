package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)
		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t,
			runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("0 tasks in slice", func(t *testing.T) {
		result := Run([]Task{}, 2, 2)
		require.Nil(t, result)
	})

	t.Run("0 goroutines to run", func(t *testing.T) {
		result := Run([]Task{}, 0, 2)
		require.Nil(t, result)
	})

	t.Run("more goroutines than tasks were started with no errors from tasks", func(t *testing.T) {
		tasksCount := 25
		tasks := make([]Task, 0, tasksCount)

		var runTaskCount int32
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTaskCount, 1)
				return nil
			})
		}

		workersCount := 50
		maxErrorsCount := 0

		result := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, result)
		require.Equal(t, int32(tasksCount), runTaskCount)
	})

	t.Run("mix of tasks with errors and half without", func(t *testing.T) {
		tasksErrorCount := 25
		tasksNoErrorCount := 25
		tasks := make([]Task, 0, tasksErrorCount+tasksNoErrorCount)

		var runTaskErrorCounter int32
		var runTaskNoErrorCounter int32

		for i := 0; i < tasksErrorCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTaskErrorCounter, 1)
				return fmt.Errorf("error from task %d", i)
			})
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTaskNoErrorCounter, 1)
				return nil
			})
		}

		workersCount := 4

		result := Run(tasks, workersCount, tasksErrorCount+1)

		require.Nil(t, result)
		require.Equal(t, int32(tasksErrorCount), runTaskErrorCounter)
		require.Equal(t, int32(tasksNoErrorCount), runTaskNoErrorCounter)
	})
}
