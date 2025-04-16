package api

import (
	"Skland/client"
	"Skland/models"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
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
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":    b.Uid,
		"player": b.NickName,
	})

	// 检查上下文是否已取消
	if err := ctx.Err(); err != nil {
		logEntry.WithError(err).Debug("上下文已取消，跳过请求")
		return nil
	}

	urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", b.Uid)
	headers, err := p.client.GetSignHeaders(urlStr, http.MethodGet, nil)
	if err != nil {
		logEntry.WithError(err).Error("获取签名头失败")
		return err
	}
	headers["Content-Type"] = "application/json"

	resp, err := p.client.DoRequest(http.MethodGet, urlStr, nil, headers)
	if err != nil {
		logEntry.WithError(err).Error("请求失败")
		return err
	}
	defer client.CloseResponse(resp)

	body, err := client.ReadResponseBody(resp)
	if err != nil {
		logEntry.WithError(err).Error("读取响应失败")
		return err
	}

	logEntry.Info("成功获取玩家信息")
	fmt.Printf("=== 玩家 %s (%s) ===\n", b.NickName, b.Uid)
	fmt.Println(string(body))
	return nil
}
