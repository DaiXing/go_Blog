package blogx

import (
	"log/slog"
	"os"
)

// 日志。
var Logger *slog.Logger

// 配置日志。
func InitLogger() {
	// 配置 JSON 格式，开启 caller 信息便于定位行号
	opts := &slog.HandlerOptions{
		AddSource: true,           // 开启 caller
		Level:     slog.LevelInfo, // 生产环境设置为 Info 级别
	}
	// handler := slog.NewJSONHandler(os.Stdout, opts)// json格式。
	handler := slog.NewTextHandler(os.Stdout, opts) // 文本格式
	logger := slog.New(handler)

	// 替换全局默认 logger
	slog.SetDefault(logger)
	Logger = logger
	Logger.Info("logger 初始化完成")
}
