package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var wg sync.WaitGroup
	var errorCount int32
	taskChan := make(chan Task, len(tasks))
	stopChan := make(chan struct{})
	stopOnce := sync.Once{}

	// Функция рабочей горутины
	worker := func() {
		defer wg.Done()
		for {
			select {
			case task, ok := <-taskChan:
				if !ok {
					return
				}
				if err := task(); err != nil {
					if atomic.AddInt32(&errorCount, 1) >= int32(m) {
						stopOnce.Do(func() {
							close(stopChan)
						})
						return
					}
				}
			case <-stopChan:
				return
			}
		}
	}

	// Запускаем n рабочих горутин
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	// Отправляем задачи в канал
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// Ожидаем завершения всех горутин
	wg.Wait()

	if errorCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
