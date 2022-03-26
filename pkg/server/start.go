package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"replace-url-gin/config"
	"replace-url-gin/router"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Start() error
}

type server struct {
	ctx   context.Context
	conf  *config.Config
	route *gin.Engine
}

func NewApp(ctx context.Context) Server {
	return &server{
		ctx:   ctx,
		conf:  config.Get(),
		route: gin.New(),
	}
}

func (a *server) Start() error {
	a.provider()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.conf.Port),
		Handler: a.route,
	}

	go func() {
		log.Printf("Server starting on port :%d\n", a.conf.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return a.handleShutdown(server)
}

func (a *server) provider() {
	a.route.Use(gin.Logger())
	a.route.Use(gin.Recovery())
	router.Register(a.ctx, a.conf, a.route).All()
}

func (a server) handleShutdown(srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return err
}
