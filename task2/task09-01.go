package main

import (
	"fmt"
	"sync"
)

//题目：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

var (
	counter int
	mutex   sync.Mutex
	wg      sync.WaitGroup
)

func main() {

	numGoroutines := 10
	incrementsPerGoroutine := 1000
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
			fmt.Printf("协程 %d 完成\n", id)
		}(i)

	}
	wg.Wait()

	// 输出最终结果
	expected := numGoroutines * incrementsPerGoroutine
	fmt.Printf("最终计数器值: %d\n", counter)
	fmt.Printf("期望值: %d\n", expected)
	fmt.Printf("结果是否正确: %v\n", counter == expected)
}
