package main

import (
	"github.com/HeathErika/Skland/pkg/app"
	"github.com/HeathErika/Skland/pkg/logger"
)

func main() {
	logger.Init()

	ctx := app.NewAppContext()
	ctx.Run()
}
