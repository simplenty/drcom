package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loadConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件不存在")
		}
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件")
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件")
	}

	return config, nil
}

func saveConfig(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return fmt.Errorf("无法编码配置数据")
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("无法写入配置文件")
	}

	return nil
}
