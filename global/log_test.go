package global

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	InitLogger()
	Logger.Error("Hi", zap.String("Name", "LiRui"))
}
