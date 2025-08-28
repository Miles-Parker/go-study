package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//题目：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。

// Task 定义任务类型
type Task func(ctx context.Context) error

// TaskResult 任务执行结果
type TaskResult struct {
	TaskID    int
	TaskName  string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Err       error
	IsSuccess bool
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	tasks         []*taskWrapper
	timeout       time.Duration
	maxConcurrent int
	mu            sync.Mutex
	taskCounter   int32
}

// taskWrapper 任务包装器，包含任务和元数据
type taskWrapper struct {
	id       int
	name     string
	task     Task
	priority int // 优先级，数字越大优先级越高
}

// SchedulerOption 调度器配置选项
type SchedulerOption func(*TaskScheduler)

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler(options ...SchedulerOption) *TaskScheduler {
	scheduler := &TaskScheduler{
		timeout:       30 * time.Second, // 默认超时时间
		maxConcurrent: 10,               // 默认最大并发数
	}

	for _, option := range options {
		option(scheduler)
	}

	return scheduler
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) SchedulerOption {
	return func(s *TaskScheduler) {
		s.timeout = timeout
	}
}

// WithMaxConcurrent 设置最大并发数
func WithMaxConcurrent(max int) SchedulerOption {
	return func(s *TaskScheduler) {
		s.maxConcurrent = max
	}
}

// AddTask 添加任务
func (ts *TaskScheduler) AddTask(name string, task Task) {
	ts.AddTaskWithPriority(name, task, 0)
}

// AddTaskWithPriority 添加带优先级的任务
func (ts *TaskScheduler) AddTaskWithPriority(name string, task Task, priority int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	id := atomic.AddInt32(&ts.taskCounter, 1)
	ts.tasks = append(ts.tasks, &taskWrapper{
		id:       int(id),
		name:     name,
		task:     task,
		priority: priority,
	})
}

// Run 执行所有任务
func (ts *TaskScheduler) Run() []TaskResult {
	if len(ts.tasks) == 0 {
		return nil
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	// 控制并发数量的信号量
	sem := make(chan struct{}, ts.maxConcurrent)
	results := make([]TaskResult, len(ts.tasks))
	var wg sync.WaitGroup

	// 按照优先级排序任务（简单的优先级处理）
	ts.sortTasksByPriority()

	for i, taskWrapper := range ts.tasks {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(idx int, tw *taskWrapper) {
			defer func() {
				<-sem // 释放信号量
				wg.Done()
			}()

			// 执行单个任务并记录结果
			results[idx] = ts.executeTask(ctx, tw)
		}(i, taskWrapper)
	}

	wg.Wait()
	return results
}

// executeTask 执行单个任务
func (ts *TaskScheduler) executeTask(ctx context.Context, tw *taskWrapper) TaskResult {
	startTime := time.Now()
	result := TaskResult{
		TaskID:    tw.id,
		TaskName:  tw.name,
		StartTime: startTime,
	}

	defer func() {
		// 捕获panic
		if r := recover(); r != nil {
			result.Err = fmt.Errorf("task panicked: %v", r)
			result.IsSuccess = false
		}
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(startTime)
	}()

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		result.Err = ctx.Err()
		result.IsSuccess = false
		return result
	default:
	}

	// 执行任务
	err := tw.task(ctx)
	result.Err = err
	result.IsSuccess = err == nil

	return result
}

// sortTasksByPriority 按优先级排序任务
func (ts *TaskScheduler) sortTasksByPriority() {
	// 简单的冒泡排序，对于任务数量不多的情况足够用
	for i := 0; i < len(ts.tasks)-1; i++ {
		for j := 0; j < len(ts.tasks)-i-1; j++ {
			if ts.tasks[j].priority < ts.tasks[j+1].priority {
				ts.tasks[j], ts.tasks[j+1] = ts.tasks[j+1], ts.tasks[j]
			}
		}
	}
}

// GetStatistics 获取统计信息
func (ts *TaskScheduler) GetStatistics(results []TaskResult) map[string]interface{} {
	stats := make(map[string]interface{})
	var totalDuration time.Duration
	successCount := 0
	failCount := 0

	for _, result := range results {
		totalDuration += result.Duration
		if result.IsSuccess {
			successCount++
		} else {
			failCount++
		}
	}

	stats["total_tasks"] = len(results)
	stats["success_count"] = successCount
	stats["fail_count"] = failCount
	stats["success_rate"] = float64(successCount) / float64(len(results)) * 100
	stats["total_duration"] = totalDuration
	stats["avg_duration"] = totalDuration / time.Duration(len(results))

	return stats
}

// PrintResults 打印任务执行结果
func (ts *TaskScheduler) PrintResults(results []TaskResult) {
	fmt.Printf("\n=== 任务执行结果统计 ===\n")
	fmt.Printf("总任务数: %d\n", len(results))

	for _, result := range results {
		status := "✓ 成功"
		if !result.IsSuccess {
			status = fmt.Sprintf("✗ 失败: %v", result.Err)
		}

		fmt.Printf("任务 %d (%s): %s | 耗时: %v\n",
			result.TaskID, result.TaskName, status, result.Duration.Round(time.Millisecond))
	}

	stats := ts.GetStatistics(results)
	fmt.Printf("\n=== 汇总统计 ===\n")
	fmt.Printf("成功任务: %d/%d\n", stats["success_count"], stats["total_tasks"])
	fmt.Printf("成功率: %.1f%%\n", stats["success_rate"])
	fmt.Printf("总执行时间: %v\n", stats["total_duration"].(time.Duration).Round(time.Millisecond))
	fmt.Printf("平均任务时间: %v\n", stats["avg_duration"].(time.Duration).Round(time.Millisecond))
}

// Reset 重置调度器状态
func (ts *TaskScheduler) Reset() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tasks = nil
	atomic.StoreInt32(&ts.taskCounter, 0)
}

// 示例任务函数
func createExampleTask(name string, duration time.Duration, shouldFail bool) Task {
	return func(ctx context.Context) error {
		// 模拟任务执行时间
		select {
		case <-time.After(duration):
			if shouldFail {
				return fmt.Errorf("任务 %s 执行失败", name)
			}
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	// 创建任务调度器：5秒超时，最大并发数为3
	scheduler := NewTaskScheduler(
		WithTimeout(5*time.Second),
		WithMaxConcurrent(3),
	)

	// 添加各种示例任务
	scheduler.AddTask("快速任务", createExampleTask("快速任务", 100*time.Millisecond, false))
	scheduler.AddTask("中等任务", createExampleTask("中等任务", 500*time.Millisecond, false))
	scheduler.AddTaskWithPriority("高优先级任务", createExampleTask("高优先级任务", 200*time.Millisecond, false), 10)
	scheduler.AddTask("慢速任务", createExampleTask("慢速任务", 2*time.Second, false))
	scheduler.AddTask("失败任务", createExampleTask("失败任务", 300*time.Millisecond, true))
	scheduler.AddTask("超长任务", createExampleTask("超长任务", 10*time.Second, false)) // 这个会超时

	// 再添加几个任务
	for i := 0; i < 3; i++ {
		taskName := fmt.Sprintf("额外任务%d", i+1)
		scheduler.AddTask(taskName, createExampleTask(taskName, time.Duration(150+i*50)*time.Millisecond, false))
	}

	fmt.Println("开始执行任务...")
	startTime := time.Now()

	// 执行所有任务
	results := scheduler.Run()

	// 打印结果
	scheduler.PrintResults(results)
	fmt.Printf("\n总调度时间: %v\n", time.Since(startTime).Round(time.Millisecond))

	// 重置调度器以便重用
	scheduler.Reset()
}
