package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a colorized console zap logger suitable for terminal output.
// If isProduction is true, it sets a higher default log level and fewer stacktraces.
func New(isProduction bool) (*zap.Logger, error) {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.TimeKey = "time"
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "msg"
	encoderCfg.CallerKey = "caller"
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	level := zap.NewAtomicLevel()
	if isProduction {
		level.SetLevel(zapcore.InfoLevel)
	} else {
		level.SetLevel(zapcore.DebugLevel)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)

	stackLevel := zapcore.ErrorLevel
	if !isProduction {
		stackLevel = zapcore.DPanicLevel
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(stackLevel))
	return logger, nil
}
