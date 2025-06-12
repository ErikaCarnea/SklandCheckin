package app

import (
	"sync"
	"time"

	"github.com/ErikaCarnea/Skland/api"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/ErikaCarnea/Skland/utils"
	"github.com/rs/zerolog/log"
)

func (ctx *AppContext) PerformSignAttendance() {
	// 获取绑定列表
	bindings, err := ctx.BindidngAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		return
	}
	for _, binding := range bindings {
		// 查询签到信息
		ctx.AttendanceAPI.QueryAttendanceInfo(binding)
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
	var errors []error

	for _, b := range bindings {
		if ctx.AttendanceAPI.QueryAttendanceInfo(b) {
			log.Info().Msgf("[%s] (%s) 今天已经签到过了", b.ChannelName, b.NickName)
			continue
		}
		result, err := ctx.AttendanceAPI.SignAttendance(b)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		log.Info().Msgf("[%s] (%s) %s",
			b.ChannelName,
			b.NickName,
			utils.FormatAttendanceResult(result))
	}
	return errors
}

func (ctx *AppContext) RunCheckinTasks() {
	var wg sync.WaitGroup
	for gameID := range api.SklandBoard {
		wg.Add(1)
		time.Sleep(2 * time.Second)
		logger := log.With().Int("gameId", gameID).Str("gameName", api.SklandBoard[gameID]).Logger()
		go func(id int) {
			defer wg.Done()
			_, err := ctx.CheckinAPI.Checkin(id)
			if err != nil {
				logger.Error().Err(err).Msg("检票失败")
				return
			}
			logger.Info().Msg("检票成功")
		}(gameID)
	}
	wg.Wait()
}
