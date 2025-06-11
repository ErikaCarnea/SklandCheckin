package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
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
			log.Info().Msg("未找到标记文件")
			return false
		}
		log.Error().Err(err).Msg("获取文件信息失败，视为今日未运行")
		return false
	}
	cst := time.FixedZone("CST", 8*60*60)
	modTime := info.ModTime()
	if modTime.Location() != cst {
		modTime = modTime.In(cst)
	}

	now := time.Now().In(cst)

	log.Debug().
		Str("文件修改时间", modTime.String()).
		Str("当前运行时间", now.String()).
		Msg("时间比对")

	return isSameDay(modTime, now)
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
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
