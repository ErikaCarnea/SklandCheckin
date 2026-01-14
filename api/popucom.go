package api

import (
	"fmt"

	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
	_ "github.com/rs/zerolog/log"
)

type PopucomAPI struct {
	client client.SklandHTTPClient
}

func NewPopucomAPI(client client.SklandHTTPClient) *PopucomAPI {
	return &PopucomAPI{client: client}
}

func (p *PopucomAPI) GetPopucomAchievement(b []models.Binding, userId string) error {

	var result models.PopucomResult
	for _, user := range b {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/popucom?uid=%s&userId=%s", user.Uid, userId)

		if err := p.client.ExecuteRequest(
			http.MethodGet,
			urlStr,
			nil,
			&result,
			client.SignedRequest,
		); err != nil {
			return fmt.Errorf("用户[%s]查询泡姆泡姆成就失败: %w", user.Uid, err)
		}
	}
	return nil
}
