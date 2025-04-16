package api

import (
	"Skland/client"
	"Skland/config"
	"Skland/models"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const AppCode = "4ca99fa6b56cc2ba"

type AuthAPI struct {
	client client.HTTPClient
}

func NewAuthAPI(c client.HTTPClient) *AuthAPI {
	return &AuthAPI{client: c}
}

func (a *AuthAPI) LoginByPassword() (*models.CredResult, error) {
	reader := bufio.NewReader(os.Stdin)

	// 记录关键操作
	logrus.WithField("action", "password_login").Info("开始密码登录流程")

	// 读取手机号
	fmt.Print("请输入手机号: ")
	phone, err := reader.ReadString('\n')
	if err != nil {
		logrus.WithError(err).Error("读取手机号失败")
		return nil, fmt.Errorf("读取手机号失败: %w", err)
	}
	phone = strings.TrimSpace(phone)

	var password string
	fmt.Print("请输入密码: ")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			logrus.WithError(err).Error("读取密码失败")
			return nil, fmt.Errorf("读取密码失败: %w", err)
		}
		password = strings.TrimSpace(string(passwordBytes))
		fmt.Println() // 换行
	} else {
		log.Println("警告：当前环境不支持密码隐藏，密码将以明文显示！") // 使用 log 输出
		fmt.Print("请输入密码（明文显示）: ")
		passwordInput, err := reader.ReadString('\n')
		if err != nil {
			logrus.WithError(err).Error("读取密码失败")
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
		"https://as.hypergryph.com/user/auth/v1/token_by_phone_password",
		reqBody,
		&loginResult,
	); err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	if err := config.SaveToken(loginResult.Data.Token); err != nil {
		return nil, fmt.Errorf("保存Token失败: %w", err)
	}
	return a.GetCredByToken(loginResult.Data.Token)
}

func (a *AuthAPI) LoginByPhoneCode() (*models.CredResult, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入手机号: ")
	phone, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取手机号失败: %w", err)
	}
	phone = strings.TrimSpace(phone)

	reqBody := map[string]any{
		"phone": phone,
		"type":  2,
	}

	resp, err := a.client.DoRequest(
		http.MethodPost,
		"https://as.hypergryph.com/general/v1/send_phone_code",
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，HTTP状态码：%d", resp.StatusCode)
	}

	var result struct {
		Status  int    `json:"status"`
		Type    string `json:"type"`
		Message string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Status != 0 {
		return nil, fmt.Errorf(result.Message)
	}

	fmt.Print("请输入验证码: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取验证码失败: %w", err)
	}
	code = strings.TrimSpace(code)

	loginReqBody := map[string]string{
		"phone": phone,
		"code":  code,
	}

	resp, err = a.client.DoRequest(
		http.MethodPost,
		"https://as.hypergryph.com/user/auth/v2/token_by_phone_code",
		loginReqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return nil, fmt.Errorf("登录请求失败: %w", err)
	}

	var loginResult models.LoginResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Status != 0 {
		return nil, fmt.Errorf(result.Message)
	}

	if err = config.SaveToken(loginResult.Data.Token); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	credResult, err := a.GetCredByToken(loginResult.Data.Token)
	if err != nil {
		return nil, fmt.Errorf("获取凭证失败: %v", err)
	}
	return credResult, nil
}

func (a *AuthAPI) LoginByCode() (*models.CredResult, error) {
	logrus.Info("开始授权码登录流程")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("登录森空岛电脑官网后请访问这个网址: https://web-api.skland.com/account/info/hg")
	fmt.Print("请输入获得的内容: ")

	code, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("读取输入失败: %w", err)
	}
	code = strings.TrimSpace(code)

	var response struct {
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(code), &response); err != nil {
		logrus.WithFields(logrus.Fields{
			"raw_input": code,
			"error":     err.Error(),
		}).Error("JSON解析失败")
		return nil, fmt.Errorf("JSON解析失败: %w (原始输入: %s)", err, code)
	}

	if response.Data.Content == "" {
		return nil, fmt.Errorf("无效的数据结构，缺少content字段 (原始输入: %s)", code)
	}

	if err = config.SaveToken(response.Data.Content); err != nil {
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
		"https://zonai.skland.com/api/v1/user/auth/generate_cred_by_code",
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var credResult models.CredResult
	if err := json.NewDecoder(resp.Body).Decode(&credResult); err != nil {
		return nil, err
	}

	if credResult.Code != 0 {
		return nil, fmt.Errorf(credResult.Message)
	}

	return &credResult, nil
}

func (a *AuthAPI) getGrantCode(token string) (string, error) {
	reqBody := map[string]any{
		"appCode": AppCode,
		"token":   token,
		"type":    0,
	}

	resp, err := a.client.DoRequest(
		http.MethodPost,
		"https://as.hypergryph.com/user/oauth2/v2/grant",
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var result models.GrantResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Status != 0 {
		return "", fmt.Errorf(result.Message)
	}
	return result.Data.Code, nil
}
