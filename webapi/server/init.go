package server

import (
	"path/filepath"

	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
	"github.com/ttrnecka/wwn_identity/webapi/shared/utils"
)

var logger zerolog.Logger

func init() {
	logging.LogLocation(filepath.Join(utils.BinaryOrBuildDir(), "logs"))
	logger = logging.SetupLogger("http")
}

const SESSION_STORE = "agentpoc"
