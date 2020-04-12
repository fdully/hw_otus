package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks
func Run(tasks []Task, n int, m int) error {
	// Place your code here
	const (
		one = 1
	)
	if len(tasks) == 0 {
		return nil
	}
	if n < one {
		return nil
	}

	var (
		taskCh     = make(chan Task)
		errCh      = make(chan error)
		wg         = &sync.WaitGroup{}
		errCounter int
		doneCh     = make(chan struct{})
	)

	for i := 0; i < n; i++ {
		wg.Add(one)
		go do(taskCh, doneCh, errCh, wg)
	}

	// passing tasks to task channel
	wg.Add(one)
	go passTasksToDo(tasks, taskCh, doneCh, wg)

	// processing errors from tasks
	for i := 0; i < len(tasks); i++ {
		err := <-errCh
		if err != nil {
			errCounter++
			if errCounter >= m {
				close(doneCh)
				wg.Wait()
				return ErrErrorsLimitExceeded
			}
		}
	}

	// exiting all goroutines
	close(doneCh)

	wg.Wait()

	return nil
}

func passTasksToDo(tasks []Task, taskCh chan Task, doneCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range tasks {
		select {
		case _, ok := <-doneCh:
			if !ok {
				return
			}
		default:
		}
		select {
		case taskCh <- v:
		case _, ok := <-doneCh:
			if !ok {
				return
			}
		}
	}
}

func do(taskCh chan Task, doneCh chan struct{}, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case _, ok := <-doneCh:
			if !ok {
				return
			}
		default:
		}

		select {
		case task := <-taskCh:
			select {
			// run task
			case errCh <- task():
			case _, ok := <-doneCh:
				if !ok {
					return
				}
			}
		case _, ok := <-doneCh:
			if !ok {
				return
			}
		}
	}
}
