// Package router はGinのエンジンを組み立てる．
//
// ルートの定義は openapi.yaml から生成されたコードが行うので，
// ここでやるのはミドルウェアの適用と，仕様に載らない設定だけ．
package router

import (
	"github.com/gin-gonic/gin"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
	"github.com/nishi240masa/meshi-iko/back/internal/interfaces/middleware"
)

// BasePath は全エンドポイントの接頭辞．openapi.yaml の servers と揃えること．
const BasePath = "/api/v1"

// openapi.yaml の security に対応するオペレーション．仕様を変えたらここも直す．
var (
	sessionTokenOperations = map[string]struct{}{
		"LogoutUser":        {},
		"GetMe":             {},
		"GetUser":           {},
		"GetAnswers":        {},
		"PutMyAnswer":       {},
		"GetPoll":           {},
		"CreateDeviceToken": {},
	}

	adminTokenOperations = map[string]struct{}{
		"UpdatePoll":       {},
		"SendNotification": {},
	}
)

// Config はルータの組み立てに必要な設定．
type Config struct {
	// CORSAllowedOrigin は許可するOrigin．"*" なら全許可．
	CORSAllowedOrigin string
	// AdminToken は管理用エンドポイントのトークン．空なら管理用は常に401になる．
	AdminToken string
}

// New はルーティング済みのGinエンジンを返す．
func New(handler openapi.StrictServerInterface, cfg Config) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.CORS(cfg.CORSAllowedOrigin))

	// 認証をオペレーション単位で切り替えたいので，operationID を受け取れる
	// StrictMiddlewareFunc を使う．
	strict := openapi.NewStrictHandler(handler, []openapi.StrictMiddlewareFunc{
		authMiddleware(cfg.AdminToken),
	})

	openapi.RegisterHandlersWithOptions(engine, strict, openapi.GinServerOptions{
		BaseURL: BasePath,
	})

	return engine
}

func authMiddleware(adminToken string) openapi.StrictMiddlewareFunc {
	session := middleware.SessionAuth()
	admin := middleware.AdminAuth(adminToken)

	return func(next openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		return func(c *gin.Context, request any) (any, error) {
			switch {
			case has(sessionTokenOperations, operationID):
				session(c)
			case has(adminTokenOperations, operationID):
				admin(c)
			default:
				return next(c, request)
			}

			// ミドルウェアが401を書き込んだ場合はハンドラを呼ばない．
			if c.IsAborted() {
				return nil, nil
			}
			return next(c, request)
		}
	}
}

func has(set map[string]struct{}, key string) bool {
	_, ok := set[key]
	return ok
}
