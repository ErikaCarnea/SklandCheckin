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

type AttendanceAPI struct {
	client client.SklandHTTPClient
}

func NewAttendanceAPI(client client.SklandHTTPClient) *AttendanceAPI {
	return &AttendanceAPI{client: client}
}

func (a *AttendanceAPI) SignAttendance(b models.Binding) (*models.AttendanceResult, error) {
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
			b.ToString(),
			err)
	}
	return &result, nil
}

func (a *AttendanceAPI) QueryAttendanceInfo(b models.Binding) bool {
	attendanceInfo, err := a.getAttendanceInfo(b)
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

func (a *AttendanceAPI) getAttendanceInfo(binding models.Binding) (*models.AttendanceInfo, error) {
	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/attendance?uid=%s&gameId=1", binding.Uid)

	var attendanceInfo models.AttendanceInfo
	if err := a.client.ExecuteRequest(
		http.MethodGet,
		urlStr,
		nil,
		&attendanceInfo,
		client.SignedRequest,
	); err != nil {
		return nil, fmt.Errorf("查询签到信息失败: %w", err)
	}

	return &attendanceInfo, nil
}
