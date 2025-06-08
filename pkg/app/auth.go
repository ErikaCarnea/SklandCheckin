package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HeathErika/Skland/config"
	"github.com/HeathErika/Skland/models"
	"github.com/rs/zerolog/log"
)

func (ctx *AppContext) Authenticate() bool {
	token, exists := config.CheckSavedToken()
	if exists {
		credResult := ctx.tryAutoLogin(token)
		if credResult != nil {
			ctx.CredResult = credResult
			return true
		}
	}

	ctx.CredResult = ctx.loginProcess()
	return ctx.CredResult != nil
}

func (ctx *AppContext) tryAutoLogin(token string) *models.CredResult {
	credResult, err := ctx.AuthAPI.GetCredByToken(token)
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

func (ctx *AppContext) loginProcess() *models.CredResult {
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
			credResult, err = ctx.AuthAPI.LoginByPassword()
		case 2:
			credResult, err = ctx.AuthAPI.LoginByPhoneCode()
		case 3:
			credResult, err = ctx.AuthAPI.LoginByCode()
		case 0:
			os.Exit(1)
		default:
			log.Warn().Msg("无效的登录选项，请重新输入")
			continue
		}
		if err != nil {
			log.Error().Err(err).Msg("登录流程失败,请重新尝试")
		} else {
			return credResult
		}
	}
}
