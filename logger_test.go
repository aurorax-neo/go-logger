package logger

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	Logger.Info("Hello World")
	Logger.Debug("Hello World")
	Logger.Error("Hello World")

	for i := 0; i < 1000; i++ {
		Logger.Info("Hello World")
		Logger.Debug("Hello World")
		Logger.Error("Hello World")
		time.Sleep(1 * time.Second)
	}
}
