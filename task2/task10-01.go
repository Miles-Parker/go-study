package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//题目：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func main() {
	var (
		counter int64
		wg      sync.WaitGroup
	)

	// 启动10个协程
	numGoroutines := 10
	incrementsPerGoroutine := 1000

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				// 使用原子操作递增计数器
				atomic.AddInt64(&counter, 1)
			}
			fmt.Printf("协程 %d 完成\n", id)
		}(i)
	}
	// 等待所有协程完成
	wg.Wait()

	// 输出最终结果
	expected := int64(numGoroutines * incrementsPerGoroutine)
	fmt.Printf("最终计数器值: %d\n", counter)
	fmt.Printf("期望值: %d\n", expected)
	fmt.Printf("结果是否正确: %v\n", counter == expected)
}
