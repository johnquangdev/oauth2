package cmd

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/delivery"
	"github.com/johnquangdev/oauth2/usecase"
	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New()
	// register validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	// register useCase
	useCase, err := usecase.NewUseCase()
	if err != nil {
		log.Fatalf("register usecase err by: %v", err)
	}
	g := e.Group("/v1")
	// register router
	err = delivery.NewDelivery(useCase, g, validate)
	if err != nil {
		log.Fatalf("register router failed by error: %v", err)
	}
	// run server
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatalf("server startup failed due to error: %v", err)
	}

}
