package app

import (
	"github.com/ErikaCarnea/Skland/api"
	"github.com/ErikaCarnea/Skland/client"
	"github.com/ErikaCarnea/Skland/models"
	"github.com/ErikaCarnea/Skland/utils"
	"github.com/rs/zerolog/log"
)

type AppContext struct {
	HttpClient    client.SklandHTTPClient
	AuthAPI       *api.AuthAPI
	BindingAPI    *api.BindingAPI
	AttendanceAPI *api.AttendanceAPI
	PlayerAPI     *api.PlayerAPI
	CheckinAPI    *api.CheckinAPI
	CredResult    *models.CredResult
	PopucomAPI    *api.PopucomAPI
	ExaAPI        *api.ExaAPI
}

func NewAppContext() *AppContext {
	httpClient := client.NewClient()
	return &AppContext{
		HttpClient:    httpClient,
		AuthAPI:       api.NewAuthAPI(httpClient),
		BindingAPI:    api.NewBindingAPI(httpClient),
		AttendanceAPI: api.NewAttendanceAPI(httpClient),
		PlayerAPI:     api.NewPlayerAPI(httpClient),
		CheckinAPI:    api.NewCheckinAPI(httpClient),
		PopucomAPI:    api.NewPopucomAPI(httpClient),
		ExaAPI:        api.NewExaAPI(httpClient),
	}
}

func (ctx *AppContext) Run() {
	if !ctx.Authenticate() {
		log.Error().Msg("登录失败，程序退出")
		return
	}

	ctx.HttpClient.SetCred(ctx.CredResult.Data.Cred)
	ctx.HttpClient.SetSignToken(ctx.CredResult.Data.Token)

	bindings, err := ctx.BindingAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		return
	}

	ctx.GetPopucomAchievement(bindings)
	ctx.GetExaAchievement(bindings)

	hasSigned := ctx.PerformSignAttendance(bindings)

	hasCheckedIn := ctx.RunCheckinTasks()

	if hasSigned && hasCheckedIn {
		log.Info().Msg("今日已签到")
		return
	}
	utils.WaitForExit()
}
