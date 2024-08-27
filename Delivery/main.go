package main

import (
	"Loan_manager/Delivery/controller"
	"Loan_manager/Delivery/router"
	"Loan_manager/Repository"
	"Loan_manager/Usecases"
	"Loan_manager/infrastructure"

	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URL")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	userDatabase := client.Database("Blog_management")
	userCollection := userDatabase.Collection("User")
	tokenCollection := userDatabase.Collection("Token")
	logCollection := userDatabase.Collection("Logs") // Collection for system logs

	userRepository := Repository.NewUserRepository(userCollection, tokenCollection)
	loanRepository := Repository.NewLoanRepository(userCollection)
	logRepository := Repository.NewLogRepository(logCollection) // Create log repository

	emailService := infrastructure.NewEmailService()

	userUsecase := Usecases.NewUserUsecase(userRepository, emailService)
	loanUsecase := Usecases.NewLoanUsecase(loanRepository)
	logUsecase := Usecases.NewLogUsecase(logRepository) // Create log usecase

	userController := controller.NewUserController(userUsecase)
	loanController := controller.NewLoanController(loanUsecase)
	logController := controller.NewLogController(logUsecase) // Create log controller

	router := router.SetupRouter(userController, loanController, logController, tokenCollection)
	log.Fatal(router.Run(":8080"))
}
