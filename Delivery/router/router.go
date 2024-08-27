package router

import (
	"Loan_manager/Delivery/controller"
	"Loan_manager/infrastructure"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(userController *controller.UserController, loanController *controller.LoanController, logController *controller.LogController, tokenCollection *mongo.Collection) *gin.Engine {
	router := gin.Default()

	// Public routes (no authentication required)
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh", userController.RefreshToken)
	router.POST("/forgot-password", userController.ForgotPassword)
	router.GET("/reset/:token", userController.ResetPassword)
	// router.GET("/verify/:token", userController.Verify)

	// Authenticated user routes
	usersRoute := router.Group("/")
	usersRoute.Use(infrastructure.AuthMiddleware(tokenCollection)) // Apply authentication middleware

	// User management routes
	usersRoute.PUT("/update/:username", userController.UpdateUser)
	usersRoute.PUT("/change_password", userController.ChangePassword)
	usersRoute.POST("/logout", userController.Logout)

	// Loan management routes
	usersRoute.POST("/loans", loanController.ApplyLoan)
	usersRoute.GET("/loans/:id", loanController.ViewLoanStatus)

	// Admin routes (requires admin role)
	adminRoute := usersRoute.Group("/admin")
	adminRoute.Use(infrastructure.RoleMiddleware("admin")) // Apply admin role middleware

	adminRoute.DELETE("/delete/:username", userController.DeleteUser)

	// Admin loan management routes
	adminRoute.GET("/loans", loanController.ViewAllLoans)
	adminRoute.PATCH("/loans/:id/status", loanController.ApproveRejectLoan)
	adminRoute.DELETE("/loans/:id", loanController.DeleteLoan)

	// Admin system logs route
	adminRoute.GET("/logs", logController.ViewSystemLogs)
	log.Fatal(router.Run(":8080"))
	return router
}
