package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run tasks in n goroutines and stops its work when receiving m errors from tasks
func Run(tasks []Task, n int, m int) error {
	// Place your code here
	if len(tasks) == 0 {
		return nil
	}
	if n < 1 { //nolint:gomnd
		return nil
	}

	var (
		taskCh      = make(chan Task)
		errCh       = make(chan error)
		wg          = &sync.WaitGroup{}
		errCounter  int
		doneCh      = make(chan struct{})
		resultError error
	)

	wg.Add(n)
	// start n workers
	for i := 0; i < n; i++ {
		go do(taskCh, doneCh, errCh, wg)
	}

	// queueing tasks for workers
	wg.Add(1) //nolint:gomnd
	go queueTasksToDo(tasks, taskCh, doneCh, wg)

	// processing errors from tasks
	wg.Add(1) //nolint:gomnd
	go processErrorsFromTasks(errCh, doneCh, m, errCounter, &resultError, wg)

	wg.Wait()
	return resultError
}

func processErrorsFromTasks(errCh chan error, doneCh chan struct{}, maxErrorsAllowed int, errCounter int, resultError *error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case err := <-errCh:
			if err != nil {
				errCounter++
				if errCounter >= maxErrorsAllowed {
					*resultError = ErrErrorsLimitExceeded
					close(doneCh)
					return
				}
			}
		case <-doneCh:
			return
		}
	}
}

func queueTasksToDo(tasks []Task, taskCh chan Task, doneCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range tasks {
		select {
		case taskCh <- v:
		case <-doneCh:
			return
		}
	}
	close(doneCh)
}

func do(taskCh chan Task, doneCh chan struct{}, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case task := <-taskCh:
			// run task
			err := task()
			if err != nil {
				select {
				case errCh <- err:
				case <-doneCh:
					return
				}
			}
		case <-doneCh:
			return
		}
	}
}
