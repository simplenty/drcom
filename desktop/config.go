package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func resolveConfigPath(filename string) (string, error) {
	if runtime.GOOS == "windows" {
		return filename, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取配置目录失败")
	}
	return filepath.Join(configDir, "drcom", filename), nil
}

func loadConfig(filename string) (*Config, error) {
	path, err := resolveConfigPath(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件不存在")
		}
		return nil, fmt.Errorf("无法读取配置文件")
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("配置文件格式错误")
	}

	return config, nil
}

func saveConfig(filename string, config Config) error {
	path, err := resolveConfigPath(filename)
	if err != nil {
		return fmt.Errorf("无法解析配置文件")
	}

	if runtime.GOOS != "windows" {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("无法创建配置目录")
		}
	}

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
