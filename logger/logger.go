package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

var MLogger *zap.Logger

var infoLevel = zap.LevelEnablerFunc(
	func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	},
)

var errorLevel = zap.LevelEnablerFunc(
	func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel
	},
)

var fatalLevel = zap.LevelEnablerFunc(
	func(level zapcore.Level) bool {
		return level == zapcore.FatalLevel
	},
)

func init() {
	logConfig := zapcore.EncoderConfig{
		MessageKey:    "Message",
		LevelKey:      "Level",
		NameKey:       "LoggerName",
		TimeKey:       "Time",
		CallerKey:     "Caller",
		StacktraceKey: "Stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.TimeEncoderOfLayout(time.RFC850),
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(logConfig)

	logInfoFile := getWriter("./log/log_Info.log")
	logErrorFile := getWriter("./log/log_Error.log")
	logFatalFile := getWriter("./log/log_Fatal.log")

	logInfoFile_Core := zapcore.NewCore(encoder, zapcore.AddSync(logInfoFile), infoLevel)
	logErrorFile_Core := zapcore.NewCore(encoder, zapcore.AddSync(logErrorFile), errorLevel)
	logFatalFile_Core := zapcore.NewCore(encoder, zapcore.AddSync(logFatalFile), fatalLevel)
	logDebugStd_Core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)

	multiIo := zapcore.NewTee(logInfoFile_Core, logErrorFile_Core, logFatalFile_Core, logDebugStd_Core)

	MLogger = zap.New(multiIo)

	MLogger = MLogger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.OnFatal(zapcore.WriteThenPanic), zap.WithCaller(true))
	MLogger = MLogger.Named("MLogger")
	zap.ReplaceGlobals(MLogger)
}

func getWriter(fileName string) io.Writer {
	hook, err := rotatelogs.New(
		//替换文件名
		path.Join("./log",strings.Replace(fileName, ".log", "", -1)+"_%Y-%m-%d.log"),
		//保存多久的日志
		rotatelogs.WithMaxAge(time.Hour*24*30),
		//多久分割一次
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		fmt.Println("getWriter err exit!!")
		os.Exit(1)
	}
	return hook
}
