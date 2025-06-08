package api

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/rs/zerolog/log"
)

type PlayerAPI struct {
	client client.SklandHttpClient
}

func NewPlayerAPI(client client.SklandHttpClient) *PlayerAPI {
	return &PlayerAPI{client: client}
}

func (p *PlayerAPI) PrintAllPlayersInfo(bindings []models.Binding) error {
	logger := log.With().Logger()
	var playerData models.PlayerResponse
	for _, binding := range bindings {
		b := binding
		respBody, err := p.fetchPlayerInfo(b)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(respBody, &playerData); err != nil {
			return err
		}
		//activityData := playerData.Data.ActivityInfoMap
		//for _, key := range activityData.Keys() {
		//	value, exists := activityData.Get(key)
		//	if !exists {
		//		continue // 或处理不存在的情况
		//	}
		//	fmt.Printf("Key: %s, Value: %v\n", key, value)
		//}
		logger.Info().
			Str("uid", b.Uid).
			Str("nickname", b.NickName).
			Msg("成功获取玩家详细信息")
		fmt.Printf("=== 玩家 %s (%s) ===\n", b.NickName, b.Uid)
	}
	return nil
}

func (p *PlayerAPI) fetchPlayerInfo(b models.Binding) ([]byte, error) {
	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", b.Uid)
	headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("获取签名头失败: %w", err)
	}
	headers["Content-Type"] = "application/json"

	resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer client.CloseResponse(resp)

	body, err := client.ReadResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var errResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, fmt.Errorf("解析响应元数据失败: %w", err)
	}

	if errResp.Code != 0 {
		return nil, fmt.Errorf("服务端返回错误: %s (code: %d)", errResp.Message, errResp.Code)
	}
	return body, nil
}
