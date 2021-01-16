package wlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = production()

func DevelopMode() {
	logger = develope()
}

func production() *zap.SugaredLogger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	filename := filepath.Base(os.Args[0])
	infoWriter := getWriter("./logs/" + filename)
	warnWriter := getWriter(fmt.Sprintf("./logs/%s_error", filename))

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func develope() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	l, err := config.Build()
	if err != nil {
		log.Fatal("日志初始化失败")
	}

	encoder := zapcore.NewConsoleEncoder(config.EncoderConfig)
	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	filename := filepath.Base(os.Args[0])
	infoWriter := getWriter("./logs/" + filename)
	warnWriter := getWriter(fmt.Sprintf("./logs/%s_error", filename))

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		l.Core(),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	ops := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour * 24 * 7),
	}

	if runtime.GOOS == "linux" {
		ops = append(ops, rotatelogs.WithLinkName(filename))
	}

	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H.log",
		ops...,
	)

	if err != nil {
		panic(err)
	}

	return hook
}

func Info(args ...interface{}) {
	logger.Info(args)
}
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

func Error(args ...interface{}) {
	logger.Error(args)
}
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}
func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}
func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}
