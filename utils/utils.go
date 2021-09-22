package utils

import (
	"fmt"

	naisiov1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetCombinedName(alert *naisiov1.Alert) string {
	return fmt.Sprintf("%s-%s", alert.Namespace, alert.Name)
}

func ZapLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	return cfg.Build()
}
