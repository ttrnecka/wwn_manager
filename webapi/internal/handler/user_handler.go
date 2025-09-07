package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"github.com/ttrnecka/wwn_identity/webapi/internal/mapper"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	"github.com/ttrnecka/wwn_identity/webapi/server/middleware"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	// func NewUserHandler() *UserHandler {
	return &UserHandler{s}
}

func (e *UserHandler) LoginUser(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := e.service.GetByName(c.Request().Context(), username)

	if err != nil {
		logger.Error().Err(err).Msg("")
		if errors.Is(err, cdb.ErrNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}
		return errorWithInternal(http.StatusInternalServerError, "Unexpected login error: ", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}
	userDTO := mapper.ToUserDTO(*user)

	sess := middleware.Session(c)
	sess.Values["user"] = userDTO
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful", "user": userDTO.Username})
}

func (e *UserHandler) LogoutUser(c echo.Context) error {
	sess := middleware.Session(c)
	sess.Options.MaxAge = -1
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out"})
}

func (e *UserHandler) User(c echo.Context) error {
	sess := middleware.Session(c)
	user, ok := sess.Values["user"].(dto.UserDTO)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Incorrect user in session, assertion failed")
	}
	return c.JSON(http.StatusOK, user)
}
