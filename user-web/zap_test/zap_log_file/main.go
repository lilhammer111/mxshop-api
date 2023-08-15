package main

import (
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"./log/user_log/user.log", "stderr", "stdout"}
	return cfg.Build()
}

func main() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	url := "https://imooc.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
