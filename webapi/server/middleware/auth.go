package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/v5/session"
	"github.com/labstack/echo/v5"
)

func SessionManager() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			sess, err := session.Get(SessionStore, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session").Wrap(err)
			}
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   3600 * 1,
				HttpOnly: true,
			}
			sess.Values["_refresh"] = time.Now().UnixNano()
			c.Set(SessionStore, sess)

			// Ensure the session is always saved *before* headers are sent
			response, err := echo.UnwrapResponse(c.Response())
			if err != nil {
				return fmt.Errorf("unwrap Echo response: %w", err)
			}
			response.Before(func() {
				if err := sess.Save(c.Request(), c.Response()); err != nil {
					c.Logger().Error("failed to save session", "error", err)
				}
			})

			if err := next(c); err != nil {
				return err
			}
			return nil
		}
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		sess := Session(c)
		if sess.Values["user"] != nil {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}
}
