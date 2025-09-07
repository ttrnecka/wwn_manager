package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionManager() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			sess, err := session.Get(SESSION_STORE, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session:", err)
			}
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   3600 * 1,
				HttpOnly: true,
			}
			sess.Values["_refresh"] = time.Now().UnixNano()
			c.Set(SESSION_STORE, sess)

			// Ensure the session is always saved *before* headers are sent
			c.Response().Before(func() {
				if err := sess.Save(c.Request(), c.Response()); err != nil {
					c.Logger().Errorf("failed to save session: %v", err)
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
	return func(c echo.Context) error {
		sess := Session(c)
		if sess.Values["user"] != nil {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}
