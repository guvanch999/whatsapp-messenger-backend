package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/repo"
	"github.com/medium-messenger/messenger-backend/internal/modules/contacts/service"
)

func InitUserContactsRouter(server *cmd.Server) {
	contactsRepository := repo.NewUserContactRepository(server.Database)
	contactsService := service.NewUserContactsService(contactsRepository)
	contactsHandler := handler.NewUserContactsHandler(contactsService)

	g := server.Echo.Group("v1/user-contacts")
	g.Use(middleware.AuthMiddleware(server.Supabase, server.Database))

	g.GET("", contactsHandler.GetMyContacts)
	g.GET("/all", contactsHandler.GetAllContacts, middleware.CheckAdminMiddleware)
	g.GET("/:guid", contactsHandler.GetContactDetail)
	g.POST("", contactsHandler.AddContact)
	g.POST("/list", contactsHandler.AddListOfContacts)
	g.POST("/validate", contactsHandler.ValidateNumber)
	g.PUT("/:guid", contactsHandler.UpdateContactDetail)
	g.DELETE("/:guid", contactsHandler.DeleteContactDetail)
}
