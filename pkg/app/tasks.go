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
	var signErrors []error
	// 打印玩家信息
	// if err := ctx.PlayerAPI.PrintAllPlayersInfo(bindings); err != nil {
	// 	log.Error().Err(err).Msg("获取玩家信息失败")
	// 	return
	// }

	// 执行签到
	for _, data := range bindings.Data.List {
		if data.AppCode == "arknights" {
			for gameID, game := range api.SklandBoard {
				if game == "明日方舟" {
					signErrors = ctx.signAttendance(data.BindingList, gameID)
				}
			}
		}
		if data.AppCode == "endfield" {
			for gameID, game := range api.SklandBoard {
				if game == "明日方舟: 终末地" {
					signErrors = ctx.signAttendance(data.BindingList, gameID)
				}
			}
		}
	}

	if len(signErrors) > 0 {
		for _, err := range signErrors {
			log.Error().Err(err).Msg("签到失败")
		}
		return false
	}
	return true
}

func (ctx *AppContext) signAttendance(bindings []models.Binding, gameID int) []error {
	var errors []error

	for _, b := range bindings {
		if ctx.AttendanceAPI.QueryAttendanceInfo(b, gameID) {
			log.Debug().Msgf("[%s] %s 今日已签到",
				b.ChannelName,
				b.NickName)
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
	allTasksCompleted := true

	// 获取所有游戏的签到状态
	statusMap, err := ctx.CheckinAPI.GetAllCheckinStatus()
	if err != nil {
		// 如果获取状态失败，我们仍然尝试签到，但记录错误
		log.Warn().Err(err).Msg("无法获取检票状态，将逐个检查")
		statusMap = nil // 设置为nil，表示无法使用状态检查
	}

	for gameID, gameName := range api.SklandBoard {
		logger := log.With().Int("gameId", gameID).Str("gameName", gameName).Logger()

		if statusMap != nil {
			checked, ok := statusMap[gameID]
			if ok && checked {
				logger.Debug().Msg("已检票")
				continue
			}
		}

		// 执行签到
		resp, err := ctx.CheckinAPI.Checkin(gameID)
		if err != nil {
			allTasksCompleted = false
			logger.Error().Err(err).Msg("检票失败")
			continue
		}

		// 检查签到结果
		switch resp.Code {
		case 0:
			// 签到成功
			logger.Info().Msg("检票成功")
			// 这里注意：签到成功也算完成任务
		case 10001:
			// 重复签到（竞态条件）
			logger.Info().Msg("今日已检票")
			// 重复签到也算完成任务
		default:
			// 其他错误
			allTasksCompleted = false
			logger.Error().Int("code", resp.Code).Str("message", resp.Message).Msg("检票失败")
		}
	}
	return allTasksCompleted
}

func (ctx *AppContext) GetPopucomAchievement() {
	bindings, err := ctx.BindidngAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
	}
	for _, data := range bindings.Data.List {
		if data.AppCode == "popucom" {
			ctx.PopucomAPI.GetPopucomAchievement(data.BindingList, ctx.CredResult.Data.UserId)
		}
	}
}

func (ctx *AppContext) GetExaAchievement() {
	bindings, err := ctx.BindidngAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
	}
	for _, data := range bindings.Data.List {
		if data.AppCode == "exa" {
			ctx.ExaAPI.GetExastrisAchievement(data.BindingList, ctx.CredResult.Data.UserId)
		}
	}
}
