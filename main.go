package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/peang/amartha-loan-service/configs"
	"github.com/peang/amartha-loan-service/handlers"
	middlewares "github.com/peang/amartha-loan-service/middlewares"
	"github.com/peang/amartha-loan-service/repositories"
	"github.com/peang/amartha-loan-service/services/file_services"
	"github.com/peang/amartha-loan-service/usecases"
)

func main() {
	conf := configs.LoadConfig()

	db := configs.LoadDatabase(conf)
	defer db.Close()

	enfocer, err := configs.NewCasbinEnfocer()
	if err != nil {
		panic(err)
	}

	middleware := middlewares.NewMiddleware(enfocer)

	e := echo.New()

	// Register Repositories
	userRepository := repositories.NewUserRepository(db)
	approvalRepository := repositories.NewApprovalRepository(db)
	loanRepository := repositories.NewLoanRepository(db, approvalRepository)
	investmentRepository := repositories.NewInvestmentRepository(db, loanRepository)

	// Register Services
	fileService := file_services.NewLocalFileService()

	// Register Usecases
	loanUsecase := usecases.NewLoanUsecase(loanRepository, investmentRepository, fileService)

	handlers.NewAuthHandler(e, userRepository)
	handlers.NewLoanHandler(e, middleware, loanUsecase)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
