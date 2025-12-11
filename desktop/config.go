package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	LoginURL       string `json:"login_url"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

func loadConfig(path string) (*Config, error) {
	// 判断文件状态
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件不存在")
		}
		return nil, err
	}
	// 读取文件
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件")
	}
	// 解析文件到 JSON
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件")
	}
	// 返回
	return config, nil
}

func saveConfig(path string, config *Config) error {
	// 编码文件到 JSON
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return fmt.Errorf("无法编码配置数据")
	}
	// 写入文件
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("无法写入配置文件")
	}
	// 返回
	return nil
}
