package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ErikaCarnea/Skland/internal/api"
	"github.com/ErikaCarnea/Skland/internal/config"
	"github.com/ErikaCarnea/Skland/internal/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

// AuthService 认证服务，封装所有登录相关逻辑。
type AuthService struct {
	api *api.AuthAPI
}

func NewAuthService(authAPI *api.AuthAPI) *AuthService {
	return &AuthService{api: authAPI}
}

// Authenticate 执行完整认证流程：先尝试自动登录，失败则进入交互式登录。
// 返回 nil 表示用户取消。
func (a *AuthService) Authenticate(ctx context.Context) *models.CredResult {
	token, exists := config.CheckSavedToken()
	if exists {
		credResult := a.tryAutoLogin(ctx, token)
		if credResult != nil {
			return credResult
		}
	}

	return a.loginProcess(ctx)
}

func (a *AuthService) tryAutoLogin(ctx context.Context, token string) *models.CredResult {
	credResult, err := a.api.GetCredByToken(ctx, token)
	if err != nil {
		log.Error().Err(err).Msg("自动登录失败，删除无效token文件")
		if err := config.DeleteTokenFile(); err != nil {
			log.Error().Err(err).Msg("删除token文件失败")
		}
		return nil
	}
	log.Info().Msg("检测到有效token，自动登录成功")
	return credResult
}

func (a *AuthService) loginProcess(ctx context.Context) *models.CredResult {
	if os.Getenv("CI") == "true" {
		log.Info().Msg("检测到 CI 环境，自动执行密码登录")
		return a.ciPasswordLogin(ctx)
	}

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
			credResult, err = a.passwordLogin(ctx, scanner)
		case 2:
			credResult, err = a.phoneCodeLogin(ctx, scanner)
		case 3:
			credResult, err = a.codeLogin(ctx, scanner)
		case 0:
			return nil
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

func (a *AuthService) passwordLogin(ctx context.Context, scanner *bufio.Scanner) (*models.CredResult, error) {
	fmt.Print("请输入手机号: ")
	scanner.Scan()
	phone := strings.TrimSpace(scanner.Text())

	var password string
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("请输入密码: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, fmt.Errorf("读取密码失败: %w", err)
		}
		password = strings.TrimSpace(string(passwordBytes))
		fmt.Println()
	} else {
		log.Warn().Msg("警告：当前环境不支持密码隐藏，密码将以明文显示！")
		fmt.Print("请输入密码（明文显示）: ")
		scanner.Scan()
		password = strings.TrimSpace(scanner.Text())
	}

	if phone == "" || password == "" {
		return nil, fmt.Errorf("手机号和密码不能为空")
	}

	return a.api.LoginByPassword(ctx, phone, password)
}

func (a *AuthService) phoneCodeLogin(ctx context.Context, scanner *bufio.Scanner) (*models.CredResult, error) {
	fmt.Print("请输入手机号: ")
	scanner.Scan()
	phone := strings.TrimSpace(scanner.Text())
	if phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}

	if err := a.api.SendPhoneCode(ctx, phone); err != nil {
		return nil, fmt.Errorf("发送验证码失败: %w", err)
	}

	fmt.Print("请输入验证码: ")
	scanner.Scan()
	code := strings.TrimSpace(scanner.Text())

	return a.api.LoginByPhoneCode(ctx, phone, code)
}

func (a *AuthService) codeLogin(ctx context.Context, scanner *bufio.Scanner) (*models.CredResult, error) {
	fmt.Println("登录森空岛电脑官网后请访问这个网址: https://web-api.skland.com/account/info/hg")
	fmt.Print("请输入获得的内容: ")

	scanner.Scan()
	code := strings.TrimSpace(scanner.Text())

	return a.api.LoginByCode(ctx, code)
}

func (a *AuthService) ciPasswordLogin(ctx context.Context) *models.CredResult {
	phone := os.Getenv("PHONE")
	password := os.Getenv("PASSWORD")
	if phone == "" {
		log.Error().Msg("CI 环境下未设置 PHONE 环境变量")
		return nil
	}
	if password == "" {
		log.Error().Msg("CI 环境下未设置 PASSWORD 环境变量")
		return nil
	}
	credResult, err := a.api.LoginByPassword(ctx, phone, password)
	if err != nil {
		log.Error().Err(err).Msg("CI 密码登录失败")
		return nil
	}
	return credResult
}
