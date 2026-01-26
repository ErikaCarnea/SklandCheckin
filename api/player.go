package api

import (
	"fmt"

	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/ErikaCarnea/Skland/models/player"
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
	var playerData player.PlayerResponse

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

		// if binding.Uid == "64847751" {
		// 	for item, activityInfo := range playerData.Data.ActivityInfoMap {
		// 		for _, actInfo := range playerData.Data.Activity {
		// 			if item == actInfo.ActId {
		// 				for _, zone := range actInfo.Zones {
		// 					if zone.ClearedStage != zone.TotalStage {
		// 						fmt.Printf("%s\n", activityInfo.Name)
		// 						for _, v := range actInfo.Zones {
		// 							fmt.Printf("%s\n", v.ZoneId)
		// 							fmt.Printf("%d/%d\n", v.ClearedStage, v.TotalStage)
		// 						}
		// 						fmt.Println("=======================================")
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// }
	}
	return nil
}
