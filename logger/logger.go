package logger

import (
	"fmt"
	"github.com/aiyouyo/bluebell/settingts"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/natefinch/lumberjack"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() (err error) {

	writeSyncer := getLogWriter(
		viper.GetString("log.filename"),
		viper.GetInt("log.max_age"),
		viper.GetInt("log.max_backups"),
		viper.GetInt("log.max_size"),
		viper.GetBool("log.compress"),
	)

	encoder := getEncoder()
	l := new(zapcore.Level)
	err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		fmt.Println("日志级别反序列化错误")
		zap.L().Error("l.UnmarshalText failed, err: ", zap.Error(err))
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	lg := zap.New(core, zap.AddCaller()) // 调用信息也会被记录
	// 替换zap 日志库的logger
	zap.ReplaceGlobals(lg)
	return
}

func InitV2(cfg *settingts.LogConfig, mode string) (err error) {

	writeSyncer := getLogWriter(
		cfg.Filename,
		cfg.MaxAge,
		cfg.MaxBackups,
		cfg.MaxSize,
		cfg.Compress,
	)

	encoder := getEncoder()
	l := new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		fmt.Println("日志级别反序列化错误")
		zap.L().Error("l.UnmarshalText failed, err: ", zap.Error(err))
		return
	}
	var core zapcore.Core
	if mode == "dev" {
		// 开发模式将日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), l),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}
	lg := zap.New(core, zap.AddCaller()) // 调用信息也会被记录

	// 替换zap 日志库的logger
	zap.ReplaceGlobals(lg)
	return
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)

		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: err check
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}

}

func getLogWriter(filename string, maxAge, maxSize, maxBackups int, compress bool) zapcore.WriteSyncer {

	file := &lumberjack.Logger{
		Filename:   filename,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		MaxSize:    maxSize,
		Compress:   compress,
	}

	return zapcore.AddSync(file)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}
