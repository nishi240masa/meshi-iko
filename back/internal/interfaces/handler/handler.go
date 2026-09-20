// Package handler はHTTPの入出力を担当する．
// 型の変換とステータスの決定だけを行い，ビジネスロジックは usecase に置く．
//
// いまは /health 以外すべてモックで，固定データを返す．
// 各メソッドのコメントに書いたパスは router.BasePath（/api/v1）からの相対．
package handler

import (
	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// Handler は各リソースのハンドラをまとめ，生成されたインタフェースを満たす．
// 埋め込みでメソッドが引き上げられるので，リソースを増やすときはフィールドを足す．
type Handler struct {
	*HealthHandler
	*UserHandler
	*AnswerHandler
	*PollHandler
	*DeviceHandler
	*NotificationHandler
}

var _ openapi.StrictServerInterface = (*Handler)(nil)

// New はHandlerを組み立てる．
func New(
	health *HealthHandler,
	user *UserHandler,
	answer *AnswerHandler,
	poll *PollHandler,
	device *DeviceHandler,
	notification *NotificationHandler,
) *Handler {
	return &Handler{
		HealthHandler:       health,
		UserHandler:         user,
		AnswerHandler:       answer,
		PollHandler:         poll,
		DeviceHandler:       device,
		NotificationHandler: notification,
	}
}
