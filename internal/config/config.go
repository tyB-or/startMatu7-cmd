package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"matu7/pkg/models"
)

// Config 存储所有配置
type Config struct {
	OfflineTools     OfflineToolsConfig
	WebTools         WebToolsConfig
	WebNotes         WebNotesConfig
	ConfigFolderPath string
}

// OfflineToolsConfig 离线工具配置
type OfflineToolsConfig struct {
	ScanPath    string               `json:"scan_path"`
	AutoRefresh bool                 `json:"auto_refresh"`
	Tools       []models.OfflineTool `json:"tools"`
}

// WebToolsConfig 网页工具配置
type WebToolsConfig struct {
	Tools []models.WebTool `json:"tools"`
}

// WebNotesConfig 网页笔记配置
type WebNotesConfig struct {
	Notes []models.Note `json:"notes"`
}

// LoadConfig 加载所有配置文件
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		ConfigFolderPath: configPath,
	}

	// 检查配置路径是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置路径不存在: %s", configPath)
	}

	// 加载离线工具配置
	offlineToolsPath := filepath.Join(configPath, "offline_tools.json")
	if _, err := os.Stat(offlineToolsPath); err == nil {
		offlineData, err := os.ReadFile(offlineToolsPath)
		if err != nil {
			return nil, fmt.Errorf("读取离线工具配置失败: %v", err)
		}

		if err := json.Unmarshal(offlineData, &config.OfflineTools); err != nil {
			return nil, fmt.Errorf("解析离线工具配置失败: %v", err)
		}
	}

	// 加载网页工具配置
	webToolsPath := filepath.Join(configPath, "web_tools.json")
	if _, err := os.Stat(webToolsPath); err == nil {
		webData, err := os.ReadFile(webToolsPath)
		if err != nil {
			return nil, fmt.Errorf("读取网页工具配置失败: %v", err)
		}

		if err := json.Unmarshal(webData, &config.WebTools); err != nil {
			return nil, fmt.Errorf("解析网页工具配置失败: %v", err)
		}
	}

	// 加载笔记配置
	webNotesPath := filepath.Join(configPath, "web_notes.json")
	if _, err := os.Stat(webNotesPath); err == nil {
		notesData, err := os.ReadFile(webNotesPath)
		if err != nil {
			return nil, fmt.Errorf("读取笔记配置失败: %v", err)
		}

		if err := json.Unmarshal(notesData, &config.WebNotes); err != nil {
			return nil, fmt.Errorf("解析笔记配置失败: %v", err)
		}
	}

	return config, nil
}

// SaveConfigPath 保存配置路径到用户主目录
func SaveConfigPath(configPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}

	// 创建.matu7目录
	matu7Dir := filepath.Join(homeDir, ".matu7")
	if err := os.MkdirAll(matu7Dir, 0755); err != nil {
		return fmt.Errorf("创建.matu7目录失败: %v", err)
	}

	// 保存配置路径
	configPathFile := filepath.Join(matu7Dir, "config_path")
	return os.WriteFile(configPathFile, []byte(configPath), 0644)
}

// GetConfigPath 获取保存的配置路径
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户主目录失败: %v", err)
	}

	configPathFile := filepath.Join(homeDir, ".matu7", "config_path")
	if _, err := os.Stat(configPathFile); os.IsNotExist(err) {
		return "", fmt.Errorf("配置路径文件不存在")
	}

	data, err := os.ReadFile(configPathFile)
	if err != nil {
		return "", fmt.Errorf("读取配置路径文件失败: %v", err)
	}

	return string(data), nil
}
