// Package repository は永続化の契約だけを定義する．
//
// 実装は infrastructure 層に置き，usecase 層はこのインタフェースにのみ
// 依存する（依存性逆転）．テストではモックを差し込めばDBなしで検証できる．
//
// 実装もこれを使う usecase もまだ無い．ロジックを書くときの設計指針として置いてある．
package repository

import (
	"context"

	"github.com/nishi240masa/meshi-iko/back/internal/domain/model"
)

// UserRepository はユーザーの永続化を担う．
type UserRepository interface {
	// Create はユーザーを保存し，採番されたIDとCreatedAtを埋めて返す．
	// 同名がすでにいる場合は apperr.KindAlreadyExists を返す．
	Create(ctx context.Context, user *model.User) (*model.User, error)

	// FindByID はIDでユーザーを取得する．無ければ apperr.KindNotFound．
	FindByID(ctx context.Context, id int64) (*model.User, error)

	// FindByToken はトークンでユーザーを取得する．無ければ apperr.KindNotFound．
	FindByToken(ctx context.Context, token string) (*model.User, error)
}
