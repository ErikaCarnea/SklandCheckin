package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
)

const CheckinURL = "https://zonai.skland.com/api/v1/score/checkin"

var SklandBoard = map[int]string{
	1:   "明日方舟",
	2:   "来自星辰",
	3:   "明日方舟: 终末地",
	4:   "泡姆泡姆",
	100: "纳斯特港",
	101: "开拓芯",
}

type CheckinAPI struct {
	client client.SklandHTTPClient
}

func NewCheckinAPI(client client.SklandHTTPClient) *CheckinAPI {
	return &CheckinAPI{client: client}
}

// Checkin 执行森空岛板块签到
// gameID: 游戏ID，对应 SklandBoard 中的键值
// 返回签到结果和错误信息
func (c *CheckinAPI) Checkin(gameID int) (*models.CheckinResponse, error) {
	reqBody := models.CheckinRequest{
		GameID: strconv.Itoa(gameID),
	}

	var resp models.CheckinResponse

	if err := c.client.ExecuteRequest(
		http.MethodPost,
		CheckinURL,
		reqBody,
		&resp,
		client.SignedRequest,
	); err != nil {
		return nil, fmt.Errorf("签到失败: %w", err)
	}

	switch resp.GetCode() {
	case 0:
		return &resp, nil
	case 10001:
		return &resp, nil
	default:
		return nil, fmt.Errorf("签到失败: %s (code: %d)", resp.Message, resp.Code)
	}
}
