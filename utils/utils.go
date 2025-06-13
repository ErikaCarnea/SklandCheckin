package utils

import (
	"fmt"
	"time"
)

func WaitForExit() {
	fmt.Println("\n签到执行完毕，按回车键退出程序（10秒后自动退出）...")

	inputCh := make(chan int)

	go func() {
		n, _ := fmt.Scanln()
		inputCh <- n
	}()

	select {
	case <-inputCh:
		fmt.Println("程序退出")
	case <-time.After(10 * time.Second):
		fmt.Println("\n等待超时，程序自动退出")
	}
}
