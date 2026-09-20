package handler

import (
	"context"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// AnswerHandler は回答のエンドポイントを扱う．
type AnswerHandler struct{}

// NewAnswerHandler はAnswerHandlerを返す．
func NewAnswerHandler() *AnswerHandler {
	return &AnswerHandler{}
}

// GetAnswers は GET /answers を処理する．
func (h *AnswerHandler) GetAnswers(_ context.Context, request openapi.GetAnswersRequestObject) (openapi.GetAnswersResponseObject, error) {
	date := today()
	if request.Params.Date != nil {
		date = *request.Params.Date
	}

	return openapi.GetAnswers200JSONResponse{
		Date:    date,
		Answers: mockAnswers(date),
	}, nil
}

// PutMyAnswer は PUT /answers/me を処理する．
func (h *AnswerHandler) PutMyAnswer(_ context.Context, request openapi.PutMyAnswerRequestObject) (openapi.PutMyAnswerResponseObject, error) {
	if !validAnswerStatus(request.Body.Status) {
		return openapi.PutMyAnswer400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_status", Message: "status must be one of undecided, available, unavailable",
		}}, nil
	}

	// 行ける人以外の時間帯は意味を持たないので，仕様どおり空にして返す．
	slots := openapi.TimeSlots{}
	if request.Body.Status == openapi.Available && request.Body.TimeSlots != nil {
		slots = *request.Body.TimeSlots
	}

	return openapi.PutMyAnswer200JSONResponse{
		Date:      today(),
		UserId:    mockUsers[0].Id,
		UserName:  mockUsers[0].Name,
		Status:    request.Body.Status,
		TimeSlots: slots,
	}, nil
}

func validAnswerStatus(s openapi.AnswerStatus) bool {
	switch s {
	case openapi.Undecided, openapi.Available, openapi.Unavailable:
		return true
	default:
		return false
	}
}
