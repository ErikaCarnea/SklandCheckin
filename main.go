package main

import (
	"Skland/api"
	"Skland/client"
	"Skland/config"
	"fmt"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 初始化客户端
	httpClient := client.NewClient()

	// 初始化各API模块
	authAPI := api.NewAuthAPI(httpClient)
	bindingAPI := api.NewBindingAPI(httpClient)
	attendanceAPI := api.NewAttendanceAPI(httpClient)

	// 登录流程
	token, err := authAPI.Login(cfg.Account)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	// 获取凭证
	credResult, err := authAPI.GetCredByToken(token)
	if err != nil {
		log.Fatalf("获取凭证失败: %v", err)
	}

	httpClient.SetCred(credResult.Data.Cred)
	httpClient.SetSignToken(credResult.Data.Token)

	// 获取绑定列表
	bindings, err := bindingAPI.GetBindingList()
	if err != nil {
		log.Fatalf("获取绑定列表失败: %v", err)
	}

	err = httpClient.PrintAllPlayersInfo(bindings)
	if err != nil {
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
