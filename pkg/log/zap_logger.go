package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/space-ark-z/infra-common/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// defaultLogger 是包级别的默认日志记录器
	defaultLogger Logger
	// once 用于确保默认日志记录器只初始化一次
	once sync.Once
	// enableConsole 控制是否启用控制台输出
	enableConsole = getConsoleOutputFromEnv()
)

// getConsoleOutputFromEnv 从环境变量获取控制台输出设置
func getConsoleOutputFromEnv() bool {
	envValue := config.LoadConfig().Get("LOG", "true")
	if envValue == "" {
		return true // 默认启用控制台输出
	}

	// 解析环境变量值
	consoleEnabled, err := strconv.ParseBool(envValue)
	if err != nil {
		// 如果解析失败，默认启用控制台输出
		return true
	}
	return consoleEnabled
}

// GetLogger 获取默认的日志记录器实例
// 外部可以直接调用此函数进行简单日志记录
func GetLogger() Logger {
	once.Do(func() {
		defaultLogger = NewZapLogger()
	})
	return defaultLogger
}

// EnableConsoleOutput 启用或禁用控制台输出
func EnableConsoleOutput(enable bool) {
	enableConsole = enable
}

// Debug 记录调试级别日志
func Debug(in map[string]any) bool {
	return GetLogger().Debug(in)
}

// Info 记录信息级别日志
func Info(in map[string]any) bool {
	return GetLogger().Info(in)
}

// Warn 记录警告级别日志
func Warn(in map[string]any) bool {
	return GetLogger().Warn(in)
}

// Error 记录错误级别日志
func Error(in map[string]any) bool {
	return GetLogger().Error(in)
}

// Fatal 记录致命错误日志并终止程序
func Fatal(in map[string]any) bool {
	return GetLogger().Fatal(in)
}

type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger 创建一个新的ZapLogger实例
// 日志将根据时间戳写入./log/目录下
func NewZapLogger() Logger {
	return NewZapLoggerWithConfig("default", enableConsole)
}

// NewZapLoggerWithModule 创建一个带模块名的ZapLogger实例
// moduleName 用于标识日志来源模块
func NewZapLoggerWithModule(moduleName string) Logger {
	return NewZapLoggerWithConfig(moduleName, enableConsole)
}

// NewZapLoggerWithConfig 创建一个自定义配置的ZapLogger实例
func NewZapLoggerWithConfig(moduleName string, consoleOutput bool) Logger {
	// 确保log目录存在
	logDir := "log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}

	// 生成基于时间戳的文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	var filename string
	if moduleName == "default" {
		filename = filepath.Join(logDir, fmt.Sprintf("log_%s.log", timestamp))
	} else {
		filename = filepath.Join(logDir, fmt.Sprintf("log_%s_%s.log", timestamp, moduleName))
	}

	// 配置zap日志
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建JSON编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建文件写入器
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to create log file: %v", err))
	}

	// 创建写入同步器
	fileWriteSyncer := zapcore.AddSync(file)

	// 创建核心
	var core zapcore.Core
	if consoleOutput {
		// 同时输出到控制台和文件
		consoleWriteSyncer := zapcore.AddSync(os.Stdout)
		// 创建多写入器
		multiWriteSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)
		core = zapcore.NewCore(encoder, multiWriteSyncer, zapcore.DebugLevel)
	} else {
		// 只输出到文件
		core = zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel)
	}

	// 创建zap logger并添加调用者信息
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &ZapLogger{
		logger: logger,
	}
}

// Debug 记录调试级别日志
func (z *ZapLogger) Debug(in map[string]any) bool {
	fields := mapToFields(in)
	z.logger.Debug("", fields...)
	return true
}

// Info 记录信息级别日志
func (z *ZapLogger) Info(in map[string]any) bool {
	fields := mapToFields(in)
	z.logger.Info("", fields...)
	return true
}

// Warn 记录警告级别日志
func (z *ZapLogger) Warn(in map[string]any) bool {
	fields := mapToFields(in)
	z.logger.Warn("", fields...)
	return true
}

// Error 记录错误级别日志
func (z *ZapLogger) Error(in map[string]any) bool {
	fields := mapToFields(in)
	z.logger.Error("", fields...)
	return true
}

// Fatal 记录致命错误日志并终止程序
func (z *ZapLogger) Fatal(in map[string]any) bool {
	fields := mapToFields(in)
	z.logger.Fatal("", fields...)
	panic(fmt.Sprintf("fatal error: %v", in))
	return true
}

// mapToFields 将map转换为zap字段
func mapToFields(data map[string]any) []zap.Field {
	fields := make([]zap.Field, 0, len(data)+1)

	// 添加PID字段
	pid := os.Getpid()
	fields = append(fields, zap.Int("pid", pid))

	// 添加用户提供的字段
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}
