package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	openapi "github.com/nishi240masa/meshi-iko/back/internal/generated/openapi"
)

// ContextKeyUser は認証情報をコンテキストに入れるときのキー．
// モックのいまは検証していないトークン文字列が入る．本実装では *model.User に変える．
//
// gin.Context がそのまま context.Context として渡るので，ハンドラ側は
// ctx.Value(ContextKeyUser) で取り出せる．
const ContextKeyUser = "authenticated_user"

// SessionAuth は Authorization: Bearer <トークン> を検証する．
// openapi.yaml の SessionToken に対応する．
//
// いまはモックなのでトークンの中身は見ず，空かどうかだけを判定する．
// 本実装ではトークンからユーザーを引く処理に差し替える．
func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
		if token == "" {
			abortUnauthorized(c, "invalid or missing token")
			return
		}

		c.Set(ContextKeyUser, token)
		c.Next()
	}
}

// AdminAuth は X-Admin-Token が expected と一致することを確認する．
// expected が空のときは，設定し忘れで全開放になるのを防ぐため常に拒否する．
func AdminAuth(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got := c.GetHeader("X-Admin-Token")
		if expected == "" || subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			abortUnauthorized(c, "invalid or missing admin token")
			return
		}
		c.Next()
	}
}

// CurrentToken はSessionAuthが通したトークンを返す．
// ハンドラがモックの間は呼び出し元が無い．
func CurrentToken(c *gin.Context) (string, bool) {
	v, ok := c.Get(ContextKeyUser)
	if !ok {
		return "", false
	}
	token, ok := v.(string)
	return token, ok
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.ErrorResponse{
		Code:    "unauthorized",
		Message: message,
	})
}
