package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init() {
	// 初始化日志配置
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Error().Err(err).Msg("创建日志目录失败")
	}

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "skland.log"),
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile)

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
