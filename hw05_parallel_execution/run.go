package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// если m < 0, то игнорирует возвращаемые ошибки при выполнении задач.
func Run(tasks []Task, n, m int) error {
	// счетчик ошибок при выполнении задач во всех обработчиках
	var errCnt atomic.Int32
	// признак остановки обработчиков из-за того, что число возникших ошибок при обработке >= maxErrCnt
	var stop atomic.Bool
	maxErrCnt := int32(m)
	tasksCh := make(chan Task, len(tasks))
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasksCh {
				if stop.Load() {
					return
				}
				if err := t(); maxErrCnt > 0 && err != nil {
					if errCnt.Add(1) >= maxErrCnt {
						stop.Store(true)
						return
					}
				}
			}
		}()
	}
	go func() {
		defer close(tasksCh)
		generateTasksCh(tasksCh, tasks, maxErrCnt, &errCnt)
	}()
	wg.Wait()
	if maxErrCnt > 0 && errCnt.Load() >= maxErrCnt {
		return ErrErrorsLimitExceeded
	}
	return nil
}

// generateTasksCh записать задачи из слайса в канал, остановить запись, если возникло maxErrCnt ошибок.
func generateTasksCh(out chan<- Task, tasks []Task, maxErrCnt int32, errCnt *atomic.Int32) {
	for _, t := range tasks {
		if maxErrCnt > 0 && errCnt.Load() >= maxErrCnt {
			return
		}
		out <- t
	}
}
