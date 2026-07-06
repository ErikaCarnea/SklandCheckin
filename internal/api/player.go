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

// GetAllPlayersInfo 查询并输出各绑定的玩家基本信息。
func (p *PlayerAPI) GetAllPlayersInfo(ctx context.Context, bindings []models.Binding) error {
	for _, binding := range bindings {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", binding.Uid)

		resp, err := client.ExecuteRequest[player.PlayerResponse](ctx, p.client, http.MethodGet, urlStr, nil, client.WithSign())
		if err != nil {
			return fmt.Errorf("获取玩家[%s]信息失败: %w", binding.Uid, err)
		}

		log.Info().
			Str("uid", binding.Uid).
			Str("nickname", resp.Data.Status.Name).
			Int("level", resp.Data.Status.Level).
			Str("mainStageProgress", resp.Data.Status.MainStageProgress).
			Int("charCnt", resp.Data.Status.CharCnt).
			Int("skinCnt", resp.Data.Status.SkinCnt).
			Msg("成功获取玩家信息")
	}
	return nil
}
