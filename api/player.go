package api

import (
	"fmt"

	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/rs/zerolog/log"
)

type PlayerAPI struct {
	client client.SklandHTTPClient
}

func NewPlayerAPI(client client.SklandHTTPClient) *PlayerAPI {
	return &PlayerAPI{client: client}
}

func (p *PlayerAPI) PrintAllPlayersInfo(bindings []models.Binding) error {
	logger := log.With().Logger()
	var playerData models.PlayerResponse

	for _, binding := range bindings {
		// 使用新的 ExecuteRequest 方法替换原有的 fetchPlayerInfo
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", binding.Uid)

		if err := p.client.ExecuteRequest(
			http.MethodGet,
			urlStr,
			nil,
			&playerData,
			client.SignedRequest,
		); err != nil {
			return fmt.Errorf("获取玩家[%s]信息失败: %w", binding.Uid, err)
		}

		logger.Info().
			Str("uid", binding.Uid).
			Str("nickname", binding.NickName).
			Msg("成功获取玩家详细信息")
		fmt.Printf("=== 玩家 %s (%s) ===\n", binding.NickName, binding.Uid)

		// 注释掉的代码保持不变
		// activityData := playerData.Data.ActivityInfoMap
		// for _, key := range activityData.Keys() {
		// 	value, exists := activityData.Get(key)
		// 	if !exists {
		// 		continue // 或处理不存在的情况
		// 	}
		// 	fmt.Printf("Key: %s, Value: %v\n", key, value)
		// }
	}
	return nil
}
