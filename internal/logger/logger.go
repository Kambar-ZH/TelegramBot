package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger, ServiceLogger, HandlerLogger *zap.Logger

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan 01, 2006  15:04:05"))
}

func CustomEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[level])
}

func CustomLevelFileEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + logLevelSeverity[level] + "]")
}

func MyCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.Function)
}

func init() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "temp1.log",
		MaxSize:    1024,
		MaxBackups: 20,
		MaxAge:     28,
		Compress:   true,
	})

	//Define config for the console output
	cfgConsole := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "severity",
		EncodeLevel:  CustomEncodeLevel,
		TimeKey:      "time",
		EncodeTime:   SyslogTimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
	}
	cfgFile := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "severity",
		EncodeLevel:  CustomLevelFileEncoder,
		TimeKey:      "time",
		EncodeTime:   SyslogTimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	consoleDebugging := zapcore.Lock(os.Stdout)
	//consoleError := zapcore.Lock(os.Stderr)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfgFile), w, zap.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(cfgConsole), consoleDebugging, zap.DebugLevel),
		//zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), consoleError, zap.ErrorLevel),
	)
	//core := zapcore.NewCore(zapcore.NewConsoleEncoder(encConsole), w, zap.DebugLevel)
	Logger = zap.New(core, zap.AddCaller())
}
