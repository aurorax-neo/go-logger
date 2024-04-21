package logger

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

// go get -u go.uber.org/zap
// go get -u github.com/joho/godotenv

var (
	Logger       *zap.Logger
	zapLogConfig = zap.NewProductionConfig()
)

func init() {
	newLogger()
}

// getLoggerEnv 设置日志环境变量
func getLoggerEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		//	创建.env文件
		_, _ = os.Create(".env")
	}
	// LOG_FILE
	logPath := os.Getenv("LOG_PATH")
	// LOG_LEVEL
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	return logPath, logLevel
}

// newLogger 初始化日志对象
func newLogger() *zap.Logger {
	// 获取日志环境变量
	logPath, logLevel := getLoggerEnv()
	// 设置日志配置
	if logPath == "" {
		zapLogConfig.OutputPaths = []string{"stdout"} // 标准输出
	} else {
		zapLogConfig.OutputPaths = []string{getLogFilePath(logPath), "stdout"} // 将日志输出到文件 和 标准输出
	}
	zapLogConfig.Encoding = "console" // 设置日志格 json console
	var LevelErr error
	zapLogConfig.Level, LevelErr = zap.ParseAtomicLevel(logLevel) // 设置日志级别
	if LevelErr != nil {
		zapLogConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zapLogConfig.EncoderConfig = zapcore.EncoderConfig{ // 创建Encoder配置
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	// 创建Logger对象
	var buildErr error
	Logger, buildErr = zapLogConfig.Build()
	if buildErr != nil {
		panic(fmt.Sprint("Failed to initialize logger: ", LevelErr))
	}
	// 在应用程序退出时调用以确保所有日志消息都被写入文件
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(Logger)
	// 检查日志文件路径是否需要更新
	if logPath != "" {
		go checkLogFilePathUpdate(logPath)
	}
	return Logger
}

// 函数以当前日期为基础创建日志文件路径
func getLogFilePath(path string) string {
	today := time.Now().Format("2006-01-02")
	filePath := filepath.Join(path, today)
	// 创建日志文件夹
	_ = os.MkdirAll(path, os.ModePerm)
	return filePath
}

// 更新日志文件路径
func updateLogFilePath(path string) {
	zapLogConfig.OutputPaths = []string{getLogFilePath(path), "stdout"}
	Logger = newLogger()
}

// 每秒检查一次日志文件路径是否需要更新
func checkLogFilePathUpdate(path string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 检查是否需要更新日志文件路径
			if time.Now().Format("2006-01-02") != filepath.Base(zapLogConfig.OutputPaths[0]) {
				updateLogFilePath(path)
			}
		}
	}
}
