package logger

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	// 初始化日志配置
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err).Msg("无法获取可执行文件路径")
	}
	exeDir := filepath.Dir(exePath)

	logDir := filepath.Join(exeDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal().Err(err).Msg("创建日志目录失败")
	}

	logFilePath := filepath.Join(logDir, "skland-%Y%m%d.log")

	rotator, err := rotatelogs.New(
		logFilePath,
		rotatelogs.WithMaxAge(3*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("创建 rotatelogs 失败")
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, rotator)

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	if os.Getenv("ENV") == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
