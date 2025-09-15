package main

import (
	"encoding/gob"
	"net/http"
	"time"

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

func main() {
	gob.Register(dto.UserDTO{})

	// db
	err := db.Init()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	srv := &http.Server{
		Addr: ":8888",
		// Handler: Router(),
		Handler:           server.Router(),
		ReadHeaderTimeout: time.Second * 10,
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
}
