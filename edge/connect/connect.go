package connect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/Ahu-Tools/example/edge/connect/hello"
	// @ahum: imports
)

type Server struct {
	srv *http.Server
}

func New() *Server {
	return &Server{}
}

func RegisterServices(mux *http.ServeMux) {
	hello.RegisterService(mux)
	// @ahum: services
}

func (s *Server) Configure() {
	mux := http.NewServeMux()

	RegisterServices(mux)

	host := viper.GetString("edges.connect.server.host")
	port := viper.GetInt("edges.connect.server.port")
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Printf("Connect server listening on %s", s.srv.Addr)

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting Connect server on %s: %v", s.srv.Addr, err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Shutting down Connect server on %s...", s.srv.Addr)
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down Connect server: %v", err)
	}
}
