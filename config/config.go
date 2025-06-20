package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const TokenFileName = "token.txt"

func SaveToken(token string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	tokenPath := filepath.Join(exeDir, TokenFileName)

	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return err
	}

	fmt.Printf("您的鹰角网络通行证已经保存在%s!\n", tokenPath)
	return nil
}

func CheckSavedToken() (string, bool) {
	exePath, err := os.Executable()
	if err != nil {
		log.Error().Err(err).Msg("无法获取可执行文件路径")
		return "", false
	}

	exeDir := filepath.Dir(exePath)

	tokenPath := filepath.Join(exeDir, TokenFileName)
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
