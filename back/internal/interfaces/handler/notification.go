package handler

import (
	"context"
	"strings"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// NotificationHandler は管理者用のプッシュ通知送信を扱う．
type NotificationHandler struct{}

// NewNotificationHandler はNotificationHandlerを返す．
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// SendNotification は POST /admin/notifications を処理する．
func (h *NotificationHandler) SendNotification(_ context.Context, request openapi.SendNotificationRequestObject) (openapi.SendNotificationResponseObject, error) {
	if strings.TrimSpace(request.Body.Title) == "" || strings.TrimSpace(request.Body.Body) == "" {
		return openapi.SendNotification400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_notification", Message: "title and body must not be empty",
		}}, nil
	}

	// openapi.yaml の default と揃える．
	target := openapi.All
	if request.Body.Target != nil {
		target = *request.Body.Target
	}

	var sentCount int
	switch target {
	case openapi.All:
		sentCount = len(mockUsers)
	case openapi.Unanswered:
		// 3人中1人だけ未回答という想定．
		// unanswered（そもそも回答していない）は undecided（回答したが未定）とは
		// 別物なので，mockAnswers の内訳とは対応しない．
		sentCount = 1
	default:
		return openapi.SendNotification400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_target", Message: "target must be one of all, unanswered",
		}}, nil
	}

	return openapi.SendNotification200JSONResponse{SentCount: sentCount}, nil
}
