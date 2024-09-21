package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/medium-messenger/messenger-backend/cmd"
	"github.com/medium-messenger/messenger-backend/docs"
	_ "github.com/medium-messenger/messenger-backend/docs"
	"github.com/medium-messenger/messenger-backend/internal"
	"github.com/medium-messenger/messenger-backend/internal/middleware"
	"github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

//			@title						Whatsapp messenger backend API
//			@version					1.0
//			@contact.name				API Support
//			@contact.email				testemail@gmail.com
//			@BasePath					/v1
//			@securityDefinitions.apikey	Bearer
//			@in							header
//			@name						Authorization
//	 		@description				Bearer token is required to access the Whatsapp messenger backend API. If you don't have a token, please register in the system first
//			@securityDefinitions.apikey	X-API-KEY
//			@in							header
//			@name						X-API-KEY
//			@description				API key is required to access the Whatsapp messenger backend API. If you don't have a token, please register in the system first
func main() {
	server := cmd.NewServer()

	server.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	middleware.RegisterMiddlewares(server)
	internal.InitRouters(server)

	if server.Config.Environment == "development" {
		docs.SwaggerInfo.Schemes = []string{"http"}
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", server.Config.Port)
	} else {
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.Host = strings.ReplaceAll(strings.ReplaceAll(server.Config.AppUrl, "https://", ""), "/v1", "")
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := server.Echo.Start(fmt.Sprintf(":%v", server.Config.Port)); err != nil && !errors.Is(
			err,
			http.ErrServerClosed,
		) {
			server.Echo.Logger.Fatal("shutting down the server. error: " + err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-signalCtx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Echo.Shutdown(ctx); err != nil {
		server.Echo.Logger.Fatal(err)
	}
}
