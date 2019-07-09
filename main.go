package main

import (
	zaplog "log/zapLog"

	"go.uber.org/zap/zapcore"
)

func main() {
	// debug mode
	debug := zaplog.Config{
		Level:              zapcore.DebugLevel, // 日志级别 debug级别下会打印所有类型日志
		EncodeLogsAsJSON:   false,              // 输出日志格式是为json格式
		FileLoggingEnabled: false,              // 输出日志是否保存到文件
		StdLoggingDisabled: false,              // 日志是否标准输入 与上一项保存到文件互斥
		MaxSize:            100000,             // 日志文件最大限制 单位Mb
		MaxBackups:         3,                  // 最大保留备份数
		MaxAge:             7,                  // 保存的天数
		IsAddCaller:        true,               // 是否开启调用追踪
		CallerSkip:         1,                  //
		Directory:          "./out/",
		Filename:           "core.log",
	}

	debug.Init()

	zaplog.Debug("[TestDebugMode] HelloWorld!")
}
