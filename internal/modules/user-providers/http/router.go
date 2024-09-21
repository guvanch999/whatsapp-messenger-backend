package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/repo"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/service"
)

func InitUserProvidersRouter(server *cmd.Server) {
	userProviderRepository := repo.NewUserProviderRepository(server.Database)
	userProviderService := service.NewUserProviderService(
		userProviderRepository,
		server.SecretManagerClient,
		server.Config,
	)
	userProviderHandler := handler.NewUserProviderHandler(userProviderService)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)
	g := server.Echo.Group("v1/user-providers", authMiddleware)

	g.GET("", userProviderHandler.GetUserProviders)
	g.GET("/all", userProviderHandler.GetAllProviders, middleware.CheckAdminMiddleware)
	g.GET("/:guid", userProviderHandler.GetDetail)
	g.POST("", userProviderHandler.CreateProvider)
	g.DELETE("/:guid", userProviderHandler.DeleteProvider)

}
