package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
)

type ExaAPI struct {
	client client.SklandHTTPClient
}

func NewExaAPI(client client.SklandHTTPClient) *ExaAPI {
	return &ExaAPI{client: client}
}

// GetExastrisAchievement 查询来自星辰游戏成就。
func (p *ExaAPI) GetExastrisAchievement(ctx context.Context, b []models.Binding, userID string) error {
	var result models.ExaResult

	for _, user := range b {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/exastris?uid=%s&userId=%s", user.Uid, userID)

		if err := p.client.ExecuteRequest(ctx,
			http.MethodGet,
			urlStr,
			nil,
			&result,
			client.WithSign(),
		); err != nil {
			return fmt.Errorf("用户[%s]查询成就失败: %w", user.Uid, err)
		}
	}
	return nil
}
