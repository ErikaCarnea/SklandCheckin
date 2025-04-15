package main

import (
	"Skland/api"
	"Skland/client"
	"Skland/config"
	"Skland/models"
	"fmt"
	"log"
	"os"
)

func main() {
	// 初始化客户端
	httpClient := client.NewClient()

	// 初始化各API模块
	authAPI := api.NewAuthAPI(httpClient)
	bindingAPI := api.NewBindingAPI(httpClient)
	attendanceAPI := api.NewAttendanceAPI(httpClient)

	if token, exists := config.CheckSavedToken(); exists {
		if credResult, err := tryAutoLogin(authAPI, token); err != nil {
			proceedWithCredential(httpClient, bindingAPI, attendanceAPI, credResult)
			return
		}
	}

	credResult := loginProcess(authAPI)
	proceedWithCredential(httpClient, bindingAPI, attendanceAPI, credResult)
}

func tryAutoLogin(authAPI *api.AuthAPI, token string) (*models.CredResult, error) {
	credResult, err := authAPI.GetCredByToken(token)
	if err != nil {
		log.Printf("自动登录失败: %v，删除无效token文件", err)
		err := os.Remove(config.TokenFileName)
		if err != nil {
			return nil, err
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
		log.Fatalf("登录失败: %v", err)
	}
	return credResult

}

func proceedWithCredential(
	httpClient *client.HttpClient,
	bindingAPI *api.BindingAPI,
	attendanceAPI *api.AttendanceAPI,
	credResult *models.CredResult) {
	// 设置凭证
	httpClient.SetCred(credResult.Data.Cred)
	httpClient.SetSignToken(credResult.Data.Token)

	// 获取绑定列表
	bindings, err := bindingAPI.GetBindingList()
	if err != nil {
		log.Fatalf("获取绑定列表失败: %v", err)
	}

	// 打印玩家信息
	if err := httpClient.PrintAllPlayersInfo(bindings); err != nil {
		log.Fatalf("获取绑定玩家数据失败: %v", err)
	}

	// 执行签到
	for _, binding := range bindings {
		result, err := attendanceAPI.SignAttendance(binding.Uid, binding.ChannelMasterId)
		if err != nil {
			log.Printf("签到失败：%v", err)
			continue
		}

		fmt.Printf("[%s] (%s) %s\n",
			binding.ChannelName,
			binding.NickName,
			api.FormatAttendanceResult(result),
		)
	}

}
