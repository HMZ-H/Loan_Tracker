package router

import (
	"Loan_manager/Delivery/controller"
	"Loan_manager/infrastructure"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(userController *controller.UserController, tokenCollection *mongo.Collection) *gin.Engine {
	router := gin.Default()

	// Public routes (no authentication required)
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh", userController.RefreshToken)
	router.POST("/forgot-password", userController.ForgotPassword)
	router.GET("/reset/:token", userController.ResetPassword)
	router.GET("/verify/:token", userController.Verify)
	// router.GET("/glogin", userController.LoginHandler)
	// router.GET("/auth/google/callback", userController.CallbackHandler)

	// Authenticated user routes
	usersRoute := router.Group("/")
	usersRoute.Use(infrastructure.AuthMiddleware(tokenCollection)) // Apply authentication middleware

	// User management routes
	usersRoute.PUT("/update/:username", userController.UpdateUser)
	usersRoute.PUT("/change_password", userController.ChangePassword)
	usersRoute.POST("/logout", userController.Logout)

	// Admin routes (requires admin role)
	adminRoute := usersRoute.Group("/")
	adminRoute.Use(infrastructure.RoleMiddleware("admin")) // Apply admin role middleware

	adminRoute.DELETE("/delete/:username", userController.DeleteUser)

	return router
}
