package handler

import (
	"context"
	"time"

	"github.com/oapi-codegen/nullable"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// PollHandler はアンケートの配信・締め切りのエンドポイントを扱う．
type PollHandler struct{}

// NewPollHandler はPollHandlerを返す．
func NewPollHandler() *PollHandler {
	return &PollHandler{}
}

// GetPoll は GET /polls/{date} を処理する．
func (h *PollHandler) GetPoll(_ context.Context, request openapi.GetPollRequestObject) (openapi.GetPollResponseObject, error) {
	return openapi.GetPoll200JSONResponse{
		Date:     request.Date,
		Status:   openapi.Open,
		OpenedAt: nullable.NewNullableWithValue(mockOpenedAt),
		ClosedAt: nullable.NewNullNullable[time.Time](),
	}, nil
}

// UpdatePoll は PUT /polls/{date} を処理する．
func (h *PollHandler) UpdatePoll(_ context.Context, request openapi.UpdatePollRequestObject) (openapi.UpdatePollResponseObject, error) {
	if !validPollStatus(request.Body.Status) {
		return openapi.UpdatePoll400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_status", Message: "status must be one of scheduled, open, closed",
		}}, nil
	}

	openedAt := nullable.NewNullNullable[time.Time]()
	closedAt := nullable.NewNullNullable[time.Time]()
	switch request.Body.Status {
	case openapi.Open:
		openedAt = nullable.NewNullableWithValue(mockOpenedAt)
	case openapi.Closed:
		openedAt = nullable.NewNullableWithValue(mockOpenedAt)
		closedAt = nullable.NewNullableWithValue(mockClosedAt)
	}

	return openapi.UpdatePoll200JSONResponse{
		Date:     request.Date,
		Status:   request.Body.Status,
		OpenedAt: openedAt,
		ClosedAt: closedAt,
	}, nil
}

var (
	mockOpenedAt = time.Date(2026, 9, 20, 16, 0, 0, 0, jst)
	mockClosedAt = time.Date(2026, 9, 20, 18, 0, 0, 0, jst)
)

func validPollStatus(s openapi.PollStatus) bool {
	switch s {
	case openapi.Scheduled, openapi.Open, openapi.Closed:
		return true
	default:
		return false
	}
}
