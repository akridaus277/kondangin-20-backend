package routes

import (
	"kondangin-backend/internal/handler"
	"kondangin-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	r.GET("/hello-world", handler.HelloWorld)

	r.POST("/register", userHandler.RegisterUser)

	r.POST("/login", userHandler.LoginUser)

	r.POST("/encrypt-password", userHandler.EncryptPassword)

	r.GET("/verify", userHandler.VerifyEmail)

	r.POST("/resend-verification", userHandler.ResendVerificationEmail)

	r.POST("/forgot-password", userHandler.ForgotPassword)

	r.POST("/reset-password", userHandler.ResetPassword)

}

func InvitationDashboardRoutes(r *gin.Engine, invitationDashboardHandler *handler.InvitationDashboardHandler) {
	inv := r.Group("/invitation-dashboard")

	// Kalau ingin route dengan JWT auth:
	auth := inv.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/create", invitationDashboardHandler.CreateInvitation)
		auth.POST("/get-data-json", invitationDashboardHandler.GetInvitationData)
		auth.POST("/add-permission", invitationDashboardHandler.AddInvitationPermission)
	}
}

func InvitationGuestRoutes(r *gin.Engine, invitationGuestHandler *handler.InvitationGuestHandler) {
	inv := r.Group("/invitation-guest")

	// Route tanpa auth
	inv.POST("/get-invitation", invitationGuestHandler.GetInvitationData)

}
