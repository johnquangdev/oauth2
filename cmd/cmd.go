package cmd

// @title oauth2 Golang Clean Architecture
// @version 1.0
// @description Backend oauth2 Golang Clean Architecture
// @host localhost:8080
// @BasePath /v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/cmd/sqlmigrate"
	"github.com/johnquangdev/oauth2/delivery"
	_ "github.com/johnquangdev/oauth2/docs"
	myMiddleware "github.com/johnquangdev/oauth2/middleware"
	"github.com/johnquangdev/oauth2/repository"
	"github.com/johnquangdev/oauth2/usecase"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
		log.Fatalf("Failed to load config: %v", err)
	}

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//connect database
	db, err := repository.ConnectPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// connect redis
	redis, err := repository.ConnectRedis(*config)
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}

	// new repository
	repo := repository.NewRepository(db, redis)

	// sql migrate
	err = sqlmigrate.RunSqlMigrate(*config, db)
	if err != nil {
		log.Fatalf("Failed to run SQL migration: %v", err)
	}

	// register useCase
	u, err := usecase.NewUseCase(*config, repo, redis)
	if err != nil {
		log.Fatalf("Failed to register usecase: %v", err)
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
