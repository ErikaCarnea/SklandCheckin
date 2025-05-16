package api

import (
	"github.com/HeathErika/Skland/client"
	"github.com/HeathErika/Skland/models"
	"net/http"
	"strconv"
)

var SklandBoard = map[int]string{
	1:   "明日方舟",
	2:   "来自星辰",
	3:   "明日方舟: 终末地",
	4:   "泡姆泡姆",
	100: "纳斯特港",
	101: "开拓芯",
}

type CheckinAPI struct {
	client client.SklandHttpClient
}

func NewCheckinAPI(c client.SklandHttpClient) *CheckinAPI {
	return &CheckinAPI{client: c}
}

func (c *CheckinAPI) Checkin(gameID int) (*models.CheckinResponse, error) {
	reqBody := models.CheckinRequest{
		GameID: strconv.Itoa(gameID),
	}

	var resp models.CheckinResponse
	if err := c.client.ExecuteRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/score/checkin",
		reqBody,
		&resp,
	); err != nil {
		return nil, err
	}
	return &resp, nil
}
