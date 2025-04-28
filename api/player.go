package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type PlayerApi struct {
	client client.HTTPClient
}

func NewPlayerAPI(client client.HTTPClient) *PlayerApi {
	return &PlayerApi{client: client}
}

func (p *PlayerApi) PrintAllPlayersInfo(bindings []models.Binding) error {
	for _, binding := range bindings {
		b := binding
		if err := p.fetchPlayerInfo(b); err != nil {
			return fmt.Errorf("获取玩家信息失败: %w", err)
		}
	}
	return nil
}

func (p *PlayerApi) fetchPlayerInfo(b models.Binding) error {
	logger := log.With().Str("uid", b.Uid).Str("player", b.NickName).Logger()

	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", b.Uid)
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

	logger.Info().Msg("成功获取玩家信息")
	fmt.Printf("=== 玩家 %s (%s) ===\n", b.NickName, b.Uid)
	fmt.Println(string(body))
	return nil
}
