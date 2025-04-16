package api

import (
	"Skland/client"
	"Skland/models"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"net/http"
)

type PlayerApi struct {
	client client.HTTPClient
}

func NewPlayerAPI(client client.HTTPClient) *PlayerApi {
	return &PlayerApi{client: client}
}

func (p *PlayerApi) PrintAllPlayersInfo(bindings []models.Binding) error {
	g, ctx := errgroup.WithContext(context.Background())
	sem := semaphore.NewWeighted(5)

	for _, binding := range bindings {
		b := binding
		if err := sem.Acquire(ctx, 1); err != nil {
			return fmt.Errorf("获取信号量失败: %w", err)
		}

		g.Go(func() error {
			defer sem.Release(1)
			return p.fetchPlayerInfo(ctx, b)
		})
	}

	// 等待所有 goroutine 完成,返回首个错误
	if err := g.Wait(); err != nil {
		return fmt.Errorf("获取玩家信息失败: %w", err)
	}
	return nil
}

func (p *PlayerApi) fetchPlayerInfo(ctx context.Context, b models.Binding) error {
	// 检查上下文是否已取消
	if err := ctx.Err(); err != nil {
		return err
	}

	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", b.Uid)
	headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("获取签名头失败(UID:%s): %w", b.Uid, err)
	}
	headers["Content-Type"] = "application/json"

	resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer client.CloseResponse(resp)

	body, err := client.ReadResponseBody(resp)
	if err != nil {
		return fmt.Errorf("读取响应失败(UID:%s): %v", b.Uid, err)
	}

	content := string(body)
	fmt.Printf("=== 玩家 %s (%s) ===\n", b.NickName, b.Uid)
	fmt.Println(content)
	fmt.Println("========================")
	return nil
}
