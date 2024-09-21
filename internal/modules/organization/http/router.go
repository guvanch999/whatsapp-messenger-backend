package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/service"
)

func InitOrganizationRouter(server *cmd.Server) {
	organizationRepository := repository.NewOrganizationRepository(server.Database)
	organizationService := service.NewOrganizationService(organizationRepository)
	organizationHandler := handler.NewOrganizationHandler(organizationService)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)
	g := server.Echo.Group("v1/organizations", authMiddleware)

	g.GET("", organizationHandler.GetUserOrganizations)
	g.GET("/all", organizationHandler.GetUserOrganizations, middleware.CheckAdminMiddleware)
	g.GET("/:guid", organizationHandler.GetDetail)
	g.POST("", organizationHandler.CreateOrganization)
	g.PUT("/:guid", organizationHandler.UpdateOrganization)
	g.DELETE("/:guid", organizationHandler.DeleteOrganization)
}
