package log

// LoggerConfig 日志配置结构体
type LoggerConfig struct{}

// Logger 日志接口，定义了基本的日志记录方法
type Logger interface {
	// Debug 记录调试级别日志
	Debug(in map[string]any) bool

	// Info 记录信息级别日志
	Info(in map[string]any) bool

	// Warn 记录警告级别日志
	Warn(in map[string]any) bool

	// Error 记录错误级别日志
	Error(in map[string]any) bool

	// Fatal 记录致命错误日志并终止程序
	Fatal(in map[string]any) bool
}
