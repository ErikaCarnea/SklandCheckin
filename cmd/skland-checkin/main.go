package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ErikaCarnea/Skland/internal/app"
	"github.com/ErikaCarnea/Skland/internal/logger"
)

func main() {
	logger.Init()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	a := app.NewApp()
	a.Run(ctx)
}
