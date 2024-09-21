package http

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/handler"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/repository"
	"github.com/medium-messenger/messenger-backend/internal/modules/users/service"
)

func InitUsersRouter(server *cmd.Server) {
	usersRepository := repository.NewUserRepository(server.Database)
	usersService := service.NewUserService(usersRepository)
	userHandler := handler.NewUserHandler(usersService)

	adminG := server.Echo.Group("v1/users")
	adminG.Use(middleware.AuthMiddleware(server.Supabase, server.Database))
	adminG.Use(middleware.CheckAdminMiddleware)
	adminG.GET("", userHandler.GetAllUsers)
	adminG.PUT("/:guid", userHandler.ChangeUserInfo)
}
