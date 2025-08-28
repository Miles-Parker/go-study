package main

import (
	"fmt"
	"sync"
	"time"
)

//题目：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。

type Task struct {
	ID       int
	Function func() error
}
type TaskResult struct {
	TaskID    int
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
	Err       error
}

type Scheduler struct {
	tasks      chan Task
	results    chan TaskResult
	wg         sync.WaitGroup
	workerPool int
}

func (s *Scheduler) Start() {
	for i := 0; i < s.workerPool; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Scheduler) worker() {
	defer s.wg.Done()
	for task := range s.tasks {
		err := task.Function()
		start := time.Now()
		end := time.Now()
		s.results <- TaskResult{
			TaskID:    task.ID,
			Duration:  end.Sub(start),
			StartTime: start,
			EndTime:   end,
			Err:       err,
		}
	}
}
func (s *Scheduler) Submit(task Task) {
	s.tasks <- task
}

func (s *Scheduler) WaitAndPrint() {
	close(s.tasks)
	s.wg.Wait()
	close(s.results)

	for result := range s.results {
		if result.Err != nil {
			fmt.Printf("任务 %d 执行失败: %v\n", result.TaskID, result.Err)
		} else {
			fmt.Printf("任务 %d 执行成功 | 耗时: %v | 开始: %v | 结束: %v\n",
				result.TaskID, result.Duration, result.StartTime.Format("15:04:05.000"), result.EndTime.Format("15:04:05.000"))
		}
	}
}

func NewScheduler(workerPool int) *Scheduler {
	return &Scheduler{
		tasks:      make(chan Task, 100),
		results:    make(chan TaskResult, 100),
		workerPool: workerPool,
	}
}
func main() {
	scheduler := NewScheduler(5)
	tasks := []Task{
		{
			ID: 1,
			Function: func() error {
				time.Sleep(1 * time.Second)
				return nil
			},
		},
		{
			ID: 2,
			Function: func() error {
				time.Sleep(2 * time.Second)
				return nil
			},
		},
		{
			ID: 3,
			Function: func() error {
				time.Sleep(500 * time.Millisecond)
				return fmt.Errorf("模拟错误")
			},
		},
	}

	scheduler.Start()

	for _, task := range tasks {
		scheduler.Submit(task)
	}

	scheduler.WaitAndPrint()

}
