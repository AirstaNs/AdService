package httpgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	grpc2 "homework10/internal/ports/grpc"
	"log"
	"net/http"
)

type HttpServer struct {
	App *http.Server
	log *log.Logger
}

func NewHTTPServer(port string, a app.App, logger *log.Logger) grpc2.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	api := handler.Group("/api/v1/")
	AppRouter(api, a, logger)
	srv := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return &HttpServer{App: srv, log: logger}
}

func (s *HttpServer) Start(network, address string) error {
	s.log.Printf("starting http server started at %s:%s", network, address)
	if err := s.App.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Fatalf("listen: %s\n", err)
		return err
	}
	return nil
}

func (s *HttpServer) Stop() error {
	return s.App.Shutdown(context.Background())
}
