package api

import (
	"Skland/client"
	"Skland/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AttendanceAPI struct {
	client client.HTTPClient
}

func NewAttendanceAPI(c client.HTTPClient) *AttendanceAPI {
	return &AttendanceAPI{client: c}
}

func (a *AttendanceAPI) SignAttendance(uid, gameId string) (*models.AttendanceResult, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":    uid,
		"gameId": gameId,
	})

	reqBody := models.AttendanceRequest{Uid: uid, GameId: gameId}
	var result models.AttendanceResult

	if err := a.client.ExecuteRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/game/attendance",
		reqBody,
		&result,
	); err != nil {
		logEntry.WithError(err).Error("签到失败")
		return nil, err
	}

	logEntry.Info("签到成功")
	return &result, nil
}
