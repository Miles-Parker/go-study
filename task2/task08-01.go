package main

import (
	"fmt"
	"sync"
)

//题目：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。

func main() {
	ch := make(chan int, 10)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 0; i < 100; i++ {
			ch <- i
			fmt.Printf("生产者发送: %d\n", i)
		}
		fmt.Println("生产者完成发送")
	}()

	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("消费者接收: %d\n", num)
		}
		fmt.Println("消费者完成接收")
	}()
	wg.Wait()
	fmt.Println("执行完毕！！")
}
