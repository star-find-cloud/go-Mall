package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/star-find-cloud/star-mall/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	AppLogger   *zap.SugaredLogger // 主日志
	MySQLLogger *zap.SugaredLogger // 数据库专用日志
	HttpLogger  *zap.SugaredLogger
	RedisLogger *zap.SugaredLogger
	EtcdLogger  *zap.SugaredLogger
)

func init() {
	encoderConig := zap.NewProductionEncoderConfig()
	encoderConig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConig.EncodeLevel = zapcore.CapitalLevelEncoder

	appCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConig), getLogWriter("app"), zap.DebugLevel)
	dbCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConig), getLogWriter("mysql"), zap.DebugLevel)
	httpCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConig), getLogWriter("http"), zap.DebugLevel)
	redisCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConig), getLogWriter("redis"), zap.DebugLevel)
	etcdCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConig), getLogWriter("etcd"), zap.DebugLevel)

	AppLogger = zap.New(appCore, zap.AddCaller()).Sugar()
	MySQLLogger = zap.New(dbCore, zap.AddCaller()).Sugar()
	HttpLogger = zap.New(httpCore, zap.AddCaller()).Sugar()
	RedisLogger = zap.New(redisCore, zap.AddCaller()).Sugar()
	EtcdLogger = zap.New(etcdCore, zap.AddCaller()).Sugar()
}

func getLogWriter(model string) zapcore.WriteSyncer {
	c := conf.GetConfig()
	fileName := fmt.Sprintf("%s/mall-%s-%s-%s.log", c.Log.Path, c.Log.Version, c.Log.Level, model)
	lumberJack := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    c.Log.MaxSize,
		MaxBackups: c.Log.MaxBackup,
		MaxAge:     c.Log.MaxAge,
		Compress:   c.Log.Compress,
	}
	return zapcore.AddSync(lumberJack)
}
