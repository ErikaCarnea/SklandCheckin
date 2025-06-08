package main

import (
	"github.com/ErikaCarnea/Skland/pkg/app"
	"github.com/ErikaCarnea/Skland/pkg/logger"
)

func main() {
	logger.Init()

	ctx := app.NewAppContext()
	ctx.Run()
}
