package api

import (
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
)

const BindingURL = "https://zonai.skland.com/api/v1/game/player/binding"

type BindingAPI struct {
	client client.SklandHTTPClient
}

func NewBindingAPI(client client.SklandHTTPClient) *BindingAPI {
	return &BindingAPI{client: client}
}

func (b *BindingAPI) GetBindingList() (models.BindingResult, error) {
	var result models.BindingResult
	if err := b.client.ExecuteRequest(
		http.MethodGet,
		BindingURL,
		nil,
		&result,
		client.SignedRequest,
	); err != nil {
		return models.BindingResult{}, err
	}

	if result.GetCode() != 0 {
		return models.BindingResult{}, fmt.Errorf("API错误: %s (code: %d)", result.Message, result.Code)
	}
	// var bindings []models.Binding
	// for _, app := range result.Data.List {
	// 	bindings = append(bindings, app.BindingList...)
	// }
	// return bindings, nil
	return result, nil
}

// HasValidBindings 检查用户是否有有效的游戏绑定
// 返回true表示有绑定，false表示没有绑定或获取失败
func (b *BindingAPI) HasValidBindings() (bool, error) {
	result, err := b.GetBindingList()
	if err != nil {
		return false, err
	}
	return result.Code == 0 && len(result.Data.List) > 0, nil
}
