package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	jobs := make(chan Task)
	stopCh := make(chan bool)
	errCh := make(chan error)

	wg := sync.WaitGroup{}

	go func() {
		defer close(jobs)

		for _, task := range tasks {
			select {
			case <-stopCh:
				return
			default:
			}

			select {
			case <-stopCh:
				return
			case jobs <- task:
			}
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)

		go execute(jobs, stopCh, errCh, &wg)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var err error
	errCount := 0

	for range errCh {
		errCount++

		if errCount >= m && err == nil {
			err = ErrErrorsLimitExceeded
			close(stopCh)
		}
	}

	return err
}

func execute(jobs <-chan Task, stopCh <-chan bool, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-stopCh:
			return
		case task, ok := <-jobs:
			if !ok {
				return
			}
			select {
			case <-stopCh:
				return
			default:
			}
			if err := task(); err != nil {
				errCh <- err
			}
		}
	}
}
