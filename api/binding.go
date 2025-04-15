package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
)

type BindingAPI struct {
	client *client.HttpClient
}

func NewBindingAPI(c *client.HttpClient) *BindingAPI {
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
	//urlStr := "https://zonai.skland.com/api/v1/game/player/binding"
	//
	//headers, err := b.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	//if err != nil {
	//	return nil, err
	//}
	//
	//resp, err := b.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	//if err != nil {
	//	return nil, err
	//}
	//defer func(Body io.ReadCloser) {
	//	if err := Body.Close(); err != nil {
	//		log.Printf("关闭响应体失败: %v", err)
	//	}
	//}(resp.Body)
	//
	//body, err := client.ReadResponseBody(resp)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var result models.BindingResult
	//if err := json.Unmarshal(body, &result); err != nil {
	//	return nil, err
	//}
	//
	//if result.Code != 0 {
	//	return nil, fmt.Errorf("API error: code %d", result.Code)
	//}
	//
	//var bindings []models.Binding
	//for _, app := range result.Data.List {
	//	if app.AppCode == "arknights" {
	//		bindings = append(bindings, app.BindingList...)
	//	}
	//}
	//return bindings, nil
}
