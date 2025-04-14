package config

import (
	"Skland/models"
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Account models.AccountInfo `yaml:"account"`
}

func LoadConfig() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	//fmt.Printf("exePath: %v\n", exePath)

	configPath := filepath.Join(filepath.Dir(exePath), "config.yaml")
	//fmt.Printf("configPath: %v\n", configPath)

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("未找到配置文件，进入初始化流程...")
		account, err := promptUserForAccount()
		if err != nil {
			return nil, fmt.Errorf("获取账户信息失败: %w", err)
		}

		// 创建并写入新配置文件
		cfg := Config{Account: account}
		data, err := yaml.Marshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("配置文件序列化失败: %w", err)
		}

		if err := os.WriteFile(configPath, data, 0600); err != nil {
			return nil, fmt.Errorf("配置文件创建失败: %w", err)
		}
		fmt.Printf("配置文件已创建在: %s\n", configPath)
		return &cfg, nil
	} else if err != nil {
		return nil, fmt.Errorf("配置文件检查失败: %w", err)
	}

	// 正常读取现有配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Account.Phone == "" || cfg.Account.Password == "" {
		return nil, fmt.Errorf("missing required account configuration")
	}

	return &cfg, nil
}

func promptUserForAccount() (models.AccountInfo, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入手机号: ")
	phone, err := reader.ReadString('\n')
	if err != nil {
		return models.AccountInfo{}, err
	}
	phone = strings.TrimSpace(phone)

	fmt.Print("请输入密码: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return models.AccountInfo{}, err
	}
	password = strings.TrimSpace(password)

	if phone == "" || password == "" {
		return models.AccountInfo{}, fmt.Errorf("手机号和密码不能为空")
	}

	return models.AccountInfo{
		Phone:    phone,
		Password: models.QuotedString(password),
	}, nil
}
