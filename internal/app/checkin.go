package app

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/ErikaCarnea/Skland/internal/api"
	"github.com/ErikaCarnea/Skland/internal/client"
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
func (c *CheckinService) RunAll(ctx context.Context, credResult *models.CredResult) error {
	bindings, err := c.bindingAPI.GetBindingList(ctx)
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		return err
	}

	c.getPopucomAchievement(ctx, bindings, credResult.Data.UserId)
	c.getExaAchievement(ctx, bindings, credResult.Data.UserId)

	signErrors := c.PerformSignAttendance(ctx, bindings)
	checkinErr := c.RunCheckinTasks(ctx)

	for _, err := range signErrors {
		log.Error().Err(err).Msg("签到失败")
	}

	allErrors := signErrors
	if checkinErr != nil {
		allErrors = append(allErrors, checkinErr)
	}

	if len(allErrors) == 0 {
		log.Info().Msg("今日已签到")
	}

	for _, err := range allErrors {
		if errors.Is(err, client.ErrUnauthorized) {
			return err
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("签到流程存在 %d 个失败", len(allErrors))
	}
	return nil
}

func (c *CheckinService) PerformSignAttendance(ctx context.Context, bindings models.BindingResult) []error {
	var signErrors []error

	for _, data := range bindings.Data.List {
		switch data.AppCode {
		case "arknights":
			signErrors = append(signErrors, c.signAttendance(ctx, data.BindingList, models.GameArknights)...)
		case "endfield":
			signErrors = append(signErrors, c.signAttendance(ctx, data.BindingList, models.GameEndfield)...)
		}
	}

	return signErrors
}

func (c *CheckinService) signAttendance(ctx context.Context, bindings []models.Binding, gameID models.GameID) []error {
	var errors []error
	switch gameID {
	case models.GameArknights:
		for _, b := range bindings {
			signed, err := c.attendanceAPI.HasArknightsSignedToday(ctx, b)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			if signed {
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
				signed, err := c.attendanceAPI.HasEndfieldSignedToday(ctx, role)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				if signed {
					log.Debug().Msgf("[%s] %s 今日已签到", role.ServerName, role.Nickname)
					continue
				}
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

func (c *CheckinService) RunCheckinTasks(ctx context.Context) error {
	statusMap, err := c.checkinAPI.GetAllCheckinStatus(ctx)
	if err != nil {
		if errors.Is(err, client.ErrUnauthorized) {
			return err
		}
		log.Warn().Err(err).Msg("无法获取检票状态，将逐个检查")
		statusMap = nil
	}

	ids := make(map[int]struct{}, len(statusMap)+len(models.GameIDNames))
	for id := range statusMap {
		ids[id] = struct{}{}
	}
	for id := range models.GameIDNames {
		ids[int(id)] = struct{}{}
	}

	boardIDs := make([]int, 0, len(ids))
	for id := range ids {
		boardIDs = append(boardIDs, id)
	}
	sort.Ints(boardIDs)

	var succeeded, already, skipped []string
	order := make([]string, 0, len(boardIDs))
	groups := make(map[string][]string)

	for _, id := range boardIDs {
		name := boardName(id)
		if checked, ok := statusMap[id]; ok && checked {
			skipped = append(skipped, name)
			continue
		}

		resp, err := c.checkinAPI.Checkin(ctx, models.GameID(id))
		if err != nil {
			key := err.Error()
			if _, ok := groups[key]; !ok {
				order = append(order, key)
			}
			groups[key] = append(groups[key], name)
			continue
		}

		switch resp.Code {
		case 0:
			succeeded = append(succeeded, name)
		case 10001:
			already = append(already, name)
		}
	}

	if len(groups) > 0 {
		total := 0
		parts := make([]string, 0, len(order))
		for _, key := range order {
			names := groups[key]
			total += len(names)
			if len(names) == 1 {
				parts = append(parts, fmt.Sprintf("%s(%s)", names[0], key))
			} else {
				parts = append(parts, fmt.Sprintf("%s x%d", key, len(names)))
			}
		}
		err := fmt.Errorf("检票失败 %d 个板块: %s", total, strings.Join(parts, "; "))
		log.Error().Err(err).Msg("检票失败")
		return err
	}

	summary := []string{fmt.Sprintf("共 %d 个板块", len(boardIDs))}
	if len(succeeded) > 0 {
		summary = append(summary, fmt.Sprintf("新检 %d 个（%s）", len(succeeded), strings.Join(succeeded, "、")))
	}
	if len(already) > 0 {
		summary = append(summary, fmt.Sprintf("已检 %d 个（%s）", len(already), strings.Join(already, "、")))
	}
	if len(skipped) > 0 {
		summary = append(summary, fmt.Sprintf("跳过 %d 个（%s）", len(skipped), strings.Join(skipped, "、")))
	}
	log.Info().Msg("检票完成：" + strings.Join(summary, "，"))
	return nil
}

func boardName(id int) string {
	if name, ok := models.GameIDNames[models.GameID(id)]; ok {
		return name
	}
	return fmt.Sprintf("未知板块(%d)", id)
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
