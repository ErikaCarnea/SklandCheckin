package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/HeathErika/Skland/client"
	"github.com/HeathErika/Skland/config"
	"github.com/HeathErika/Skland/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

const (
	AppCode             = "4ca99fa6b56cc2ba"
	PhonePasswordURL    = "https://as.hypergryph.com/user/auth/v1/token_by_phone_password"
	SendPhoneCodeURL    = "https://as.hypergryph.com/general/v1/send_phone_code"
	TokenByPhoneCodeUrl = "https://as.hypergryph.com/user/auth/v2/token_by_phone_code"
	GenerateCredURL     = "https://zonai.skland.com/api/v1/user/auth/generate_cred_by_code"
	GrantURL            = "https://as.hypergryph.com/user/oauth2/v2/grant"
)

type AuthAPI struct {
	client client.SklandHttpClient
}

func NewAuthAPI(c client.SklandHttpClient) *AuthAPI {
	return &AuthAPI{client: c}
}

func (a *AuthAPI) LoginByPassword() (*models.CredResult, error) {
	reader := bufio.NewReader(os.Stdin)

	log.Info().Str("action", "password_login").Msg("开始密码登录流程")

	// 读取手机号
	fmt.Print("请输入手机号: ")
	phone, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取手机号失败: %w", err)
	}
	phone = strings.TrimSpace(phone)

	var password string
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("请输入密码: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, fmt.Errorf("读取密码失败: %w", err)
		}
		password = strings.TrimSpace(string(passwordBytes))
		fmt.Println() // 换行
	} else {
		log.Warn().Msg("警告：当前环境不支持密码隐藏，密码将以明文显示！") // 使用 log 输出
		fmt.Print("请输入密码（明文显示）: ")
		passwordInput, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("读取密码失败: %w", err)
		}
		password = strings.TrimSpace(passwordInput)
	}

	if phone == "" || password == "" {
		return nil, fmt.Errorf("手机号和密码不能为空")
	}

	reqBody := map[string]string{
		"phone":    phone,
		"password": password,
	}
	var loginResult models.LoginResult
	if err := a.client.ExecuteRequest(
		http.MethodPost,
		PhonePasswordURL,
		reqBody,
		&loginResult,
	); err != nil {
		return nil, err
	}

	if err := config.SaveToken(loginResult.Data.Token); err != nil {
		log.Error().Err(err).Msg("保存Token失败")
	}
	return a.GetCredByToken(loginResult.Data.Token)
}

func (a *AuthAPI) LoginByPhoneCode() (*models.CredResult, error) {
	reader := bufio.NewReader(os.Stdin)
	log.Info().Str("action", "phone_code_login").Msg("开始手机验证码登录流程")

	fmt.Print("请输入手机号: ")
	phone, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取手机号失败: %w", err)
	}
	phone = strings.TrimSpace(phone)

	var sendCodeResult models.SendCodeResult
	if err := a.client.ExecuteRequest(
		http.MethodPost,
		SendPhoneCodeURL,
		map[string]any{"phone": phone, "type": 2},
		&sendCodeResult,
	); err != nil {
		return nil, fmt.Errorf("发送验证码失败: %w", err)
	}

	if sendCodeResult.GetCode() != 0 {
		return nil, fmt.Errorf("API错误: %s", sendCodeResult.GetMessage())
	}

	fmt.Print("请输入验证码: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取验证码失败: %w", err)
	}
	code = strings.TrimSpace(code)

	var loginResult models.LoginResult
	if err := a.client.ExecuteRequest(
		http.MethodPost,
		TokenByPhoneCodeUrl,
		map[string]string{"phone": phone, "code": code},
		&loginResult,
	); err != nil {
		return nil, err
	}

	if err := config.SaveToken(loginResult.Data.Token); err != nil {
		log.Error().Err(err).Msg("保存Token失败")
	}

	return a.GetCredByToken(loginResult.Data.Token)
}

func (a *AuthAPI) LoginByCode() (*models.CredResult, error) {
	log.Info().Str("action", "token_login").Msg("开始授权码登录流程")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("登录森空岛电脑官网后请访问这个网址: https://web-api.skland.com/account/info/hg")
	fmt.Print("请输入获得的内容: ")

	code, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("token读取失败: %w", err)
	}
	code = strings.TrimSpace(code)

	var response struct {
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(code), &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	if response.Data.Content == "" {
		return nil, fmt.Errorf("无效的数据结构，缺少content字段")
	}

	if err := config.SaveToken(response.Data.Content); err != nil {
		return nil, fmt.Errorf("保存Token失败: %w", err)
	}

	return a.GetCredByToken(response.Data.Content)
}

func (a *AuthAPI) GetCredByToken(token string) (*models.CredResult, error) {
	grantCode, err := a.getGrantCode(token)
	if err != nil {
		return nil, err
	}

	reqBody := map[string]any{
		"code": grantCode,
		"kind": 1,
	}

	resp, err := a.client.DoRequest(
		http.MethodPost,
		GenerateCredURL,
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var credResult models.CredResult
	if err := json.NewDecoder(resp.Body).Decode(&credResult); err != nil {
		return nil, err
	}

	if credResult.Code != 0 {
		return nil, fmt.Errorf("%s", credResult.GetMessage())
	}

	return &credResult, nil
}

func (a *AuthAPI) getGrantCode(token string) (string, error) {
	var result models.GrantResult
	if err := a.client.ExecuteRequest(
		http.MethodPost,
		GrantURL,
		map[string]any{"appCode": AppCode, "token": token, "type": 0},
		&result,
	); err != nil {
		return "", err
	}
	return result.Data.Code, nil
}
