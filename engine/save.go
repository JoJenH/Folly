package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveGame 原子写入：先写临时文件再 rename
func SaveGame(state GameState, path string) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("重命名失败: %w", err)
	}

	return nil
}

// LoadGame 从文件读取并反序列化游戏状态
func LoadGame(path string) (GameState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GameState{}, fmt.Errorf("读取存档失败: %w", err)
	}

	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return GameState{}, fmt.Errorf("解析存档失败: %w", err)
	}

	return state, nil
}

// SavePath 返回默认存档路径
func SavePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".balatro-cli-save.json"
	}
	return filepath.Join(home, ".config", "balatro-cli", "save.json")
}
