package handler

import (
	"context"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
	"github.com/nishi240masa/meshi-iko/back/internal/usecase"
)

// HealthHandler は死活監視エンドポイントを扱う．
type HealthHandler struct {
	healthUsecase *usecase.HealthUsecase
}

// NewHealthHandler はHealthHandlerを返す．
func NewHealthHandler(healthUsecase *usecase.HealthUsecase) *HealthHandler {
	return &HealthHandler{healthUsecase: healthUsecase}
}

// GetHealth は GET /health を処理する．
//
// DBが落ちていても200を返し，status で "degraded" と伝える．
// ヘルスチェックが叩くため，プロセスの生死とDBの生死を分けたい．
func (h *HealthHandler) GetHealth(ctx context.Context, _ openapi.GetHealthRequestObject) (openapi.GetHealthResponseObject, error) {
	status := openapi.HealthResponseStatus("degraded")
	if h.healthUsecase.Check(ctx) {
		status = openapi.HealthResponseStatus("ok")
	}
	return openapi.GetHealth200JSONResponse{Status: status}, nil
}
