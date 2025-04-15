package config

import (
	"Skland/models"
	"fmt"
	"log"
	"os"
	"strings"
)

const TokenFileName = "token.txt"

type Config struct {
	Account models.AccountInfo `yaml:"account"`
}

//func LoadConfig() (*Config, error) {
//	exePath, err := os.Executable()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get executable path: %w", err)
//	}
//	//fmt.Printf("exePath: %v\n", exePath)
//
//	configPath := filepath.Join(filepath.Dir(exePath), "config.yaml")
//	//fmt.Printf("configPath: %v\n", configPath)
//
//	// 检查配置文件是否存在
//	if _, err := os.Stat(configPath); os.IsNotExist(err) {
//		fmt.Println("未找到配置文件，进入初始化流程...")
//		account, err := promptUserForAccount()
//		if err != nil {
//			return nil, fmt.Errorf("获取账户信息失败: %w", err)
//		}
//
//		// 创建并写入新配置文件
//		cfg := Config{Account: account}
//		data, err := yaml.Marshal(&cfg)
//		if err != nil {
//			return nil, fmt.Errorf("配置文件序列化失败: %w", err)
//		}
//
//		if err := os.WriteFile(configPath, data, 0600); err != nil {
//			return nil, fmt.Errorf("配置文件创建失败: %w", err)
//		}
//		fmt.Printf("配置文件已创建在: %s\n", configPath)
//		return &cfg, nil
//	} else if err != nil {
//		return nil, fmt.Errorf("配置文件检查失败: %w", err)
//	}
//
//	// 正常读取现有配置文件
//	data, err := os.ReadFile(configPath)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read config file: %w", err)
//	}
//
//	var cfg Config
//	if err := yaml.Unmarshal(data, &cfg); err != nil {
//		return nil, fmt.Errorf("failed to parse config: %w", err)
//	}
//
//	if cfg.Account.Phone == "" || cfg.Account.Password == "" {
//		return nil, fmt.Errorf("missing required account configuration")
//	}
//
//	return &cfg, nil
//}

func SaveToken(token string) {
	err := os.WriteFile(TokenFileName, []byte(token), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"您的鹰角网络通行证已经保存在%s, 打开这个可以把它复制到云函数服务器上执行!\n"+
			"如果需要再次运行，删除创建的这个文件即可\n",
		TokenFileName,
	)
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
