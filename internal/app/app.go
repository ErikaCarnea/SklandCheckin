package app

import (
	"context"
	"errors"
	"os"

	"github.com/ErikaCarnea/Skland/internal/api"
	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/models"
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
		exitIfCI()
		return
	}

	a.setCred(credResult)

	if err := a.auth.CheckCred(ctx); err != nil {
		if !errors.Is(err, client.ErrUnauthorized) {
			log.Error().Err(err).Msg("cred 校验失败，程序退出")
			exitIfCI()
			return
		}

		log.Warn().Err(err).Msg("cred 已失效，重新登录")

		credResult = a.auth.Reauthenticate(ctx)
		if credResult == nil {
			log.Error().Msg("重新登录失败，程序退出")
			exitIfCI()
			return
		}

		a.setCred(credResult)

		if err := a.auth.CheckCred(ctx); err != nil {
			log.Error().Err(err).Msg("重新登录后 cred 仍然无效，程序退出")
			exitIfCI()
			return
		}
	}

	err := a.checkin.RunAll(ctx, credResult)
	if errors.Is(err, client.ErrUnauthorized) {
		log.Warn().Msg("运行中 cred 失效，重新登录后重试")

		credResult = a.auth.Reauthenticate(ctx)
		if credResult == nil {
			log.Error().Msg("重新登录失败，程序退出")
			exitIfCI()
			return
		}

		a.setCred(credResult)
		err = a.checkin.RunAll(ctx, credResult)
	}

	if err != nil {
		log.Error().Err(err).Msg("签到流程失败，程序退出")
		exitIfCI()
		return
	}

	utils.WaitForExit()
}

func (a *App) setCred(credResult *models.CredResult) {
	a.httpClient.SetCred(credResult.Data.Cred)
	a.httpClient.SetSignToken(credResult.Data.Token)
}

func exitIfCI() {
	if os.Getenv("CI") == "true" {
		os.Exit(1)
	}
}
