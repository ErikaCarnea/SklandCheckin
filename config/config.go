package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const TokenFileName = "token.txt"

func SaveToken(token string) error {
	if err := os.WriteFile(TokenFileName, []byte(token), 0600); err != nil {
		return err
	}

	fmt.Printf("您的鹰角网络通行证已经保存在%s!\n", TokenFileName)
	return nil
}

func CheckSavedToken() (string, bool) {
	if _, err := os.Stat(TokenFileName); os.IsNotExist(err) {
		return "", false
	}

	tokenBytes, err := os.ReadFile(TokenFileName)
	if err != nil {
		log.Error().Err(err).Msg("读取token文件失败")
		return "", false
	}
	return strings.TrimSpace(string(tokenBytes)), true
}
