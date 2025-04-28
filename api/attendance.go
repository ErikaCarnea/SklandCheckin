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

func (a *AttendanceAPI) SignAttendance(b models.Binding) (*models.AttendanceResult, error) {
	//logEntry := logrus.WithFields(logrus.Fields{
	//	"account": b,
	//})

	reqBody := models.AttendanceRequest{Uid: b.Uid, GameId: b.ChannelMasterId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/game/attendance",
		reqBody,
		&result,
	); err != nil {
		//logEntry.WithError(err).Error("签到失败")
		return nil, fmt.Errorf("account: %v | %w",
			b.ToString(),
			err)
	}

	//logEntry.Info("签到成功")
	return &result, nil
}
