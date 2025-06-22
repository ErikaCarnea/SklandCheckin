package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
)

type ExaAPI struct {
	client client.SklandHTTPClient
}

func NewExaAPI(client client.SklandHTTPClient) *ExaAPI {
	return &ExaAPI{client: client}
}

func (p *ExaAPI) GetExastrisAchievement(b []models.Binding, userId string) error {
	// logger := log.With().Logger()
	var result models.ExaResult
	for _, user := range b {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/exastris?uid=%s&userId=%s", user.Uid, userId)

		headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
		if err != nil {
			return fmt.Errorf("获取签名头失败: %w", err)
		}
		headers["Content-Type"] = "application/json"

		resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
		if err != nil {
			return fmt.Errorf("请求失败: %w", err)
		}
		defer client.CloseResponse(resp)

		body, err := client.ReadResponseBody(resp)
		if err != nil {
			return fmt.Errorf("读取响应失败: %w", err)
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("解析响应失败: %w", err)
		}
		if result.Code != 0 {
			return fmt.Errorf("服务端返回错误: %s (code: %d)", result.Message, result.Code)
		}
	}
	return nil
}
