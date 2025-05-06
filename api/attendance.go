package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
)

type AttendanceAPI struct {
	client client.SklandHttpClient
}

func NewAttendanceAPI(c client.SklandHttpClient) *AttendanceAPI {
	return &AttendanceAPI{client: c}
}

func (a *AttendanceAPI) SignAttendance(b models.Binding) (*models.AttendanceResult, error) {
	reqBody := models.AttendanceRequest{Uid: b.Uid, GameId: b.ChannelMasterId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/game/attendance",
		reqBody,
		&result,
	); err != nil {
		return nil, fmt.Errorf("account: %v | %w",
			b.ToString(),
			err)
	}
	return &result, nil
}
