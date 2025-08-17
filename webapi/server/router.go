package server

import (
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

func Router() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(mid.SessionManager())

	// db layer
	users := entity.Users()
	fcEntries := entity.FCEntries()
	rule := entity.Rules()

	// repositories
	usersRepo := repository.NewUserRepository(users)
	fcEntryRepo := repository.NewFCEntryRepository(fcEntries)
	ruleRepo := repository.NewRuleRepository(rule)

	// services
	userSvc := service.NewUserService(usersRepo)
	fcEntrySvc := service.NewFCEntryService(fcEntryRepo)
	ruleSvc := service.NewRuleService(ruleRepo)

	//handlers
	userHandler := handler.NewUserHandler(userSvc)
	fcEntryHandler := handler.NewFCEntryHandler(fcEntrySvc, ruleSvc)
	ruleHandler := handler.NewRuleHandler(ruleSvc)

	e.POST("/api/login", userHandler.LoginUser)
	e.GET("/api/logout", userHandler.LogoutUser)
	e.GET("/api/user", userHandler.User, mid.AuthMiddleware)

	api := e.Group("/api/v1", mid.AuthMiddleware)

	//rules amd entries

	api.POST("/import", fcEntryHandler.ImportHandler)
	api.GET("/customers", fcEntryHandler.ListCustomers)
	api.GET("/customers/:name/rules", ruleHandler.GetRules)
	api.POST("/customers/:name/rules", ruleHandler.CreateUpdateRule)
	api.DELETE("/customers/:name/rules/:id", ruleHandler.DeleteRule)
	api.GET("/customers/:name/entries", fcEntryHandler.FCEntries)

	return e

}
