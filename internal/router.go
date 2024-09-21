package internal

import (
	"github.com/medium-messenger/messenger-backend/cmd"
	. "github.com/medium-messenger/messenger-backend/internal/modules/api-keys/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/auth/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contact-list/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/contacts/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/messaging/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/organization/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/templates/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/user-providers/http"
	. "github.com/medium-messenger/messenger-backend/internal/modules/users/http"
)

func InitRouters(server *cmd.Server) {
	InitAuthRouter(server)
	InitUsersRouter(server)
	InitUserContactsRouter(server)
	InitContactListRouter(server)
	InitTemplatesRouter(server)
	InitOrganizationRouter(server)
	InitUserProvidersRouter(server)
	InitMessagingRouter(server)
	InitApiKeysRouter(server)
}
