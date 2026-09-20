// Package registry は依存関係の組み立て（DI）を1か所に集める．
//
// どの実装を使うかを決めるのはここだけ．各レイヤーは互いを直接newせず，
// コンストラクタ引数で受け取るので，テストでは差し替えが効く．
package registry

import (
	"github.com/nishi240masa/meshi-iko/back/internal/infrastructure/database"
	"github.com/nishi240masa/meshi-iko/back/internal/interfaces/handler"
	"github.com/nishi240masa/meshi-iko/back/internal/usecase"
)

// Registry は組み立て済みの依存を保持する．
type Registry struct {
	db *database.DB
}

// New はRegistryを返す．
func New(db *database.DB) *Registry {
	return &Registry{db: db}
}

// NewHandler は組み立て済みのハンドラを返す．
//
// いまは /health だけが本物のDBを見るので，ほかのハンドラは依存を持たない．
// 実装を入れるときは repository -> usecase -> handler の順に足していく．
func (r *Registry) NewHandler() *handler.Handler {
	healthUsecase := usecase.NewHealthUsecase(r.db)

	return handler.New(
		handler.NewHealthHandler(healthUsecase),
		handler.NewUserHandler(),
		handler.NewAnswerHandler(),
		handler.NewPollHandler(),
		handler.NewDeviceHandler(),
		handler.NewNotificationHandler(),
	)
}
