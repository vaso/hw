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
	tChan := make(chan Task) // channel with array of tasks to complete
	errCount.Store(0)

	log.Println("start Run")

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer(tasks, tChan, int32(m))
	}()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer(tChan, int32(m))
		}()
	}
	wg.Wait()
	if errCount.Load() >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func producer(tasks []Task, tChan chan<- Task, maxErrCount int32) {
	for _, task := range tasks {
		if errCount.Load() >= maxErrCount {
			close(tChan) // closing tasks channel is a flag for all consumers to stop.
			return
		}

		tChan <- task
	}
	close(tChan)
}

func consumer(tChan <-chan Task, maxErrCount int32) {
	for {
		task, ok := <-tChan
		if !ok {
			return
		}

		err := task()
		if err == nil {
			continue
		}

		errCount.Add(int32(1))
		if errCount.Load() >= maxErrCount {
			select {
			case <-tChan: // read from tChan to release producer for next iteration
			default:
			}
			return
		}
	}
}
