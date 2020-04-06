package main

import (
	"go.uber.org/zap"
)

var SL *zap.SugaredLogger

func init() {
	// 初始化Logger
	logger, _ := zap.NewProduction()
	SL = logger.Sugar()
}
