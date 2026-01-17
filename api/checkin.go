package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
)

const CheckinURL = "https://zonai.skland.com/api/v1/score/checkin"
const IsCheckinURL = "https://zonai.skland.com/api/v1/score/ischeckin"

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

	err := c.client.ExecuteRequest(
		http.MethodPost,
		CheckinURL,
		reqBody,
		&resp,
		client.SignedRequest,
	)

	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "code: 10001") || strings.Contains(errStr, "重复签到") {
			return &resp, nil
		}
		return nil, fmt.Errorf("签到失败: %w", err)
	}

	// 正常情况
	return &resp, nil
}

func (c *CheckinAPI) CheckIsCheckin(gameID int) (bool, error) {
	var result models.IsCheckinResponse
	if err := c.client.ExecuteRequest(
		http.MethodGet,
		IsCheckinURL,
		nil,
		&result,
		client.SignedRequest,
	); err != nil {
		return false, err
	}

	for _, item := range result.Data.List {
		if item.GameId == gameID {
			return item.Checked == 1, nil
		}
	}
	return false, nil
}

func (c *CheckinAPI) GetAllCheckinStatus() (map[int]bool, error) {
	var result models.IsCheckinResponse
	if err := c.client.ExecuteRequest(
		http.MethodGet,
		IsCheckinURL,
		nil,
		&result,
		client.SignedRequest,
	); err != nil {
		return nil, err
	}

	statusMap := make(map[int]bool)
	for _, item := range result.Data.List {
		statusMap[item.GameId] = item.Checked == 1
	}
	return statusMap, nil
}
