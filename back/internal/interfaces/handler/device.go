package handler

import (
	"context"
	"strings"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// DeviceHandler はプッシュ通知の宛先登録を扱う．
type DeviceHandler struct{}

// NewDeviceHandler はDeviceHandlerを返す．
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

// CreateDeviceToken は POST /devices を処理する．
func (h *DeviceHandler) CreateDeviceToken(_ context.Context, request openapi.CreateDeviceTokenRequestObject) (openapi.CreateDeviceTokenResponseObject, error) {
	if strings.TrimSpace(request.Body.DeviceToken) == "" {
		return openapi.CreateDeviceToken400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_device_token", Message: "deviceToken must not be empty",
		}}, nil
	}
	if request.Body.OsType != openapi.IOS && request.Body.OsType != openapi.Android {
		return openapi.CreateDeviceToken400JSONResponse{BadRequestJSONResponse: openapi.BadRequestJSONResponse{
			Code: "invalid_os_type", Message: "osType must be one of iOS, Android",
		}}, nil
	}

	return openapi.CreateDeviceToken201Response{}, nil
}
