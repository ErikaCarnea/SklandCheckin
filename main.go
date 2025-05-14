package main

import (
	"Skland/api"
	"Skland/client"
	"Skland/config"
	"Skland/models"
	"Skland/utils"
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CredentialContext struct {
	HttpClient    client.SklandHttpClient
	BindingAPI    *api.BindingAPI
	AttendanceAPI *api.AttendanceAPI
	PlayerAPI     *api.PlayerAPI
	CredResult    *models.CredResult
}

func main() {
	// 初始化日志配置
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// 检测今日是否已经运行过
	//if utils.HasRunToday() {
	//	log.Info().Msg("今日已运行，程序退出")
	//	os.Exit(0)
	//}
	// 创建标记文件
	if err := utils.MarkRun(); err != nil {
		log.Error().Err(err).Msg("无法创建运行标记文件，程序继续执行")
	}

	// 初始化客户端和各API模块
	httpClient := client.NewClient()
	authAPI := api.NewAuthAPI(httpClient)
	bindingAPI := api.NewBindingAPI(httpClient)
	attendanceAPI := api.NewAttendanceAPI(httpClient)
	playerAPI := api.NewPlayerAPI(httpClient)
	checkinAPI := api.NewCheckinAPI(httpClient)

	ctx := &CredentialContext{
		HttpClient:    httpClient,
		BindingAPI:    bindingAPI,
		AttendanceAPI: attendanceAPI,
		PlayerAPI:     playerAPI,
	}
	var credResult *models.CredResult
	token, exists := config.CheckSavedToken()
	if exists {
		credResult = tryAutoLogin(authAPI, token)
		if credResult == nil {
			credResult = loginProcess(authAPI)
		}
	} else {
		credResult = loginProcess(authAPI)
	}
	ctx.CredResult = credResult
	proceedWithCredential(ctx)
	runCheckinTasks(checkinAPI)
	waitForExit()
}

func waitForExit() {
	fmt.Println("\n签到执行完毕，按回车键退出程序...")
	_, _ = fmt.Scanln()
}

func tryAutoLogin(authAPI *api.AuthAPI, token string) *models.CredResult {
	credResult, err := authAPI.GetCredByToken(token)
	if err != nil {
		log.Error().Err(err).Msg("自动登录失败，删除无效token文件")
		if err := os.Remove(config.TokenFileName); err != nil {
			log.Error().Err(err).Msg("删除token文件失败")
		}
		return nil
	}
	log.Info().Msg("检测到有效token，自动登录成功")
	return credResult
}

func loginProcess(authAPI *api.AuthAPI) *models.CredResult {
	// 显示登录选项
	fmt.Println("请选择登录方式:")
	fmt.Println("1. 密码登录 (可能触发人机验证)")
	fmt.Println("2. 手机验证码登录 (可能触发人机验证)")
	fmt.Println("3. 授权码登录")
	fmt.Println("0. 输入\"0\"退出")
	scanner := bufio.NewScanner(os.Stdin)
	var (
		choice     int
		credResult *models.CredResult
		err        error
	)
	for {
		fmt.Print("请输入选项: ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			log.Error().Msg("输入不能为空，请重新输入")
			continue
		}
		if choice, err = strconv.Atoi(input); err != nil {
			log.Error().Msg("请输入有效的数字")
			continue
		}

		switch choice {
		case 1:
			credResult, err = authAPI.LoginByPassword()
		case 2:
			credResult, err = authAPI.LoginByPhoneCode()
		case 3:
			credResult, err = authAPI.LoginByCode()
		case 0:
			os.Exit(1)
		default:
			log.Warn().Msg("无效的登录选项，请重新输入")
			continue
		}
		if err != nil {
			log.Error().Err(err).Msg("登录流程失败,请重新尝试")
		} else {
			break
		}

	}
	return credResult

}

func proceedWithCredential(ctx *CredentialContext) {
	// 设置凭证
	ctx.HttpClient.SetCred(ctx.CredResult.Data.Cred)
	ctx.HttpClient.SetSignToken(ctx.CredResult.Data.Token)

	// 获取绑定列表
	bindings, err := ctx.BindingAPI.GetBindingList()
	if err != nil {
		log.Error().Err(err).Msg("获取绑定列表失败")
		waitForExit()
		os.Exit(1)
	}

	// 打印玩家信息
	//if err := ctx.PlayerAPI.PrintAllPlayersInfo(bindings); err != nil {
	//	log.Error().Err(err).Msg("获取玩家信息失败")
	//	waitForExit()
	//	os.Exit(1)
	//}

	// 执行签到
	signErrors := signAttendance(ctx, bindings)
	if len(signErrors) > 0 {
		for _, err = range signErrors {
			log.Error().Err(err).Msg("签到失败")
		}
	}
}

func signAttendance(ctx *CredentialContext, bindings []models.Binding) []error {
	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, len(bindings))
		errors []error
	)

	for _, binding := range bindings {
		wg.Add(1)
		go func(b models.Binding) {
			defer wg.Done()
			result, err := ctx.AttendanceAPI.SignAttendance(b)
			if err != nil {
				errCh <- err
				return
			}
			log.Info().Msgf("[%s] (%s) %s",
				b.ChannelName,
				b.NickName,
				utils.FormatAttendanceResult(result))
		}(binding)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		errors = append(errors, err)
	}
	return errors
}

func runCheckinTasks(a *api.CheckinAPI) {
	var wg sync.WaitGroup
	for gameID := range api.SklandBoard {
		wg.Add(1)
		time.Sleep(2 * time.Second)
		logger := log.With().Int("gameID", gameID).Str("gameName", api.SklandBoard[gameID]).Logger()
		go func(id int) {
			defer wg.Done()
			_, err := a.Checkin(id)
			if err != nil {
				logger.Error().Err(err).Msg("检票失败")
				return
			}
			logger.Info().Msg("检票成功")
		}(gameID)
	}
	wg.Wait()
}
