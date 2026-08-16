package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const TokenFileName = "token.txt"

// TokenFilePath 返回 token 文件的完整路径。
func TokenFilePath() (string, error) {
	if os.Getenv("CI") == "true" {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(wd, TokenFileName), nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exePath), TokenFileName), nil
}

// SaveToken 保存 token 到可执行文件同目录下的 token.txt。
// 返回保存的完整路径，由调用方决定是否通知用户。
func SaveToken(token string) (string, error) {
	tokenPath, err := TokenFilePath()
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return "", err
	}

	return tokenPath, nil
}

// DeleteTokenFile 删除 token 文件。用于 token 失效时清理。
func DeleteTokenFile() error {
	tokenPath, err := TokenFilePath()
	if err != nil {
		return err
	}
	return os.Remove(tokenPath)
}

// CheckSavedToken 检查并读取已保存的 token。
// 返回 token 字符串和是否存在。
func CheckSavedToken() (string, bool) {
	tokenPath, err := TokenFilePath()
	if err != nil {
		log.Error().Err(err).Msg("无法获取可执行文件路径")
		return "", false
	}

	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		log.Info().Msg("Token文件不存在")
		return "", false
	}

	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Error().Err(err).Msg("读取token文件失败")
		return "", false
	}
	return strings.TrimSpace(string(tokenBytes)), true
}
