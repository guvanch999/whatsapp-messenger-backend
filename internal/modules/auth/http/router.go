package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/auth/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/auth/service"
)

func InitAuthRouter(server *cmd.Server) {
	authService := service.NewAuthService(server.Supabase)
	authHandler := handler.NewAuthHandler(authService, server.Config)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)

	g := server.Echo.Group("v1/auth")
	g.POST("/register", authHandler.RegisterHandler)
	g.POST("/login", authHandler.LoginHandler)
	g.POST("/refresh", authHandler.RefreshHandler, authMiddleware)
	g.POST("/change-password", authHandler.ChangePassword, authMiddleware)
	g.GET("/me", authHandler.GetMe, authMiddleware)
}
