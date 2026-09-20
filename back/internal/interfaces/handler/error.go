package handler

import (
	"errors"
	"log/slog"

	"github.com/nishi240masa/meshi-iko/back/internal/domain/apperr"
	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// errorBody はエラーをレスポンスボディに変換する．
// 原因不明のエラーは内部事情が漏れないよう固定の文言にし，詳細はログに残す．
//
// ハンドラがモックの間は呼び出し元が無い．ロジックを実装したら各ハンドラから使う．
//
//nolint:unused // ロジック実装まで未使用．使い始めたらこの行を消す．
func errorBody(err error) openapi.ErrorResponse {
	var appErr *apperr.Error
	if errors.As(err, &appErr) && appErr.Kind != apperr.KindUnknown {
		return openapi.ErrorResponse{Code: appErr.Code, Message: appErr.Message}
	}

	slog.Error("unhandled error", "error", err)
	return openapi.ErrorResponse{Code: "internal_error", Message: "internal server error"}
}
