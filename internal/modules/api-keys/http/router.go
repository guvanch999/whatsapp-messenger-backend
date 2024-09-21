package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/repo"
	"github.com/medium-messenger/messenger-backend/internal/modules/api-keys/service"
)

func InitApiKeysRouter(server *cmd.Server) {
	apiKeyRepository := repo.NewApiKeyRepository(server.Database)
	apiKeyService := service.NewApiKeysService(server.Config, apiKeyRepository)
	apiKeyHandler := handler.NewApiKeysHandler(apiKeyService)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)
	g := server.Echo.Group("v1/api-keys", authMiddleware)

	g.GET("", apiKeyHandler.GetUserApiKeys)
	g.GET("/:guid", apiKeyHandler.GetDetail)
	g.GET("/value/:guid", apiKeyHandler.GetApiKeyValue)
	g.POST("", apiKeyHandler.AddApiKey)
	g.DELETE("/:guid", apiKeyHandler.DeleteApiKey)
}
