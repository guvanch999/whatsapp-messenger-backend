package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	repository2 "github.com/medium-messenger/messenger-backend/internal/modules/contact-list/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/messaging/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/messaging/service"
	"github.com/medium-messenger/messenger-backend/internal/modules/templates/repository"
	service2 "github.com/medium-messenger/messenger-backend/internal/modules/templates/service"
)

func InitMessagingRouter(server *cmd.Server) {
	templateRepository := repository.NewTemplateRepository(server.Database)
	templateService := service2.NewTemplateService(server.Database, server.SecretManagerClient, templateRepository)

	contactListRepository := repository2.NewContactListRepository(server.Database)

	messageService := service.NewMessageService(
		server.Database,
		server.SecretManagerClient,
		templateService,
		contactListRepository,
	)
	messageHandler := handler.NewMessageHandler(messageService)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)
	g := server.Echo.Group("v1/messages", authMiddleware)

	g.POST("", messageHandler.SendMessage)
	g.POST("/to-list", messageHandler.SendMessageList)

}
