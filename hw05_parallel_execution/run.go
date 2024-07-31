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
	// будут выполнены все задачи независимо от количества возвращаемых ошибок
	noErrLimit := m < 0
	if noErrLimit {
		return executeTasksIgnoreErrs(tasks, n)
	}
	return executeTasksCheckErrs(tasks, n, m)
}

// executeTasksCheckErrs выполнить задачи tasks в n горутинах и остановить выполнение, если возникло m ошибок.
// Если возникло m ошибок, то вернет ErrErrorsLimitExceeded.
func executeTasksCheckErrs(tasks []Task, n, m int) error {
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
				if err := t(); err != nil {
					if errCnt.Add(1) >= maxErrCnt {
						stop.Store(true)
					}
				}
			}
		}()
	}
	for _, t := range tasks {
		tasksCh <- t
	}
	close(tasksCh)
	wg.Wait()
	if errCnt.Load() >= maxErrCnt {
		return ErrErrorsLimitExceeded
	}
	return nil
}

// executeTasksCheckErrs выполнить задачи tasks в n горутинах, игнорируя ошибки.
func executeTasksIgnoreErrs(tasks []Task, n int) error {
	tasksCh := make(chan Task, len(tasks))
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasksCh {
				_ = t()
			}
		}()
	}
	for _, t := range tasks {
		tasksCh <- t
	}
	close(tasksCh)
	wg.Wait()
	return nil
}
