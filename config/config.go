package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const TokenFileName = "token.txt"

func SaveToken(token string) error {
	if err := os.WriteFile(TokenFileName, []byte(token), 0644); err != nil {
		return fmt.Errorf("保存Token到文件失败: %w", err)
	}

	fmt.Printf(
		"您的鹰角网络通行证已经保存在%s, 打开这个可以把它复制到云函数服务器上执行!\n"+
			"如果需要再次运行，删除创建的这个文件即可\n",
		TokenFileName,
	)
	return nil
}

func CheckSavedToken() (string, bool) {
	if _, err := os.Stat(TokenFileName); os.IsNotExist(err) {
		return "", false
	}

	tokenBytes, err := os.ReadFile(TokenFileName)
	if err != nil {
		log.Printf("读取token文件失败: %v", err)
		return "", false
	}
	return strings.TrimSpace(string(tokenBytes)), true
}
