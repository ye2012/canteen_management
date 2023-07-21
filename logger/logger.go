package logger

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

func init() {
	// 日志级别
	logLevel := "DEBUG"

	atomicLevel := zap.NewAtomicLevel()
	switch logLevel {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "msg",
		FunctionKey:    "",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 日志轮转
	writer := &lumberjack.Logger{
		Filename:   "server.log",
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   true,
	}
	mw := io.MultiWriter(writer, os.Stdout)

	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(mw),
		atomicLevel,
	)

	logger = zap.New(zapCore)
}

func Error(tag, template string, args ...interface{} ) {
	defer logger.Sync()
	message := fmt.Sprintf(template, args...)
	logger.Sugar().Errorw(message, "tag", tag )
}

func Warn(tag, template string, args ...interface{}) {
	defer logger.Sync()
	message := fmt.Sprintf(template, args...)
	logger.Sugar().Warnw(message ,"tag", tag)
}

func Info(tag, template string, args ...interface{}) {
	defer logger.Sync()
	message := fmt.Sprintf(template, args...)
	logger.Sugar().Infow(message ,"tag", tag)
}

func Debug(tag, template string, args ...interface{}) {
	defer logger.Sync()
	message := fmt.Sprintf(template, args...)
	logger.Sugar().Debugw(message, "tag", tag)
}