package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
)

type AttendanceAPI struct {
	client client.HTTPClient
}

func NewAttendanceAPI(c client.HTTPClient) *AttendanceAPI {
	return &AttendanceAPI{client: c}
}

func (a *AttendanceAPI) SignAttendance(uid, gameId string) (*models.AttendanceResult, error) {
	reqBody := models.AttendanceRequest{Uid: uid, GameId: gameId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/game/attendance",
		reqBody,
		&result,
	); err != nil {
		return nil, fmt.Errorf("签到失败: %w", err)
	}

	return &result, nil
}
