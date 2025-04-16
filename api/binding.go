package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
)

type BindingAPI struct {
	client client.HTTPClient
}

func NewBindingAPI(c client.HTTPClient) *BindingAPI {
	return &BindingAPI{client: c}
}

func (b *BindingAPI) GetBindingList() ([]models.Binding, error) {
	var result models.BindingResult
	if err := b.client.ExecuteRequest(
		http.MethodGet,
		"https://zonai.skland.com/api/v1/game/player/binding",
		nil,
		&result,
	); err != nil {
		return nil, fmt.Errorf("获取绑定列表失败: %w", err)
	}

	var bindings []models.Binding
	for _, app := range result.Data.List {
		if app.AppCode == "arknights" {
			bindings = append(bindings, app.BindingList...)
		}
	}
	return bindings, nil
}
