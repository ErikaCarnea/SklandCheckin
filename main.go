package main

import (
	"Skland/api"
	"Skland/client"
	"Skland/config"
	"Skland/models"
	"Skland/utils"
	"fmt"
	"log"
	"os"
	"sync"
)

type CredentialContext struct {
	HttpClient    *client.HttpClient
	BindingAPI    *api.BindingAPI
	AttendanceAPI *api.AttendanceAPI
	PlayerAPI     *api.PlayerApi
	CredResult    *models.CredResult
}

func main() {
	// 初始化客户端
	httpClient := client.NewClient()

	// 初始化各API模块
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
		log.Printf("自动登录失败: %v，删除无效token文件", err)
		if err := os.Remove(config.TokenFileName); err != nil {
			log.Printf("删除token文件失败: %v", err)
		}
		return nil, err
	}
	log.Println("检测到有效token，自动登录成功")
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
		log.Fatalf("无效的选项，请输入 1-3 之间的数字")
	}

	if err != nil {
		log.Printf("%v", err)
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
		log.Printf("获取绑定列表失败: %v", err)
		waitForExit()
		os.Exit(1)
	}

	// 打印玩家信息
	if err := ctx.PlayerAPI.PrintAllPlayersInfo(bindings); err != nil {
		log.Printf("获取绑定玩家数据失败: %v", err)
		waitForExit()
		os.Exit(1)
	}

	// 执行签到
	// 并发执行签到
	var wg sync.WaitGroup
	for _, binding := range bindings {
		wg.Add(1)
		go func(b models.Binding) {
			defer wg.Done()
			result, err := ctx.AttendanceAPI.SignAttendance(b.Uid, b.ChannelMasterId)
			if err != nil {
				log.Printf("签到失败：%v", err)
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
