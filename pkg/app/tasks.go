package app

import (
	"sync"
	"time"

	"github.com/HeathErika/Skland/api"
	"github.com/HeathErika/Skland/models"
	"github.com/HeathErika/Skland/utils"
	"github.com/rs/zerolog/log"
)

func (ctx *AppContext) PerformSignAttendance() {
	// 获取绑定列表
	bindings, err := ctx.BindidngAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		return
	}

	// 打印玩家信息
	// if err := ctx.PlayerAPI.PrintAllPlayersInfo(bindings); err != nil {
	// 	log.Error().Err(err).Msg("获取玩家信息失败")
	// 	return
	// }

	// 执行签到
	signErrors := ctx.signAttendance(bindings)
	if len(signErrors) > 0 {
		for _, err := range signErrors {
			log.Error().Err(err).Msg("签到失败")
		}
	}
}

func (ctx *AppContext) signAttendance(bindings []models.Binding) []error {
	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, len(bindings))
		errors []error
	)

	for _, binding := range bindings {
		wg.Add(1)
		go func(b models.Binding) {
			defer wg.Done()
			result, err := ctx.AttendanceAPI.SignAttendance(b)
			if err != nil {
				errCh <- err
				return
			}
			log.Info().Msgf("[%s] (%s) %s",
				b.ChannelName,
				b.NickName,
				utils.FormatAttendanceResult(result))
		}(binding)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		errors = append(errors, err)
	}
	return errors
}

func (ctx *AppContext) RunCheckinTasks() {
	var wg sync.WaitGroup
	for gameId := range api.SklandBoard {
		wg.Add(1)
		time.Sleep(2 * time.Second)
		logger := log.With().Int("gameId", gameId).Str("gameName", api.SklandBoard[gameId]).Logger()
		go func(id int) {
			defer wg.Done()
			_, err := ctx.CheckinAPI.Checkin(id)
			if err != nil {
				logger.Error().Err(err).Msg("检票失败")
				return
			}
			logger.Info().Msg("检票成功")
		}(gameId)
	}
	wg.Wait()
}
