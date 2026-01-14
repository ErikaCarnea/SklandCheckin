package api

import (
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

// GetExastrisAchievement 查询来自星辰游戏成就
// b: 用户绑定信息列表
// userId: 用户ID
// 返回错误信息
func (p *ExaAPI) GetExastrisAchievement(b []models.Binding, userId string) error {
	var result models.ExaResult

	for _, user := range b {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/exastris?uid=%s&userId=%s", user.Uid, userId)

		if err := p.client.ExecuteRequest(
			http.MethodGet,
			urlStr,
			nil,
			&result,
			client.SignedRequest,
		); err != nil {
			return fmt.Errorf("用户[%s]查询成就失败: %w", user.Uid, err)
		}
	}
	return nil
}
