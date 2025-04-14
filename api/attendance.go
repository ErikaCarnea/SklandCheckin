package api

import (
	"Skland/client"
	"Skland/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type AttendanceAPI struct {
	client *client.HttpClient
}

func NewAttendanceAPI(c *client.HttpClient) *AttendanceAPI {
	return &AttendanceAPI{client: c}
}

func (a *AttendanceAPI) SignAttendance(uid, gameId string) (*models.AttendanceResult, error) {
	urlStr := "https://zonai.skland.com/api/v1/game/attendance"
	reqBody := models.AttendanceRequest{
		Uid:    uid,
		GameId: gameId,
	}

	headers, err := a.client.GetSignHeaders(urlStr, http.MethodPost, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get sign header: %v", err)
	}
	headers["Content-Type"] = "application/json"

	resp, err := a.client.DoRequest(http.MethodPost, urlStr, reqBody, headers)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := client.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var result models.AttendanceResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return &result, fmt.Errorf("API error: %s (code: %d)", result.Message, result.Code)
	}

	return &result, nil
}

func FormatAttendanceResult(result *models.AttendanceResult) string {
	ts, _ := strconv.ParseInt(result.Data.Timestamp, 10, 64)
	localTime := time.Unix(ts, 0).Local().Format("2006-01-02 15:04:05")

	var awards string
	for _, award := range result.Data.Awards {
		awards += fmt.Sprintf("获得奖励：%s x%d (类型：%s)\n",
			award.Resource.Name,
			award.Count,
			award.Type)
	}

	return fmt.Sprintf("[%s] 签到成功！\n%s", localTime, awards)
}
