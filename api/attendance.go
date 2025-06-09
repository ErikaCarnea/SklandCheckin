package api

import (
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
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
