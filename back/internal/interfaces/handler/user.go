package handler

import (
	"context"
	"strings"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// UserHandler はユーザー関連のエンドポイントを扱う．
type UserHandler struct{}

// NewUserHandler はUserHandlerを返す．
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// CreateUser は POST /users を処理する．
func (h *UserHandler) CreateUser(_ context.Context, request openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	name := strings.TrimSpace(request.Body.Name)
	if name == "" {
		return openapi.CreateUser400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_name", Message: "name must not be empty",
		}}, nil
	}

	return openapi.CreateUser201JSONResponse{
		User:  openapi.User{Id: 1, Name: name, CreatedAt: mockCreatedAt(0)},
		Token: mockToken,
	}, nil
}

// LoginUser は POST /users/login を処理する．
func (h *UserHandler) LoginUser(_ context.Context, request openapi.LoginUserRequestObject) (openapi.LoginUserResponseObject, error) {
	name := strings.TrimSpace(request.Body.Name)
	if name == "" {
		return openapi.LoginUser400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_name", Message: "name must not be empty",
		}}, nil
	}

	return openapi.LoginUser200JSONResponse{
		User:  openapi.User{Id: 1, Name: name, CreatedAt: mockCreatedAt(0)},
		Token: mockToken,
	}, nil
}

// LogoutUser は POST /users/logout を処理する．
func (h *UserHandler) LogoutUser(_ context.Context, _ openapi.LogoutUserRequestObject) (openapi.LogoutUserResponseObject, error) {
	return openapi.LogoutUser204Response{}, nil
}

// GetMe は GET /users/me を処理する．
func (h *UserHandler) GetMe(_ context.Context, _ openapi.GetMeRequestObject) (openapi.GetMeResponseObject, error) {
	return openapi.GetMe200JSONResponse(mockUsers[0]), nil
}

// GetUser は GET /users/{userId} を処理する．
func (h *UserHandler) GetUser(_ context.Context, request openapi.GetUserRequestObject) (openapi.GetUserResponseObject, error) {
	user, ok := mockUserByID(request.UserId)
	if !ok {
		return openapi.GetUser404JSONResponse{NotFoundJSONResponse: openapi.NotFoundJSONResponse{
			Code: "not_found", Message: "user not found",
		}}, nil
	}
	return openapi.GetUser200JSONResponse(user), nil
}
