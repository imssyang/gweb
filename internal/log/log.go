package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(&appWriter{
		logger:     log.New(os.Stdout, "", 0),
		prefix:     sysFlag,
		hasTime:    true,
		resetLevel: false,
	})

	Zap = zapLogger(&appWriter{
		logger:     log.New(os.Stdout, "", 0),
		prefix:     zapFlag,
		hasTime:    true,
		resetLevel: true,
	})
}

const (
	sysFlag = "SYS"
	ginFlag = "GIN"
	zapFlag = "ZAP"
)

type appWriter struct {
	logger     *log.Logger
	prefix     string
	hasTime    bool
	resetLevel bool
}

func (w *appWriter) splitLevel(p []byte) (zapcore.Level, []byte) {
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}
	for _, level := range levels {
		prefix := []byte(level.String())
		trimmed := bytes.TrimPrefix(p, prefix)
		if len(trimmed) < len(p) {
			return level, trimmed
		}
	}

	return zapcore.InvalidLevel, p
}

func (w *appWriter) Write(p []byte) (n int, err error) {
	var timestamp string
	if w.hasTime {
		timestamp = time.Now().Format("2006/1/2 15:04:05.000 ")
	}

	var level string
	if w.resetLevel {
		l, trimmed := w.splitLevel(p)
		if l != zapcore.InfoLevel {
			level = "-" + l.String()
		}
		p = trimmed
	}

	s := fmt.Sprintf("[%s%s] %s%s", w.prefix, level, timestamp, string(p))
	return len(p), w.logger.Output(2, s)
}

func zapLogger(w io.Writer) *zap.SugaredLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncodeTime = nil
	cfg.ConsoleSeparator = "\t"

	level := zapcore.InfoLevel
	if conf.App.Debug {
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(w),
		level,
	)

	return zap.New(core).Sugar()
}

var Zap *zap.SugaredLogger

func logFormatter(param gin.LogFormatterParams) string {
	timestamp := param.TimeStamp.Format("2006/1/2 15:04:05.000")
	return fmt.Sprintf("[%s] %s [%s %s] %d %s %s %s - %s\n",
		ginFlag,
		timestamp,
		param.ClientIP,
		param.Request.Proto,
		param.StatusCode,
		param.Method,
		param.Path,
		param.Latency,
		param.ErrorMessage,
	)
}

func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(logFormatter)
}
