// Package persistence は domain/repository のインタフェースをGORMで実装する．
//
// DB用の構造体（*Record）をエンティティとは別に持ち，相互変換をこの層に
// 閉じ込める．カラム追加やGORMタグの都合をドメイン層に漏らさないため．
//
// いまはテーブル定義だけで，リポジトリの実装はまだ無い．
package persistence

import "time"

// UserRecord は users テーブルの1行．
type UserRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:32;not null;uniqueIndex"`
	Token     string    `gorm:"size:64;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName はGORMが使うテーブル名．
func (UserRecord) TableName() string { return "users" }

// AllModels はAutoMigrateの対象一覧．テーブルを増やしたらここに追加する．
func AllModels() []any {
	return []any{
		&UserRecord{},
	}
}
