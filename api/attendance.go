package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
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

func (a *AttendanceAPI) SignArknights(b models.Binding) (*models.AttendanceResult, error) {
	reqBody := models.AttendanceRequest{Uid: b.Uid, GameId: b.ChannelMasterId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		AttendanceURL,
		reqBody,
		&result,
		client.SignedRequest,
	); err != nil {
		return nil, fmt.Errorf("%v | %w",
			b.ToString(b.GameId),
			err)
	}
	return &result, nil
}

func (a *AttendanceAPI) QueryAttendanceInfo(b models.Binding, gameID int) bool {
	attendanceInfo, err := a.getAttendanceInfo(b, gameID)
	if err != nil {
		log.Error().Err(err).Msg("查询签到信息失败")
		return false
	}

	if attendanceInfo.GetCode() != 0 {
		log.Error().Int("code", attendanceInfo.GetCode()).Msg("API返回错误")
		return false
	}

	cstZone := time.FixedZone("CST", 8*60*60)
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

func (a *AttendanceAPI) getAttendanceInfo(binding models.Binding, gameID int) (*models.AttendanceInfo, error) {
	var urlStr string
	var attendanceInfo models.AttendanceInfo
	switch gameID {
	case 1:
		urlStr = fmt.Sprintf("https://zonai.skland.com/api/v1/game/attendance?uid=%s&gameId=%d", binding.Uid, gameID)
		if err := a.client.ExecuteRequest(
			http.MethodGet,
			urlStr,
			nil,
			&attendanceInfo,
			client.SignedRequest,
		); err != nil {
			return nil, fmt.Errorf("查询签到信息失败: %w", err)
		}
	case 3:
		urlStr = fmt.Sprintf("https://zonai.skland.com/web/v1/game/endfield/attendance?uid=%s&gameId=%d", binding.Roles[0].RoleId, gameID)
	default:
		return nil, fmt.Errorf("游戏ID错误: %s %d", SklandBoard[gameID], gameID)
	}
	return &attendanceInfo, nil
}

func (a *AttendanceAPI) SignEndfield(b models.Binding) (*models.EndfieldResult, error) {
	var result models.EndfieldResult

	for _, role := range b.Roles {
		skGameRole := fmt.Sprintf("3_%s_%s", role.RoleId, role.ServerId)

		opts := client.SignedRequest
		if opts.Headers == nil {
			opts.Headers = make(map[string]string)
		}
		// 添加sk-game-role头
		opts.Headers["sk-game-role"] = skGameRole
		// 还可以添加其他Endfield需要的头，比如referer和origin
		opts.Headers["referer"] = "https://game.skland.com/"
		opts.Headers["origin"] = "https://game.skland.com/"
		opts.Headers["Content-Type"] = "application/json"
		if err := a.client.ExecuteRequest(
			http.MethodPost,
			EndfieldSignURL, // 使用正确的URL
			nil,             // 请求体为空
			&result,
			opts, // 使用包含sk-game-role的选项
		); err != nil {
			return nil, fmt.Errorf("%v | %w", b.ToString(b.GameId), err)
		}
		if result.GetCode() == 10001 {
			return &result, nil
		}
		if result.GetCode() == 10000 {
			return nil, fmt.Errorf("%v | %s", b.ToString(b.GameId), result.Message)
		} else if result.GetCode() == 10002 {
			return nil, fmt.Errorf("%v | %s", b.ToString(b.GameId), result.Message)
		} else if result.GetCode() != 0 {
			return nil, fmt.Errorf("%v | %s", b.ToString(b.GameId), result.Message)
		}
	}
	return &result, nil
}
