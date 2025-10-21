package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kardianos/service"
	"github.com/rs/zerolog"
	"github.com/ttrnecka/wwn_identity/webapi/db"
	"github.com/ttrnecka/wwn_identity/webapi/server"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"

	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger

func init() {
	logger = logging.SetupLogger("http")
}

type program struct {
	exit   chan struct{}
	server *http.Server
	wg     sync.WaitGroup
}

func (p *program) runServer() {
	gob.Register(dto.UserDTO{})

	ex, err := os.Executable()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	exPath := filepath.Dir(ex)
	envFile := filepath.Join(exPath, ".env")

	// if err := godotenv.Load(envFile); err != nil {
	// 	logger.Info().Msg("No .env file found, using default config")
	// }

	envMap, err := godotenv.Read(envFile)
	if err != nil {
		log.Println("No .env file found")
	}

	// Only set env vars that aren't already set
	for k, v := range envMap {
		if _, exists := os.LookupEnv(k); !exists {
			err := os.Setenv(k, v)
			if err != nil {
				logger.Fatal().Err(err).Msg("")
			}
		}
	}

	// db
	err = db.Init()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	address, ok := os.LookupEnv("ADDRESS")
	if !ok {
		address = ":8888"
	}

	srv := &http.Server{
		Addr: address,
		// Handler: Router(),
		Handler:           server.Router(),
		ReadHeaderTimeout: time.Second * 10,
	}

	p.server = srv
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
}

func (p *program) Start(s service.Service) error {
	// Start should not block.
	p.exit = make(chan struct{})
	p.wg.Add(1)
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	p.wg.Wait()
	return nil
}

func (p *program) run() {
	defer p.wg.Done()
	go func() {
		p.runServer()
	}()
	// Wait for stop signal
	<-p.exit
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Error during server shutdown: ")
	} else {
		logger.Info().Msg("Echo server shut down cleanly..")
	}
}

func main() {
	svcConfig := &service.Config{
		Name:        "WWNManager",
		DisplayName: "WWN Manager",
		Description: "WWN Manager Web service",
		Dependencies: []string{
			"Tcpip", // ensures TCP/IP is running
			"MongoDB",
		},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install", "uninstall", "start", "stop", "restart":
			err = service.Control(s, os.Args[1])
			if err != nil {
				logger.Fatal().Err(err).Msg(fmt.Sprintf("Failed to %s service", os.Args[1]))
			}
			logger.Info().Msg(fmt.Sprintf("Service %s successfully.", os.Args[1]))
			return
		}
	}

	if !service.Interactive() {
		// Running as a service
		err = s.Run()
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
		return
	}

	// Console mode for development/testing
	logger.Info().Msg("Running in console mode.")
	err = prg.Start(s)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Interrupt received, stopping service...")
	err = prg.Stop(s)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
