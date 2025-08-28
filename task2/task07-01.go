package main

import (
	"fmt"
	"time"
)

//题目：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来

func main() {
	// 创建无缓冲通道
	ch := make(chan int)

	// 发送者协程：发送1-10到通道
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond) // 模拟处理延迟
		}
		close(ch) // 发送完成后关闭通道
	}()

	// 接收者协程：从通道接收并打印
	go func() {
		for num := range ch {
			fmt.Printf("Received: %d\n", num)
		}
	}()

	// 主协程等待防止程序立即退出
	time.Sleep(2 * time.Second)
}
