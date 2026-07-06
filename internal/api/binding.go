package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
)

const BindingURL = "https://zonai.skland.com/api/v1/game/player/binding"

type BindingAPI struct {
	client client.SklandHTTPClient
}

func NewBindingAPI(client client.SklandHTTPClient) *BindingAPI {
	return &BindingAPI{client: client}
}

func (b *BindingAPI) GetBindingList(ctx context.Context) (models.BindingResult, error) {
	var result models.BindingResult
	if err := b.client.ExecuteRequest(ctx,
		http.MethodGet,
		BindingURL,
		nil,
		&result,
		client.WithSign(),
	); err != nil {
		return models.BindingResult{}, err
	}

	if result.GetCode() != 0 {
		return models.BindingResult{}, fmt.Errorf("API错误: %s (code: %d)", result.Message, result.Code)
	}
	return result, nil
}
