// Package database はDB接続を扱う．GORMに依存してよい唯一の層．
package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nishi240masa/meshi-iko/back/internal/infrastructure/persistence"
)

// DB はGORMの接続を包む．
type DB struct {
	*gorm.DB
}

// New はSQLiteに接続する．dbPathの親ディレクトリが無ければ作る．
// dbPathに ":memory:" を渡すとインメモリDBになる（テスト用）．
func New(dbPath string, debug bool) (*DB, error) {
	var dsn string
	if dbPath == ":memory:" {
		// インメモリでも複数コネクションから同じDBを見えるようにする．
		dsn = "file::memory:?cache=shared"
	} else {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, fmt.Errorf("DBディレクトリの作成に失敗しました: %w", err)
		}
		// _pragma=... はglebarez/sqliteのDSNオプション．
		// 外部キー制約を有効にし，書き込みと読み込みが競合しにくいWALモードにする．
		dsn = dbPath + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	}

	logLevel := logger.Warn
	if debug {
		logLevel = logger.Info
	}

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		// UNIQUE制約違反などを gorm.ErrDuplicatedKey に変換させる．
		// これがないと，リポジトリ実装がドライバ固有のエラー文字列を見る羽目になる．
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("SQLiteへの接続に失敗しました: %w", err)
	}

	return &DB{gormDB}, nil
}

// Migrate はテーブルを作成・更新する．
//
// カラム削除やリネームはAutoMigrateでは反映されない．必要になったら
// マイグレーションツールの導入を検討すること．
func (db *DB) Migrate() error {
	if err := db.AutoMigrate(persistence.AllModels()...); err != nil {
		return fmt.Errorf("マイグレーションに失敗しました: %w", err)
	}
	return nil
}

// Ping はDBに疎通できるかを確認する．usecase.Pinger を満たす．
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Close は接続を閉じる．
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
