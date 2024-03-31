package logger

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

// go get -u go.uber.org/zap
// go get -u github.com/joho/godotenv

// Logger 全局日志对象
var Logger *zap.Logger

func init() {
	getLoggerEnv()
	initLogger()
	go checkLogFilePathUpdate()
}

// _loggerEnv 日志环境变量
type _loggerEnv struct {
	LogPath  string
	LogLevel string
}

// loggerEnv 日志环境变量
var loggerEnv = &_loggerEnv{}

// 日志配置
var zapLogConfig = zap.NewProductionConfig()

// getLoggerEnv 设置日志环境变量
func getLoggerEnv() {
	err := godotenv.Load()
	if err != nil {
		//	创建.env文件
		_, _ = os.Create(".env")
	}
	// LOG_FILE
	loggerEnv.LogPath = os.Getenv("LOG_PATH")
	if loggerEnv.LogPath == "" {
		loggerEnv.LogPath = "logs"
	}
	// LOG_LEVEL
	loggerEnv.LogLevel = os.Getenv("LOG_LEVEL")
	if loggerEnv.LogLevel == "" {
		loggerEnv.LogLevel = "info"
	}
}

// initLogger 初始化日志对象
func initLogger() {
	zapLogConfig.OutputPaths = []string{getLogFilePath(), "stdout"} // 将日志输出到文件 和 标准输出
	zapLogConfig.Encoding = "console"                               // 设置日志格 json console
	var LevelErr error
	zapLogConfig.Level, LevelErr = zap.ParseAtomicLevel(loggerEnv.LogLevel) // 设置日志级别
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
	//zapLogConfig.Sampling = nil

	// 创建Logger对象
	var buildErr error
	Logger, buildErr = zapLogConfig.Build()
	if buildErr != nil {
		panic("Failed to initialize logger: " + LevelErr.Error())
	}
	// 在应用程序退出时调用以确保所有日志消息都被写入文件
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(Logger)
}

// 函数以当前日期为基础创建日志文件路径
func getLogFilePath() string {
	today := time.Now().Format("2006-01-02")
	filePath := filepath.Join(loggerEnv.LogPath, today)
	// 创建日志文件夹
	_ = os.MkdirAll(loggerEnv.LogPath, os.ModePerm)
	return filePath
}

// 更新日志文件路径
func updateLogFilePath() {
	zapLogConfig.OutputPaths = []string{getLogFilePath(), "stdout"}
	initLogger()
}

// 每秒检查一次日志文件路径是否需要更新
func checkLogFilePathUpdate() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 检查是否需要更新日志文件路径
			if time.Now().Format("2006-01-02") != filepath.Base(zapLogConfig.OutputPaths[0]) {
				updateLogFilePath()
			}
		}
	}
}
