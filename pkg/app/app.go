package app

import (
	"os"

	"github.com/ErikaCarnea/Skland/api"
	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/ErikaCarnea/Skland/utils"
	"github.com/rs/zerolog/log"
)

type AppContext struct {
	HttpClient    client.SklandHttpClient
	AuthAPI       *api.AuthAPI
	BindidngAPI   *api.BindingAPI
	AttendanceAPI *api.AttendanceAPI
	PlayerAPI     *api.PlayerAPI
	CheckinAPI    *api.CheckinAPI
	CredResult    *models.CredResult
}

func NewAppContext() *AppContext {
	httpClient := client.NewClient()
	return &AppContext{
		HttpClient:    httpClient,
		AuthAPI:       api.NewAuthAPI(httpClient),
		BindidngAPI:   api.NewBindingAPI(httpClient),
		AttendanceAPI: api.NewAttendanceAPI(httpClient),
		PlayerAPI:     api.NewPlayerAPI(httpClient),
		CheckinAPI:    api.NewCheckinAPI(httpClient),
	}
}

func (ctx *AppContext) Run() {
	// 检测今日是否已经运行过
	if utils.HasRunToday() {
		log.Info().Msg("今日已运行，程序退出")
		os.Exit(0)
	}
	// 创建标记文件
	if err := utils.MarkRun(); err != nil {
		log.Error().Err(err).Msg("无法创建运行标记文件，程序继续执行")
	}

	if !ctx.Authenticate() {
		log.Error().Msg("登录失败，程序退出")
		utils.WaitForExit()
		return
	}

	ctx.HttpClient.SetCred(ctx.CredResult.Data.Cred)
	ctx.HttpClient.SetSignToken(ctx.CredResult.Data.Token)

	ctx.PerformSignAttendance()

	ctx.RunCheckinTasks()

	utils.WaitForExit()
}
