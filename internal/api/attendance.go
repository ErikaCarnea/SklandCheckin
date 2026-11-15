package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
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
	reqBody := models.AttendanceRequest{Uid: b.Uid, ChannelMasterId: b.ChannelMasterId}
	result, err := client.ExecuteRequest[models.AttendanceResult](ctx, a.client, http.MethodPost, AttendanceURL, reqBody, client.WithSign())
	if err != nil {
		return nil, fmt.Errorf("%v | %w", b.ToString(), err)
	}

	switch result.GetCode() {
	case 0:
		return result, nil
	case 10001:
		if strings.Contains(result.GetMessage(), "重复") {
			return result, nil
		}
	}
	return nil, fmt.Errorf("%v | %s (code: %d)", b.ToString(), result.GetMessage(), result.GetCode())
}

func (a *AttendanceAPI) HasArknightsSignedToday(ctx context.Context, b models.Binding) (bool, error) {
	attendanceInfo, err := a.getAttendanceInfo(ctx, b)
	if err != nil {
		return false, err
	}

	if attendanceInfo.GetCode() != 0 {
		return false, fmt.Errorf("%v | %s (code: %d)", b.ToString(), attendanceInfo.GetMessage(), attendanceInfo.GetCode())
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
			return true, nil
		}
	}
	return false, nil
}

func (a *AttendanceAPI) getAttendanceInfo(ctx context.Context, binding models.Binding) (*models.AttendanceInfo, error) {
	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/attendance?uid=%s&gameId=%d", binding.Uid, models.GameArknights)
	attendanceInfo, err := client.ExecuteRequest[models.AttendanceInfo](ctx, a.client, http.MethodGet, urlStr, nil, client.WithSign())
	if err != nil {
		return nil, fmt.Errorf("查询签到信息失败: %w", err)
	}
	return attendanceInfo, nil
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
	case 0:
		return result, nil
	case 10001:
		if strings.Contains(result.GetMessage(), "重复") {
			return result, nil
		}
	}
	return nil, fmt.Errorf("%v | %s (code: %d)", role.ToString(), result.GetMessage(), result.GetCode())
}

func (a *AttendanceAPI) HasEndfieldSignedToday(ctx context.Context, role models.Role) (bool, error) {
	skGameRole := fmt.Sprintf("3_%s_%s", role.RoleId, role.ServerId)
	result, err := client.ExecuteRequest[models.EndfieldAttendanceStatus](ctx, a.client, http.MethodGet, EndfieldSignURL, nil,
		client.WithSign(),
		client.WithHeader("sk-game-role", skGameRole),
	)
	if err != nil {
		return false, fmt.Errorf("%v | %w", role.ToString(), err)
	}

	if result.GetCode() != 0 {
		return false, fmt.Errorf("%v | %s (code: %d)", role.ToString(), result.GetMessage(), result.GetCode())
	}
	return result.Data.HasToday, nil
}
