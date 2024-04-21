package logger

import "testing"

func Test(t *testing.T) {
	Logger.Info("Hello World")
	Logger.Debug("Hello World")
	Logger.Error("Hello World")
}
