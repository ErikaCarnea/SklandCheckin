package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
)

type PopucomAPI struct {
	client client.SklandHTTPClient
}

func NewPopucomAPI(client client.SklandHTTPClient) *PopucomAPI {
	return &PopucomAPI{client: client}
}

func (p *PopucomAPI) GetPopucomAchievement(ctx context.Context, b []models.Binding, userID string) error {
	var result models.PopucomResult
	for _, user := range b {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/popucom?uid=%s&userId=%s", user.Uid, userID)

		if err := p.client.ExecuteRequest(ctx,
			http.MethodGet,
			urlStr,
			nil,
			&result,
			client.WithSign(),
		); err != nil {
			return fmt.Errorf("用户[%s]查询泡姆泡姆成就失败: %w", user.Uid, err)
		}
	}
	return nil
}
