package hw05parallelexecution

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var errCount atomic.Int32

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	tChan := make(chan Task)          // channel with array of tasks to complete
	errChan := make(chan struct{}, 1) // flag to stop all go routines due to error limit exceeding
	errCount.Store(0)

	log.Println("start Run")

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer(tasks, tChan, errChan)
	}()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer(tChan, errChan, int32(m))
		}()
	}
	wg.Wait()
	close(errChan)
	if errCount.Load() >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func producer(tasks []Task, tChan chan<- Task, errChan chan struct{}) {
	for _, task := range tasks {
		select {
		case <-errChan: // read errChan for exit flag
			close(tChan) // closing tasks channel is a flag for all consumers to stop.
			return
		case tChan <- task: // send task to tChan, blocking
			continue
		}
	}
	close(tChan)
}

// push err flag to errChan. If it already has a flag, no need to push it again.
func writeToErrChan(errChan chan struct{}) {
	select {
	case errChan <- struct{}{}:
	default:
	}
}

func consumer(tChan <-chan Task, errChan chan struct{}, maxErrCount int32) {
	for {
		select {
		case <-errChan:
			writeToErrChan(errChan) // we read err flag, but we need to re-enable it until producer gets it.
			return
		case task, ok := <-tChan: // read task to execute, handle task
			if !ok {
				return
			}
			err := task()
			if err == nil {
				continue
			}

			errCount.Add(int32(1))
			if errCount.Load() >= maxErrCount {
				writeToErrChan(errChan)
				return
			}
		}
	}
}
