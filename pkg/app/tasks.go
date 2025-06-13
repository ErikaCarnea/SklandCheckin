package app

import (
	"github.com/ErikaCarnea/Skland/api"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/ErikaCarnea/Skland/utils"
	"github.com/rs/zerolog/log"
)

func (ctx *AppContext) PerformSignAttendance() bool {
	// 获取绑定列表
	bindings, err := ctx.BindidngAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		return false
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
		return false
	}
	return true
}

func (ctx *AppContext) signAttendance(bindings []models.Binding) []error {
	var errors []error

	for _, b := range bindings {
		if ctx.AttendanceAPI.QueryAttendanceInfo(b) {
			continue
		}
		result, err := ctx.AttendanceAPI.SignAttendance(b)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		log.Info().Msgf("[%s] %s %s",
			b.ChannelName,
			b.NickName,
			utils.FormatAttendanceResult(result))
	}
	return errors
}

func (ctx *AppContext) RunCheckinTasks() bool {
	allRepeated := true

	for gameID, gameName := range api.SklandBoard {
		logger := log.With().Int("gameId", gameID).Str("gameName", gameName).Logger()
		resp, err := ctx.CheckinAPI.Checkin(gameID)

		if resp == nil {
			allRepeated = false
			logger.Error().Err(err).Msg("API请求失败")
			continue
		}

		if resp.Code == 10001 {
			continue
		}

		allRepeated = false
		if err != nil {
			logger.Error().Err(err).Msg("检票失败")
		} else {
			logger.Info().Msg("检票成功")
		}
	}
	return allRepeated
}
