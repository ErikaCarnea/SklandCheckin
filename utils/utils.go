package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func getLastRunFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return filepath.Join("mark", "SklandMarkFile")
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "mark", "SklandMarkFile")
}

func HasRunToday() bool {
	filename := getLastRunFilePath()
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	now := time.Now().Local()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return info.ModTime().After(todayStart)
}

func MarkRun() error {
	filename := getLastRunFilePath()
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

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
