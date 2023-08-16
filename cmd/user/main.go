package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wlrudi19/elastic-engine/app/user/api"
	"github.com/wlrudi19/elastic-engine/app/user/repository"
	"github.com/wlrudi19/elastic-engine/app/user/service"
	"github.com/wlrudi19/elastic-engine/config"
)

func main() {
	loadConfig := config.LoanConfig()
	connDB, err := config.ConnectConfig(loadConfig.Database)

	if err != nil {
		log.Fatalf("error connecting to postgres :%v", err)
		return
	}
	defer connDB.Close()

	fmt.Println("ELASTIC ENGINE PROJECT")
	log.Printf("connected to postgres successfulyy")

	userRepository := repository.NewUserRepository()
	userLogic := service.NewUserLogic(userRepository, connDB)
	userHandler := api.NewUserHandler(userLogic)
	userRouter := api.NewProductRouter(userHandler)

	server := http.Server{
		Addr:    "localhost:7655",
		Handler: userRouter,
	}

	fmt.Println("starting server on port 7654...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
