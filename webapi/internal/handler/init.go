package handler

import (
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
	"github.com/ttrnecka/wwn_identity/webapi/shared/utils"
)

var logger zerolog.Logger
var validate *validator.Validate

func init() {
	logging.LogLocation(filepath.Join(utils.BinaryOrBuildDir(), "logs"))
	logger = logging.SetupLogger("handler")
	validate = validator.New(validator.WithRequiredStructEnabled())
}
