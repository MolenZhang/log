package log

import (
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Config Configuration for logging
type Config struct {
	Level              zapcore.Level `yaml:"level"`
	EncodeLogsAsJSON   bool          `yaml:"encodeLogsAsJSON"`
	StdLoggingDisabled bool          `yaml:"stdLoggingDisabled"`
	FileLoggingEnabled bool          `yaml:"fileLoggingEnabled"`
	Directory          string        `yaml:"directory"`
	Filename           string        `yaml:"filename"`
	MaxSize            int           `yaml:"maxSize"`    // MB
	MaxBackups         int           `yaml:"maxBackups"` // Days
	MaxAge             int           `yaml:"maxAge"`
	IsAddCaller        bool          `yaml:"isAddCaller "`
	CallerSkip         int           `yaml:"callerSkip"`
	//IsLocalTime        bool
}

var defaultCfg = Config{
	EncodeLogsAsJSON: false,
	Level:            zapcore.InfoLevel,
	IsAddCaller:      true,
	CallerSkip:       1,
}

// DefaultZapLogger 默认logger
var DefaultZapLogger = defaultCfg.newZapLogger(os.Stdout)

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func Sugar() *zap.SugaredLogger {
	return DefaultZapLogger.Sugar()
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Debug(msg, fields...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	DefaultZapLogger.Sugar().Debugf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keyAndValues ...interface{}) {
	DefaultZapLogger.Sugar().Debugw(msg, keyAndValues...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Info(msg, fields...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keyAndValues ...interface{}) {
	DefaultZapLogger.Sugar().Infow(msg, keyAndValues...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	DefaultZapLogger.Sugar().Infof(template, args...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Warn(msg, fields...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	DefaultZapLogger.Sugar().Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keyAndValues ...interface{}) {
	DefaultZapLogger.Sugar().Warnw(msg, keyAndValues...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//func Error(msg string, fields ...zapcore.Field) {
func Error(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Error(msg, fields...)
	//zap.Error()
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	DefaultZapLogger.Sugar().Errorf(template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keyAndValues ...interface{}) {
	DefaultZapLogger.Sugar().Errorw(msg, keyAndValues...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Panic(msg, fields...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	DefaultZapLogger.Sugar().Panicf(template, args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keyAndValues ...interface{}) {
	DefaultZapLogger.Sugar().Panicw(msg, keyAndValues...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Fatal(msg, fields...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	DefaultZapLogger.Sugar().Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keyAndValues ...interface{}) {
	DefaultZapLogger.Sugar().Fatalw(msg, keyAndValues...)
}

// Init 初始化配置log
func (c *Config) Init() {
	var writers []zapcore.WriteSyncer

	if !c.StdLoggingDisabled {
		writers = append(writers, os.Stdout)
	}

	if c.FileLoggingEnabled {
		writers = append(writers, newRollingFile(*c))
	}

	DefaultZapLogger = c.newZapLogger(zapcore.NewMultiWriteSyncer(writers...))

	// zap.RedirectStdLog(DefaultZapLogger)
	Info("logging configured",
		zap.Bool("stdLogging", !c.StdLoggingDisabled),
		zap.Bool("fileLogging", c.FileLoggingEnabled),
		zap.Bool("jsonLogOutput", c.EncodeLogsAsJSON),
		zap.String("logDirectory", c.Directory),
		zap.String("fileName", c.Filename),
		zap.Int("maxSizeMB", c.MaxSize),
		zap.Int("maxBackups", c.MaxBackups),
		zap.Int("maxAgeInDays", c.MaxAge))

	return
}

func newRollingFile(config Config) zapcore.WriteSyncer {
	if err := os.MkdirAll(config.Directory, os.ModePerm); err != nil {
		Error("failed create log directory", zap.Error(err), zap.String("path", config.Directory))
		panic(err)
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,    //megabytes
		MaxAge:     config.MaxAge,     //days
		MaxBackups: config.MaxBackups, //files
	})
}

func (c *Config) newZapLogger(output zapcore.WriteSyncer) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encCfg)
	if c.EncodeLogsAsJSON {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	if c.IsAddCaller {
		opt1 := zap.AddCaller()
		opt2 := zap.AddCallerSkip(c.CallerSkip)
		return zap.New(zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(c.Level)), opt1, opt2)
	}
	return zap.New(zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(c.Level)))
}
