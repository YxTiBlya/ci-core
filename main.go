package main

import "github.com/YxTiBlya/ci-core/logger"

func main() {
	logger.Init(logger.DevelopmentConfig)

	log := logger.New("main")
	log.Info().Msg("hello world")
}
