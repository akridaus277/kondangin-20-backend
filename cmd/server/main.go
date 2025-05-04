package main

import (
	"fmt"
	"kondangin-backend/config"
	"kondangin-backend/internal/handler"
	models "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"
	routes "kondangin-backend/internal/route"
	"kondangin-backend/internal/service"
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
	config.DB.AutoMigrate(&models.User{}, &models.Invitation{}) // <- ini migrasi tabel User
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	invitationRepo := repository.NewInvitationRepository(db)
	invitationPermissionRepo := repository.NewInvitationPermissionRepository(db)
	invitationDashboardService := service.NewInvitationDashboardService(invitationRepo, invitationPermissionRepo)
	invitationDashboardHandler := handler.NewInvitationDashboardHandler(invitationDashboardService)
	invitationGuestService := service.NewInvitationGuestService(invitationRepo)
	invitationGuestHandler := handler.NewInvitationGuestHandler(invitationGuestService)

	// Register routes
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	routes.UserRoutes(router, userHandler)
	routes.InvitationDashboardRoutes(router, invitationDashboardHandler)
	routes.InvitationGuestRoutes(router, invitationGuestHandler)

	// Start server
	router.Run(":8080")
}
