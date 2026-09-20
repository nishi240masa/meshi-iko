// Package middleware はログ・CORS・認証といった横断的関心事を扱う．
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger はリクエストごとに1行のアクセスログを出す．
// 出力形式をアプリ全体で揃えるため，gin標準のロガーの代わりに使う．
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		slog.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}

// CORS はブラウザからの呼び出しを許可する．allowedOrigin が "*" なら全許可．
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := allowedOrigin
		if allowedOrigin != "*" && c.GetHeader("Origin") != allowedOrigin {
			// 許可していないOriginにはCORSヘッダを返さない（ブラウザ側で弾かれる）．
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
