package server

import (
	"Trivia_Quiz/pkg/authenticate"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ApiV1 = "/api/v1"

type GinServer struct {
	Modules    []Module
	HttpServer *http.Server
}

// NewGinServer Create a new server instance
func NewGinServer(modules []Module, auth *authenticate.Auth, ctx *context.Context) *GinServer {
	router := gin.Default()
	// Register routes of modules
	unAuth := router.Group(ApiV1).Use(auth.GinMiddleware(ctx))
	for _, m := range modules {
		for _, r := range m.GetRoutes() {
			unAuth.Handle(r.Method, r.Path, r.Handler)
		}
	}

	// Create the HTTP server instance
	httpServer := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	return &GinServer{
		Modules:    modules,
		HttpServer: httpServer,
	}
}
