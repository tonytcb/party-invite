package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/api/http"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/config"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/customerfile"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/logger"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/usecase"
)

func main() {
	var log = logger.NewLogger(os.Stdout)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("error to configuration: %v", err)
	}

	log.Infof("Starting application %s", cfg.AppName)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	var (
		filterCustomers = http.NewFilterCustomersHandler(
			log,
			cfg,
			customerfile.NewCustomersFileParser(),
			usecase.NewFilterCustomers(log),
		)
		httpServer = http.NewServer(log, filterCustomers)
	)

	if err = httpServer.Start(cfg.HTTPPort); err != nil {
		log.Fatalf(err.Error())
	}

	<-done

	if err = httpServer.Stop(context.Background()); err != nil {
		log.Fatalf("error to shutdown http server: %v", err)
	}

	log.Infof("Shutting down application %s", cfg.AppName)
}

func loadConfig() (*config.Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Errorf("error to load current directory: %v", err)
	}

	cfg, err := config.Load(currentDir)
	if err != nil {
		return nil, errors.Errorf("error to load current directory: %v", err)
	}

	if err = cfg.IsValid(); err != nil {
		return nil, errors.Errorf("invalid configuration: %v", err)
	}

	return cfg, nil
}
