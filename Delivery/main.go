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
	userRepository := Repository.NewUserRepository(userCollection, tokenCollection)

	// Initialize the Email Service
	emailService := infrastructure.NewEmailService()

	// Initialize the User Usecase with the User Repository and Email Service
	userUsecase := Usecases.NewUserUsecase(userRepository, emailService)
	userController := controller.NewUserController(userUsecase)

	router := router.SetupRouter(userController, tokenCollection)
	log.Fatal(router.Run(":8080"))

}
