package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
	"github.com/rs/zerolog/log"
)

const AttendanceURL = "https://zonai.skland.com/api/v1/game/attendance"
const EndfieldSignURL = "https://zonai.skland.com/web/v1/game/endfield/attendance"

type AttendanceAPI struct {
	client client.SklandHTTPClient
}

func NewAttendanceAPI(client client.SklandHTTPClient) *AttendanceAPI {
	return &AttendanceAPI{client: client}
}

func (a *AttendanceAPI) SignArknights(ctx context.Context, b models.Binding) (*models.AttendanceResult, error) {
	reqBody := models.AttendanceRequest{Uid: b.Uid, GameId: b.ChannelMasterId}
	result, err := client.ExecuteRequest[models.AttendanceResult](ctx, a.client, http.MethodPost, AttendanceURL, reqBody, client.WithSign())
	if err != nil {
		return nil, fmt.Errorf("%v | %w", b.ToString(), err)
	}
	return result, nil
}

func (a *AttendanceAPI) QueryAttendanceInfo(ctx context.Context, b models.Binding, gameID models.GameID) bool {
	attendanceInfo, err := a.getAttendanceInfo(ctx, b, gameID)
	if err != nil {
		log.Error().Err(err).Msg("查询签到信息失败")
		return false
	}

	if attendanceInfo.GetCode() != 0 {
		log.Error().Int("code", attendanceInfo.GetCode()).Msg("API返回错误")
		return false
	}

	cstZone := time.FixedZone("Asia/Shanghai", 8*60*60)
	now := time.Now().In(cstZone)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, cstZone)
	todayUnix := today.Unix()

	for _, record := range attendanceInfo.Data.Records {
		ts, err := strconv.ParseInt(record.Ts, 10, 64)
		if err != nil {
			continue
		}
		if ts == todayUnix {
			return true
		}
	}
	return false
}

func (a *AttendanceAPI) getAttendanceInfo(ctx context.Context, binding models.Binding, gameID models.GameID) (*models.AttendanceInfo, error) {
	switch gameID {
	case models.GameArknights:
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/attendance?uid=%s&gameId=%d", binding.Uid, gameID)
		attendanceInfo, err := client.ExecuteRequest[models.AttendanceInfo](ctx, a.client, http.MethodGet, urlStr, nil, client.WithSign())
		if err != nil {
			return nil, fmt.Errorf("查询签到信息失败: %w", err)
		}
		return attendanceInfo, nil
	default:
		return nil, fmt.Errorf("不支持的游戏ID: %s %d", gameID.String(), gameID)
	}
}

func (a *AttendanceAPI) SignEndfield(ctx context.Context, role models.Role) (*models.EndfieldResult, error) {
	skGameRole := fmt.Sprintf("3_%s_%s", role.RoleId, role.ServerId)
	result, err := client.ExecuteRequest[models.EndfieldResult](ctx, a.client, http.MethodPost, EndfieldSignURL, nil,
		client.WithSign(),
		client.WithHeader("sk-game-role", skGameRole),
		client.WithHeader("referer", "https://game.skland.com/"),
		client.WithHeader("origin", "https://game.skland.com/"),
		client.WithHeader("Content-Type", "application/json"),
	)
	if err != nil {
		return nil, fmt.Errorf("%v | %w", role.ToString(), err)
	}

	switch result.GetCode() {
	case 0, 10001:
		return result, nil
	default:
		return nil, fmt.Errorf("%v | %s", role.ToString(), result.GetMessage())
	}
}
