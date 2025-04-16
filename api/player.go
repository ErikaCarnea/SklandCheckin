package api

import (
	"Skland/client"
	"Skland/models"
	"fmt"
	"net/http"
	"sync"
)

type PlayerApi struct {
	client client.HTTPClient
}

func NewPlayerAPI(client client.HTTPClient) *PlayerApi {
	return &PlayerApi{client: client}
}

func (p *PlayerApi) PrintAllPlayersInfo(bindings []models.Binding) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(bindings))

	for _, binding := range bindings {
		wg.Add(1)
		go func(b models.Binding) {
			defer wg.Done()

			//发起单个请求
			urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", b.Uid)

			headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
			if err != nil {
				errChan <- fmt.Errorf("获取签名头失败(UID:%s): %w", b.Uid, err)
				return
			}
			headers["Content-Type"] = "application/json"

			resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
			if err != nil {
				errChan <- err
				return
			}
			defer func() {
				if closeErr := resp.Body.Close(); closeErr != nil {
					errChan <- fmt.Errorf("关闭响应体失败(UID:%s): %v", b.Uid, closeErr)
				}
			}()

			body, err := client.ReadResponseBody(resp)
			if err != nil {
				errChan <- fmt.Errorf("读取响应失败(UID:%s): %v", b.Uid, err)
				return
			}

			content := string(body)
			fmt.Printf("=== 玩家 %s (%s) ===\n", b.NickName, b.Uid)
			fmt.Println(content)
			fmt.Println("========================")
		}(binding)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return fmt.Errorf("共 %d 个错误，首个错误: %v", len(errors), errors[0])
	}

	return nil

	//for _, binding := range bindings {
	//	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", binding.Uid)
	//
	//	headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	//	if err != nil {
	//		return fmt.Errorf("获取签名头失败(UID:%s): %w", binding.Uid, err)
	//	}
	//	headers["Content-Type"] = "application/json"
	//
	//	resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	//	if err != nil {
	//		return err
	//	}
	//
	//	func() {
	//		defer func() {
	//			if closeErr := resp.Body.Close(); closeErr != nil {
	//				fmt.Printf("警告：关闭响应体失败(UID:%s): %v\n", binding.Uid, closeErr)
	//			}
	//		}()
	//
	//		body, err := client.ReadResponseBody(resp)
	//		if err != nil {
	//			fmt.Printf("读取响应失败: %v\n", err)
	//			return
	//		}
	//
	//		content := string(body)
	//		fmt.Printf("=== 玩家 %s (%s) ===\n", binding.NickName, binding.Uid)
	//		fmt.Println(content)
	//		fmt.Println("========================")
	//
	//	}()
	//
	//}
	//return nil
}
