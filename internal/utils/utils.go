package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// WaitForExit 等待用户按回车退出（10 秒超时）。
// 使用 zerolog 而非 fmt.Println，保持输出一致性。
func WaitForExit() {
	log.Info().Msg("签到执行完毕，按回车键退出程序（10秒后自动退出）...")

	inputCh := make(chan int, 1)
	go func() {
		n, err := fmt.Fscanln(os.Stdin)
		if err != nil {
			n = 0
		}
		inputCh <- n
	}()

	select {
	case <-inputCh:
		log.Info().Msg("程序退出")
	case <-time.After(10 * time.Second):
		log.Info().Msg("等待超时，程序自动退出")
	}
}
