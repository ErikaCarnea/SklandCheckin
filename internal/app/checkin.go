package app

import (
	"context"

	"github.com/ErikaCarnea/Skland/internal/api"
	"github.com/ErikaCarnea/Skland/internal/models"
	"github.com/ErikaCarnea/Skland/internal/utils"
	"github.com/rs/zerolog/log"
)

// CheckinService 签到服务，封装所有签到、检票、成就相关逻辑。
type CheckinService struct {
	bindingAPI    *api.BindingAPI
	attendanceAPI *api.AttendanceAPI
	checkinAPI    *api.CheckinAPI
	playerAPI     *api.PlayerAPI
	popucomAPI    *api.PopucomAPI
	exaAPI        *api.ExaAPI
}

func NewCheckinService(
	bindingAPI *api.BindingAPI,
	attendanceAPI *api.AttendanceAPI,
	checkinAPI *api.CheckinAPI,
	playerAPI *api.PlayerAPI,
	popucomAPI *api.PopucomAPI,
	exaAPI *api.ExaAPI,
) *CheckinService {
	return &CheckinService{
		bindingAPI:    bindingAPI,
		attendanceAPI: attendanceAPI,
		checkinAPI:    checkinAPI,
		playerAPI:     playerAPI,
		popucomAPI:    popucomAPI,
		exaAPI:        exaAPI,
	}
}

// RunAll 执行全部签到流程。
func (c *CheckinService) RunAll(ctx context.Context, credResult *models.CredResult) {
	bindings, err := c.bindingAPI.GetBindingList(ctx)
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		return
	}

	c.getPopucomAchievement(ctx, bindings, credResult.Data.UserId)
	c.getExaAchievement(ctx, bindings, credResult.Data.UserId)

	hasSigned := c.PerformSignAttendance(ctx, bindings)
	hasCheckedIn := c.RunCheckinTasks(ctx)

	if hasSigned && hasCheckedIn {
		log.Info().Msg("今日已签到")
	}
}

func (c *CheckinService) PerformSignAttendance(ctx context.Context, bindings models.BindingResult) bool {
	var signErrors []error

	for _, data := range bindings.Data.List {
		switch data.AppCode {
		case "arknights":
			signErrors = append(signErrors, c.signAttendance(ctx, data.BindingList, models.GameArknights)...)
		case "endfield":
			signErrors = append(signErrors, c.signAttendance(ctx, data.BindingList, models.GameEndfield)...)
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

func (c *CheckinService) signAttendance(ctx context.Context, bindings []models.Binding, gameID models.GameID) []error {
	var errors []error
	switch gameID {
	case models.GameArknights:
		for _, b := range bindings {
			if c.attendanceAPI.QueryAttendanceInfo(ctx, b, gameID) {
				log.Debug().Msgf("[%s] %s 今日已签到", b.ChannelName, b.NickName)
				continue
			}
			result, err := c.attendanceAPI.SignArknights(ctx, b)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			log.Info().Msgf("[%s] %s %s", b.ChannelName, b.NickName, utils.FormatArknightSignResult(result))
		}
	case models.GameEndfield:
		for _, b := range bindings {
			for _, role := range b.Roles {
				result, err := c.attendanceAPI.SignEndfield(ctx, role)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				log.Info().Msgf("[%s] %s %v", role.ServerName, role.Nickname, utils.FormatEndfieldSignResult(result))
			}
		}
	}
	return errors
}

func (c *CheckinService) RunCheckinTasks(ctx context.Context) bool {
	allTasksCompleted := true

	statusMap, err := c.checkinAPI.GetAllCheckinStatus(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("无法获取检票状态，将逐个检查")
		statusMap = nil
	}

	for gameID, gameName := range models.GameIDNames {
		logger := log.With().Int("gameId", int(gameID)).Str("gameName", gameName).Logger()

		if statusMap != nil {
			checked, ok := statusMap[int(gameID)]
			if ok && checked {
				logger.Debug().Msg("已检票")
				continue
			}
		}

		resp, err := c.checkinAPI.Checkin(ctx, gameID)
		if err != nil {
			allTasksCompleted = false
			logger.Error().Err(err).Msg("检票失败")
			continue
		}

		switch resp.Code {
		case 0:
			logger.Info().Msg("检票成功")
		case 10001:
			logger.Info().Msg("今日已检票")
		default:
			allTasksCompleted = false
			logger.Error().Int("code", resp.Code).Str("message", resp.Message).Msg("检票失败")
		}
	}
	return allTasksCompleted
}

func (c *CheckinService) getPopucomAchievement(ctx context.Context, bindings models.BindingResult, userID string) {
	for _, data := range bindings.Data.List {
		if data.AppCode == "popucom" {
			c.popucomAPI.GetPopucomAchievement(ctx, data.BindingList, userID)
		}
	}
}

func (c *CheckinService) getExaAchievement(ctx context.Context, bindings models.BindingResult, userID string) {
	for _, data := range bindings.Data.List {
		if data.AppCode == "exa" {
			c.exaAPI.GetExastrisAchievement(ctx, data.BindingList, userID)
		}
	}
}
