package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
)

const CheckinURL = "https://zonai.skland.com/api/v1/score/checkin"
const IsCheckinURL = "https://zonai.skland.com/api/v1/score/ischeckin"

type CheckinAPI struct {
	client client.SklandHTTPClient
}

func NewCheckinAPI(client client.SklandHTTPClient) *CheckinAPI {
	return &CheckinAPI{client: client}
}

// Checkin 执行森空岛板块签到。gameID 使用 models.GameID 枚举。
func (c *CheckinAPI) Checkin(ctx context.Context, gameID models.GameID) (*models.CheckinResponse, error) {
	reqBody := models.CheckinRequest{
		GameID: strconv.Itoa(int(gameID)),
	}

	resp, err := client.ExecuteRequest[models.CheckinResponse](ctx, c.client, http.MethodPost, CheckinURL, reqBody, client.WithSign())

	if err != nil {
		return nil, fmt.Errorf("签到失败: %w", err)
	}

	// code: 10001 表示已签到，不是错误
	if resp.Code != 0 && resp.Code != 10001 {
		return nil, fmt.Errorf("签到失败: %s (code: %d)", resp.Message, resp.Code)
	}

	return resp, nil
}

func (c *CheckinAPI) GetAllCheckinStatus(ctx context.Context) (map[int]bool, error) {
	result, err := client.ExecuteRequest[models.IsCheckinResponse](ctx, c.client, http.MethodGet, IsCheckinURL, nil, client.WithSign())
	if err != nil {
		return nil, err
	}

	statusMap := make(map[int]bool)
	for _, item := range result.Data.List {
		statusMap[item.GameId] = item.Checked == 1
	}
	return statusMap, nil
}
