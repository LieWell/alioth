package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var Logger *zap.SugaredLogger

func InitZap() {

	cfg := GlobalConfig.Zap

	// 同时写入文件和控制台
	opts := []zapcore.WriteSyncer{
		zapcore.AddSync(&lumberjack.Logger{
			Filename:  cfg.File,
			MaxSize:   cfg.MaxSize,
			MaxAge:    cfg.MaxAge,
			LocalTime: true,
			Compress:  true,
		}),
		zapcore.AddSync(os.Stdout),
	}

	// 打印格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts", // 打印 json 格式时的键名
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,                       // 输出[日志级别]字段
		EncodeTime:     zapcore.TimeEncoderOfLayout(DefaultLogTimeFormat), // 输出[时间]格式
		EncodeDuration: zapcore.SecondsDurationEncoder,                    // 执行消耗的时间转化成浮点型的秒
		EncodeCaller:   zapcore.ShortCallerEncoder,                        // 使用短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(opts...),
			mapperLogLevel(cfg.Level),
		),
		zap.AddCaller(), // 打印文件名与行号
	)
	Logger = logger.Sugar()
}

func mapperLogLevel(l string) zapcore.Level {
	ll := strings.ToLower(l)
	var level zapcore.Level
	switch ll {
	case "debug", "debugger":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	return level
}
