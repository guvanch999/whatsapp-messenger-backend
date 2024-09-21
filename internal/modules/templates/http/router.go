package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/service"
)

// InitTemplatesRouter todo user own provider
func InitTemplatesRouter(server *cmd.Server) {
	templateRepository := repository.NewTemplateRepository(server.Database)
	templateService := service.NewTemplateService(server.Database, server.SecretManagerClient, templateRepository)
	templateHandler := handler.NewTemplateHandler(templateService)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)
	g := server.Echo.Group("v1/templates")

	g.GET("", templateHandler.GetMyTemplates, authMiddleware)
	g.GET("/all", templateHandler.GetAllTemplates, authMiddleware, middleware.CheckAdminMiddleware)
	g.GET("/:guid", templateHandler.GetDetail, authMiddleware)
	g.POST("", templateHandler.CreateTemplate, authMiddleware)
	g.DELETE("/:guid", templateHandler.DeleteTemplate, authMiddleware)

	g.POST("/approve/:guid", templateHandler.ApproveTemplate, authMiddleware)

	g.GET("/sync", templateHandler.Sync)

}
