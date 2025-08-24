package server

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/handler"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	mid "github.com/ttrnecka/wwn_identity/webapi/server/middleware"
)

func origin() string {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(3, pcs) // skip 3 frames (Callers, origin, handler)
	frames := runtime.CallersFrames(pcs[:n])
	if frame, more := frames.Next(); more {
		return fmt.Sprintf("%s:%d", frame.File, frame.Line)
	}
	return "unknown"
}

func Router() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(mid.SessionManager())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// Capture request details
		req := c.Request()
		path := req.URL.Path

		// Log full trace with timestamp and route
		fmt.Printf("\n=== ERROR ===\nTime: %s\nPath: %s\nError: %v\nStack Trace:\n%s\n",
			time.Now().Format(time.RFC3339), path, err, origin())

		// Respond with generic error to client
		if !c.Response().Committed {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error",
			})
		}
	}

	// db layer
	users := entity.Users()
	fcWWNEntries := entity.FCWWNEntries()
	rule := entity.Rules()

	// repositories
	usersRepo := repository.NewUserRepository(users)
	fcWWNEntryRepo := repository.NewFCWWNEntryRepository(fcWWNEntries)
	ruleRepo := repository.NewRuleRepository(rule)

	// services
	userSvc := service.NewUserService(usersRepo)
	fcwWWNEntrySvc := service.NewFCWWNEntryService(fcWWNEntryRepo)
	ruleSvc := service.NewRuleService(ruleRepo)

	//handlers
	userHandler := handler.NewUserHandler(userSvc)
	fcWWNEntryHandler := handler.NewFCWWNEntryHandler(fcwWWNEntrySvc, ruleSvc)
	ruleHandler := handler.NewRuleHandler(ruleSvc, fcwWWNEntrySvc)

	e.POST("/api/login", userHandler.LoginUser)
	e.GET("/api/logout", userHandler.LogoutUser)
	e.GET("/api/user", userHandler.User, mid.AuthMiddleware)

	api := e.Group("/api/v1", mid.AuthMiddleware)

	//rules amd entries

	api.POST("/import", fcWWNEntryHandler.ImportHandler)
	api.GET("/customers", fcWWNEntryHandler.ListCustomers)
	api.GET("/rules", ruleHandler.Rules)
	api.GET("/rules/export", ruleHandler.ExportRules)
	api.GET("/customers/:name/rules", ruleHandler.GetRules)
	api.POST("/customers/:name/rules", ruleHandler.CreateUpdateRule)
	api.DELETE("/customers/:name/rules/:id", ruleHandler.DeleteRule)
	api.POST("/entries/:id/reconcile", ruleHandler.SetupAndApplyReconcileRules)
	api.GET("/customers/:name/entries", fcWWNEntryHandler.FCWWNEntries)

	return e

}
