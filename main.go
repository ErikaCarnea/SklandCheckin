package main

import (
	"Skland/api"
	"Skland/client"
	"Skland/config"
	"Skland/models"
	"Skland/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
	"time"
)

type CredentialContext struct {
	HttpClient    client.HTTPClient
	BindingAPI    *api.BindingAPI
	AttendanceAPI *api.AttendanceAPI
	PlayerAPI     *api.PlayerApi
	CredResult    *models.CredResult
}

func main() {
	//初始化日志配置
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// 初始化客户端和各API模块
	httpClient := client.NewClient()
	authAPI := api.NewAuthAPI(httpClient)
	bindingAPI := api.NewBindingAPI(httpClient)
	attendanceAPI := api.NewAttendanceAPI(httpClient)
	playerAPI := api.NewPlayerAPI(httpClient)

	ctx := &CredentialContext{
		HttpClient:    httpClient,
		BindingAPI:    bindingAPI,
		AttendanceAPI: attendanceAPI,
		PlayerAPI:     playerAPI,
	}

	if token, exists := config.CheckSavedToken(); exists {
		if credResult, err := tryAutoLogin(authAPI, token); err == nil {
			ctx.CredResult = credResult
			proceedWithCredential(ctx)
			waitForExit()
			return
		}
	}

	credResult := loginProcess(authAPI)
	ctx.CredResult = credResult
	proceedWithCredential(ctx)
	time.Sleep(3000 * time.Millisecond)
	waitForExit()
}

func waitForExit() {
	fmt.Println("\n签到执行完毕，按回车键退出程序...")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
}

func tryAutoLogin(authAPI *api.AuthAPI, token string) (*models.CredResult, error) {
	credResult, err := authAPI.GetCredByToken(token)
	if err != nil {
		logrus.WithError(err).Error("自动登录失败，删除无效token文件")
		if err := os.Remove(config.TokenFileName); err != nil {
			logrus.WithError(err).Error("删除token文件失败")
		}
		return nil, err
	}
	logrus.Info("检测到有效token，自动登录成功")
	return credResult, nil
}

func loginProcess(authAPI *api.AuthAPI) *models.CredResult {
	// 显示登录选项
	fmt.Println("请选择登录方式:")
	fmt.Println("1. 密码登录")
	fmt.Println("2. 手机验证码登录")
	fmt.Println("3. 授权码登录")
	fmt.Print("请输入选项数字 (1-3): ")

	var choice int
	if _, err := fmt.Scanln(&choice); err != nil {
		log.Fatalf("输入读取失败: %v", err)
	}
	var (
		credResult *models.CredResult
		err        error
	)

	// 登录流程
	switch choice {
	case 1:
		credResult, err = authAPI.LoginByPassword()
	case 2:
		credResult, err = authAPI.LoginByPhoneCode()
	case 3:
		credResult, err = authAPI.LoginByCode()
	default:
		logrus.Error("无效的登录选项")
		os.Exit(1)
	}

	if err != nil {
		logrus.WithError(err).Error("登录流程失败")
		waitForExit()
		os.Exit(1)
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
		logrus.WithError(err).Error("获取绑定列表失败")
		waitForExit()
		os.Exit(1)
	}

	// 打印玩家信息
	//if err := ctx.PlayerAPI.PrintAllPlayersInfo(bindings); err != nil {
	//	logrus.WithError(err).Error("获取玩家信息失败")
	//	waitForExit()
	//	os.Exit(1)
	//}

	// 执行签到
	// 并发执行签到
	var wg sync.WaitGroup
	for _, binding := range bindings {
		wg.Add(1)
		go func(b models.Binding) {
			defer wg.Done()
			result, err := ctx.AttendanceAPI.SignAttendance(b)
			if err != nil {
				return
			}

			fmt.Printf("[%s] (%s) %s\n",
				b.ChannelName,
				b.NickName,
				utils.FormatAttendanceResult(result),
			)
		}(binding)
	}
	wg.Wait()

}
