package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
	"github.com/ErikaCarnea/Skland/internal/models/player"
	"github.com/rs/zerolog/log"
)

type PlayerAPI struct {
	client client.SklandHTTPClient
}

func NewPlayerAPI(client client.SklandHTTPClient) *PlayerAPI {
	return &PlayerAPI{client: client}
}

// PrintAllPlayersInfo 查询并输出玩家详细信息。
// 输出委托给调用方 — 通过 zerolog 记录结构化日志，避免 fmt.Printf 混入 API 层。
func (p *PlayerAPI) PrintAllPlayersInfo(ctx context.Context, bindings []models.Binding) error {
	for _, binding := range bindings {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", binding.Uid)

		if _, err := client.ExecuteRequest[player.PlayerResponse](ctx, p.client, http.MethodGet, urlStr, nil, client.WithSign()); err != nil {
			return fmt.Errorf("获取玩家[%s]信息失败: %w", binding.Uid, err)
		}

		log.Info().
			Str("uid", binding.Uid).
			Str("nickname", binding.NickName).
			Msg("成功获取玩家详细信息")
	}
	return nil
}
