package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()

	//logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	//logger.Info("failed to fetch URL")

	sugar := logger.Sugar()
	url := "https://baidu.com"
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
