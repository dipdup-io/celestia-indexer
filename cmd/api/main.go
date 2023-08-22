package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "DipDup Verticals | Celestia API",
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		return
	}

	if err = initLogger(cfg.LogLevel); err != nil {
		return
	}

	db := initDatabase(cfg.Database)
	e := initEcho(cfg.ApiConfig)
	initHandlers(e, cfg.ApiConfig, db)

	go func() {
		if err := e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
