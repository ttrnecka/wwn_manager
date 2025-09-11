package server

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ttrnecka/wwn_identity/webapi/db"
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/handler"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
	"github.com/ttrnecka/wwn_identity/webapi/internal/service"
	mid "github.com/ttrnecka/wwn_identity/webapi/server/middleware"
)

func Router() *echo.Echo {
	e := echo.New()

	e.Use(mid.RequestBodyCaptureMiddleware())
	e.Use(mid.ResponseBodyCaptureMiddleware())
	e.Use(mid.RequestLogger(logger))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(mid.SessionManager())

	// db layer
	users := entity.Users(db.Database())
	fcWWNEntries := entity.FCWWNEntries(db.Database())
	rule := entity.Rules(db.Database())
	snapshot := entity.Snapshots(db.Database())

	// repositories
	usersRepo := repository.NewUserRepository(users)
	fcWWNEntryRepo := repository.NewFCWWNEntryRepository(fcWWNEntries)
	ruleRepo := repository.NewRuleRepository(rule)
	snapshotRepo := repository.NewSnapshotRepository(snapshot)

	// services
	userSvc := service.NewUserService(usersRepo)
	fcwWWNEntrySvc := service.NewFCWWNEntryService(fcWWNEntryRepo)
	ruleSvc := service.NewRuleService(ruleRepo)
	snapshotSvc := service.NewSnapshotService(snapshotRepo, fcWWNEntryRepo)

	//handlers
	userHandler := handler.NewUserHandler(userSvc)
	fcWWNEntryHandler := handler.NewFCWWNEntryHandler(fcwWWNEntrySvc, ruleSvc)
	ruleHandler := handler.NewRuleHandler(ruleSvc, fcwWWNEntrySvc)
	snapshotHandler := handler.NewSnapshotHandler(snapshotSvc, fcwWWNEntrySvc)

	e.POST("/api/login", userHandler.LoginUser)
	e.GET("/api/logout", userHandler.LogoutUser)
	e.GET("/api/user", userHandler.User, mid.AuthMiddleware)

	api := e.Group("/api/v1", mid.AuthMiddleware)

	//rules amd entries

	api.POST("/import", fcWWNEntryHandler.ImportHandler)
	api.GET("/customers", fcWWNEntryHandler.ListCustomers)
	api.GET("/rules", ruleHandler.Rules)
	api.GET("/rules/export", ruleHandler.ExportRules)
	api.POST("/rules/import", ruleHandler.ImportHandler)
	api.POST("/rules/apply", ruleHandler.ApplyRules)
	api.GET("/entries/export/map", fcWWNEntryHandler.ExportHostWWNMap)
	api.GET("/entries/export/override", fcWWNEntryHandler.ExportCustomerWWNMap)
	api.GET("/customers/:name/rules", ruleHandler.GetRules)
	api.POST("/customers/:name/rules", ruleHandler.CreateUpdateRule)
	api.DELETE("/customers/:name/rules/:id", ruleHandler.DeleteRule)
	api.POST("/entries/:id/reconcile", ruleHandler.SetupAndApplyReconcileRules)
	api.POST("/entries/:id/softdelete", fcWWNEntryHandler.SoftDeleteFCWWNEntry)
	api.POST("/entries/:id/restore", fcWWNEntryHandler.RestoreFCWWNEntry)
	api.GET("/customers/:name/entries", fcWWNEntryHandler.FCWWNEntries)
	api.GET("/snapshots", snapshotHandler.Snapshots)
	api.GET("/snapshots/:id", snapshotHandler.GetSnapshotEntries)
	api.POST("/snapshots", snapshotHandler.CreateSnapshot)

	return e

}
