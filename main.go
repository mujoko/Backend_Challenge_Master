package main

import (
	"fmt"
	"go-inventory/config"
	"go-inventory/router"
	lr "go-inventory/util/logger"
	"net/http"
)

func main() {

	appConf := config.AppConfig()
	address := fmt.Sprintf(":%d", appConf.Server.Port)

	logger := lr.New(appConf.Debug)

	r := router.New()
	// log.Printf("Starting server %s\n", address)
	logger.Info().Msgf("Starting server %v", address)

	s := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failed")
	}

}
