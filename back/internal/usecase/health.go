package usecase

import "context"

// Pinger は死活確認できる依存．infrastructure 層の実装に直接依存しないよう，
// 必要な分だけをここで定義する．
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthUsecase はサービスが正常に動いているかを判定する．
type HealthUsecase struct {
	db Pinger
}

// NewHealthUsecase はHealthUsecaseを返す．
func NewHealthUsecase(db Pinger) *HealthUsecase {
	return &HealthUsecase{db: db}
}

// Check はDBに疎通できればtrueを返す．
func (u *HealthUsecase) Check(ctx context.Context) bool {
	return u.db.Ping(ctx) == nil
}
