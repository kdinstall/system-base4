package config

import "os"

// Env はアプリケーション設定を保持する構造体
type Env struct {
	AppName      string
	ServerPort   string
	PlaybooksDir string
}

// GetEnv は設定を返す（環境変数でオーバーライド可能）
func GetEnv() Env {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	playbooksDir := os.Getenv("PLAYBOOKS_DIR")
	if playbooksDir == "" {
		playbooksDir = "../../containers"
	}

	return Env{
		AppName:      "Docker管理",
		ServerPort:   port,
		PlaybooksDir: playbooksDir,
	}
}
