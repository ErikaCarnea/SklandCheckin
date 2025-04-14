package api

import (
	"Skland/client"
	"Skland/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BindingAPI struct {
	client *client.HttpClient
}

func NewBindingAPI(c *client.HttpClient) *BindingAPI {
	return &BindingAPI{client: c}
}

func (b *BindingAPI) GetBindingList() ([]models.Binding, error) {
	urlStr := "https://zonai.skland.com/api/v1/game/player/binding"

	headers, err := b.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := client.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var result models.BindingResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: code %d", result.Code)
	}

	var bindings []models.Binding
	for _, app := range result.Data.List {
		if app.AppCode == "arknights" {
			bindings = append(bindings, app.BindingList...)
		}
	}
	return bindings, nil
}
