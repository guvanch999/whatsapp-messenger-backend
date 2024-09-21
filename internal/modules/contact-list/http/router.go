package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/contact-list/service"
)

func InitContactListRouter(server *cmd.Server) {
	contactListRepository := repository.NewContactListRepository(server.Database)
	contactListService := service.NewContactListService(contactListRepository)
	contactListHandler := handler.NewContactListHandler(contactListService)

	authMiddleware := middleware.AuthMiddleware(server.Supabase, server.Database)
	g := server.Echo.Group("v1/contact-list", authMiddleware)

	g.GET("", contactListHandler.GetUserContactLists)
	g.GET("/all", contactListHandler.GetAllContactList, middleware.CheckAdminMiddleware)
	g.GET("/:guid", contactListHandler.GetDetail)
	g.POST("", contactListHandler.CreateContactList)
	g.POST("/add-contact/:guid", contactListHandler.AddContactToList)
	g.POST("/remove-contact/:guid", contactListHandler.DeleteContactFromList)
	g.PUT("/name/:guid", contactListHandler.UpdateContactListName)
	g.DELETE("/:guid", contactListHandler.DeleteContactList)

}
