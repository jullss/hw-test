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

	stopCh, errCh := make(chan bool), make(chan error)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	curIndex := 0

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-stopCh:
					return
				default:
				}

				mu.Lock()

				if curIndex >= len(tasks) {
					mu.Unlock()
					break
				}

				task := tasks[curIndex]
				curIndex++
				mu.Unlock()

				if err := task(); err != nil {
					errCh <- err
				}
			}
		}()
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
