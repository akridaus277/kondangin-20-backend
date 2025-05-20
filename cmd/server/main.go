package main

import (
	"fmt"
	"kondangin-backend/config"
	"kondangin-backend/internal/handler"
	invitationWeddingHandler "kondangin-backend/internal/handler/invitation/wedding"
	models "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	routes "kondangin-backend/internal/route"
	"kondangin-backend/internal/service"
	invitationWeddingService "kondangin-backend/internal/service/invitation/wedding"
	invitationWeddingValidatorService "kondangin-backend/internal/service/invitation/wedding/validator"
	validatorService "kondangin-backend/internal/service/validator"
	"os"
	"strings"

	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	config.LoadConfig()

	config.ConnectDatabase()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Invitation{},
		&models.InvitationPermission{},
		&models.EventType{},
		&models.InvitationTemplate{},
		&models.Package{}) // <- ini migrasi db
	db := config.GetDB()
	// Init router
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		fmt.Println("Origin:", c.Request.Header.Get("Origin"))
		c.Next()
	})

	// CORS Configuration
	frontendURL := os.Getenv("FRONTEND_URL") // e.g., http://localhost:3000 or https://example.com

	u, err := url.Parse(frontendURL)
	if err != nil {
		panic("Invalid FRONTEND_URL: " + frontendURL)
	}
	host := u.Hostname()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			parsedOrigin, err := url.Parse(origin)
			if err != nil {
				return false
			}
			originHost := parsedOrigin.Hostname()

			// exact match
			if origin == frontendURL {
				return true
			}

			// match root domain
			if originHost == host {
				return true
			}

			// match subdomains of root domain
			if strings.HasSuffix(originHost, "."+host) {
				return true
			}

			return false
		},
	}

	router.Use(cors.New(corsConfig))

	// Repository
	userRepo := repository.NewUserRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)
	invitationPermissionRepo := repository.NewInvitationPermissionRepository(db)
	eventTypeRepo := repository.NewEventTypeRepository(db)
	// packageRepo := repository.NewPackageRepository(db)
	// Model Service
	userService := service.NewUserService(userRepo)
	invitationService := service.NewInvitationService(invitationRepo)
	invitationPermissionService := service.NewInvitationPermissionService(invitationPermissionRepo)
	eventTypeService := service.NewEventTypeService(eventTypeRepo)
	// Validator Service
	memberDashboardValidatorService := validatorService.NewMemberDashboardValidator(invitationRepo, invitationPermissionRepo, eventTypeRepo)
	invitationWeddingValidatorService := invitationWeddingValidatorService.NewWeddingValidator(invitationRepo, invitationPermissionService, eventTypeRepo)

	// Feature Service
	invitationDashboardService := service.NewInvitationDashboardService(invitationRepo, invitationPermissionRepo)
	memberDashboardService := service.NewMemberDashboardService(invitationService, invitationPermissionService, eventTypeService, memberDashboardValidatorService)
	invitationGuestService := service.NewInvitationGuestService(invitationRepo)
	invitationWeddingService := invitationWeddingService.NewWeddingService(invitationWeddingValidatorService, invitationService)
	authService := service.NewAuthService(userRepo)
	// Handler
	authHandler := handler.NewAuthHandler(authService)
	invitationDashboardHandler := handler.NewInvitationDashboardHandler(invitationDashboardService)
	invitationGuestHandler := handler.NewInvitationGuestHandler(invitationGuestService)
	memberDashboardHandler := handler.NewMemberDashboardHandler(memberDashboardService)
	invitationWeddingHandler := invitationWeddingHandler.NewInvitationWeddingHandler(invitationWeddingService)

	// Expose folder uploads ke public
	router.Static("/uploads", "./uploads")
	// Register routes
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	routes.AuthRoutes(router, authHandler)
	routes.InvitationDashboardRoutes(router, userService, invitationDashboardHandler, invitationWeddingHandler)
	routes.InvitationGuestRoutes(router, invitationGuestHandler)
	routes.MemberDashboardRoutes(router, userService, memberDashboardHandler)

	// Start server
	router.Run(":8080")
}
