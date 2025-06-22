package api

import (
	"net/http"

	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
)

type BindingAPI struct {
	client client.SklandHTTPClient
}

func NewBindingAPI(c client.SklandHTTPClient) *BindingAPI {
	return &BindingAPI{client: c}
}

func (b *BindingAPI) GetBindingList() (models.BindingResult, error) {
	var result models.BindingResult
	if err := b.client.ExecuteRequest(
		http.MethodGet,
		"https://zonai.skland.com/api/v1/game/player/binding",
		nil,
		&result,
	); err != nil {
		return models.BindingResult{}, err
	}

	// var bindings []models.Binding
	// for _, app := range result.Data.List {
	// 	bindings = append(bindings, app.BindingList...)
	// }
	// return bindings, nil
	return result, nil
}
