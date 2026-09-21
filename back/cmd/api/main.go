// Command api はメシイコのHTTP APIサーバを起動する．
// 組み立てと起動だけを行い，具体的な処理は internal 以下にある．
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nishi240masa/meshi-iko/back/internal/config"
	"github.com/nishi240masa/meshi-iko/back/internal/infrastructure/database"
	"github.com/nishi240masa/meshi-iko/back/internal/interfaces/router"
	"github.com/nishi240masa/meshi-iko/back/internal/registry"
)

func main() {
	if err := run(); err != nil {
		slog.Error("サーバの起動に失敗しました", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	gin.SetMode(cfg.GinMode)

	db, err := database.New(cfg.DBPath, cfg.GinMode == gin.DebugMode)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("DBのクローズに失敗しました", "error", err)
		}
	}()

	if err := db.Migrate(); err != nil {
		return err
	}

	handler := registry.New(db).NewHandler()
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: router.New(handler, router.Config{
			CORSAllowedOrigin: cfg.CORSAllowedOrigin,
			AdminToken:        cfg.AdminToken,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 処理中のリクエストを待ってから終了する（docker stop でログが切れないように）．
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("サーバを起動しました", "addr", server.Addr, "db", cfg.DBPath)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("シャットダウンを開始します")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
