// Package config は環境変数から設定を読み込む．
//
// Dockerで環境差をなくす前提なので，設定は環境変数のみから受け取り，
// ローカル用の既定値をここに1か所だけ置く．
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config はアプリケーション全体の設定．
type Config struct {
	Port              int           // HTTPサーバの待ち受けポート
	DBPath            string        // SQLiteファイルのパス．":memory:" ならインメモリ
	GinMode           string        // "debug" / "release" / "test"
	ShutdownTimeout   time.Duration // グレースフルシャットダウンの猶予
	CORSAllowedOrigin string        // 許可するOrigin．"*" なら全許可
	// AdminToken は管理用エンドポイントのトークン．空なら管理用は常に401．
	// 設定し忘れで誰でも叩ける状態になるのを防ぐため，既定値は置かない．
	AdminToken string
}

// Load は環境変数を読み，足りないものは既定値で埋める．
func Load() (*Config, error) {
	port, err := envInt("PORT", 8080)
	if err != nil {
		return nil, err
	}
	timeout, err := envInt("SHUTDOWN_TIMEOUT_SECONDS", 10)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:              port,
		DBPath:            envString("DB_PATH", "./data/meshi-iko.db"),
		GinMode:           envString("GIN_MODE", "debug"),
		ShutdownTimeout:   time.Duration(timeout) * time.Second,
		CORSAllowedOrigin: envString("CORS_ALLOWED_ORIGIN", "*"),
		AdminToken:        envString("ADMIN_TOKEN", ""),
	}, nil
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("環境変数 %s は整数である必要があります: %q", key, v)
	}
	return n, nil
}
