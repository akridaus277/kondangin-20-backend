package routes

import (
	"kondangin-backend/internal/handler"
	"kondangin-backend/internal/middleware"
	"kondangin-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	r.GET("/hello-world", handler.HelloWorld)

	r.POST("/register", authHandler.RegisterUser)

	r.POST("/login", authHandler.LoginUser)

	r.POST("/encrypt-password", authHandler.EncryptPassword)

	r.GET("/verify", authHandler.VerifyEmail)

	r.POST("/resend-verification", authHandler.ResendVerificationEmail)

	r.POST("/forgot-password", authHandler.ForgotPassword)

	r.POST("/reset-password", authHandler.ResetPassword)

}

func InvitationDashboardRoutes(r *gin.Engine, userService service.UserService, invitationDashboardHandler *handler.InvitationDashboardHandler) {
	inv := r.Group("/invitation-dashboard")

	// Kalau ingin route dengan JWT auth:
	auth := inv.Group("/")
	auth.Use(middleware.JWTAuthMiddleware(userService))
	{
		auth.POST("/get-data-json", invitationDashboardHandler.GetInvitationData)
		auth.POST("/add-permission", invitationDashboardHandler.AddInvitationPermission)
	}
}

func InvitationGuestRoutes(r *gin.Engine, invitationGuestHandler *handler.InvitationGuestHandler) {
	inv := r.Group("/invitation-guest")

	// Route tanpa auth
	inv.POST("/get-invitation", invitationGuestHandler.GetInvitationData)

}

func MemberDashboardRoutes(r *gin.Engine, userService service.UserService, memberDashboardHandler *handler.MemberDashboardHandler) {
	inv := r.Group("/member-dashboard")

	// Kalau ingin route dengan JWT auth:
	auth := inv.Group("/")
	auth.Use(middleware.JWTAuthMiddleware(userService))
	{
		auth.POST("/create-invitation", memberDashboardHandler.CreateInvitation)
		auth.GET("/get-invitation", memberDashboardHandler.GetInvitation)
		auth.GET("/lov-event-type", memberDashboardHandler.GetLovEventType)
	}
}
