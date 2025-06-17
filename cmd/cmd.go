package cmd

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/delivery"
	myMiddleware "github.com/johnquangdev/oauth2/middleware"
	"github.com/johnquangdev/oauth2/repository"
	google "github.com/johnquangdev/oauth2/service/google"
	"github.com/johnquangdev/oauth2/usecase"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {
	e := echo.New()
	// middlerware
	e.Use(middleware.Logger())
	// register validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	//load config
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("load config failed by err: %v", err)
	}
	//connect database
	db, err := repository.ConnectPostgres(*config)
	if err != nil {
		log.Fatalf("connect database failed by err: %v", err)
	}
	// connect redis
	redis, err := repository.ConnectRedis(*config)
	if err != nil {
		log.Fatalf("connect redis failed by err: %v", err)
	}
	// new repository
	repo := repository.NewRepository(db, redis)
	// new service oauth2
	oauth2Service := google.NewGoogleOAuthService()
	// register useCase
	u, err := usecase.NewUseCase(repo, db, oauth2Service, redis)
	if err != nil {
		log.Fatalf("register usecase err by: %v", err)
	}
	middleware := myMiddleware.NewMiddleware(*config, repo)
	// register router
	g := e.Group("/v1")
	delivery.NewDelivery(u, g, validate, *config, middleware)
	// run server
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatalf("server startup failed due to error: %v", err)
	}
}
