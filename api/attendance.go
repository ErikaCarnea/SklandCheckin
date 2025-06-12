package api

import (
	"encoding/json"
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

func NewAttendanceAPI(c client.SklandHTTPClient) *AttendanceAPI {
	return &AttendanceAPI{client: c}
}

func (a *AttendanceAPI) SignAttendance(b models.Binding) (*models.AttendanceResult, error) {
	reqBody := models.AttendanceRequest{Uid: b.Uid, GameId: b.ChannelMasterId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		AttendanceURL,
		reqBody,
		&result,
	); err != nil {
		return nil, fmt.Errorf("%v | %w",
			b.ToString(),
			err)
	}
	return &result, nil
}

func (a *AttendanceAPI) QueryAttendanceInfo(b models.Binding) bool {
	var attendanceInfo models.AttendanceInfo
	body, err := a.queryAttendance(b)
	if err != nil {
		log.Error().Err(err).Msg("查询签到信息失败")
	}

	if err := json.Unmarshal(body, &attendanceInfo); err != nil {
		log.Error().Err(err).Msg("解析签到信息失败")
	}

	if attendanceInfo.Code == 0 && attendanceInfo.Message == "OK" {
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
	}
	return false
}

func (a *AttendanceAPI) queryAttendance(binding models.Binding) ([]byte, error) {
	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/attendance?uid=%s&gameId=1", binding.Uid)
	headers, err := a.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	if err != nil {

		return nil, fmt.Errorf("获取签名头失败: %w", err)
	}
	headers["Content-Type"] = "application/json"

	resp, err := a.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer client.CloseResponse(resp)

	body, err := client.ReadResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}
