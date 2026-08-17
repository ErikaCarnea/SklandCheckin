package app

import (
	"context"
	"os"

	"github.com/ErikaCarnea/Skland/internal/api"
	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/utils"
	"github.com/rs/zerolog/log"
)

// App 应用入口，仅负责依赖组装和流程编排。
type App struct {
	httpClient client.SklandHTTPClient
	auth       *AuthService
	checkin    *CheckinService
}

func NewApp() *App {
	httpClient := client.NewClient()
	return &App{
		httpClient: httpClient,
		auth:       NewAuthService(api.NewAuthAPI(httpClient)),
		checkin: NewCheckinService(
			api.NewBindingAPI(httpClient),
			api.NewAttendanceAPI(httpClient),
			api.NewCheckinAPI(httpClient),
			api.NewPlayerAPI(httpClient),
			api.NewPopucomAPI(httpClient),
			api.NewExaAPI(httpClient),
		),
	}
}

func (a *App) Run(ctx context.Context) {
	credResult := a.auth.Authenticate(ctx)
	if credResult == nil {
		log.Error().Msg("登录失败，程序退出")
		// CI 环境下以非零码退出，让 workflow 标记为失败并触发通知
		if os.Getenv("CI") == "true" {
			os.Exit(1)
		}
		return
	}

	a.httpClient.SetCred(credResult.Data.Cred)
	a.httpClient.SetSignToken(credResult.Data.Token)

	a.checkin.RunAll(ctx, credResult)

	utils.WaitForExit()
}
