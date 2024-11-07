package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultPrefix = "v"
	defaultSuffix = ""
)

type Config struct {
	DefaultPrefix string `json:"default_prefix"`
	DefaultSuffix string `json:"default_suffix"`
}

// getConfigDir 获取配置文件目录
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".tagger"), nil
}

// getConfigPath 获取配置文件路径
func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// ensureConfigDir 确保配置目录存在
func ensureConfigDir() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	
	// 检查目录是否存在，不存在则创建
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return os.MkdirAll(configDir, 0755)
	}
	return nil
}

// loadConfig 加载配置文件
func loadConfig() (*Config, error) {
	// 确保配置目录存在
	if err := ensureConfigDir(); err != nil {
		return nil, err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// 默认配置
	config := &Config{
		DefaultPrefix: defaultPrefix,
		DefaultSuffix: defaultSuffix,
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，保存默认配置
			if err := saveConfig(config); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}

	// 解析配置文件
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

// saveConfig 保存配置到文件
func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 确保配置目录存在
	if err := ensureConfigDir(); err != nil {
		return err
	}

	// 格式化 JSON 数据
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(configPath, data, 0644)
}

// SetDefaultPrefix 设置默认前缀
func SetDefaultPrefix(prefix string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}
	config.DefaultPrefix = prefix
	return saveConfig(config)
}

// SetDefaultSuffix 设置默认后缀
func SetDefaultSuffix(suffix string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}
	config.DefaultSuffix = suffix
	return saveConfig(config)
}

// GetDefaultPrefix 获取默认前缀
func GetDefaultPrefix() (string, error) {
	config, err := loadConfig()
	if err != nil {
		return defaultPrefix, err
	}
	return config.DefaultPrefix, nil
}

// GetDefaultSuffix 获取默认后缀
func GetDefaultSuffix() (string, error) {
	config, err := loadConfig()
	if err != nil {
		return defaultSuffix, err
	}
	return config.DefaultSuffix, nil
}
