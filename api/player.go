package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
)

type PlayerApi struct {
	client *client.HttpClient
}

func NewPlayerAPI(client *client.HttpClient) *PlayerApi {
	return &PlayerApi{client: client}
}

func (p *PlayerApi) PrintAllPlayersInfo(bindings []models.Binding) error {
	for _, binding := range bindings {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", binding.Uid)

		headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
		if err != nil {
			return fmt.Errorf("获取签名头失败(UID:%s): %w", binding.Uid, err)
		}
		headers["Content-Type"] = "application/json"

		resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
		if err != nil {
			return err
		}

		func() {
			defer func() {
				if closeErr := resp.Body.Close(); closeErr != nil {
					fmt.Printf("警告：关闭响应体失败(UID:%s): %v\n", binding.Uid, closeErr)
				}
			}()

			body, err := client.ReadResponseBody(resp)
			if err != nil {
				fmt.Printf("读取响应失败: %v\n", err)
				return
			}

			content := string(body)
			fmt.Printf("=== 玩家 %s (%s) ===\n", binding.NickName, binding.Uid)
			fmt.Println(content)
			fmt.Println("========================")

		}()

	}
	return nil
}
