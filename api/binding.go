package api

import (
	"net/http"

	"github.com/HeathErika/Skland/client"
	"github.com/HeathErika/Skland/models"
)

type BindingAPI struct {
	client client.SklandHttpClient
}

func NewBindingAPI(c client.SklandHttpClient) *BindingAPI {
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
		return nil, err
	}

	var bindings []models.Binding
	for _, app := range result.Data.List {
		if app.AppCode == "arknights" {
			bindings = append(bindings, app.BindingList...)
		}
	}
	return bindings, nil
}
