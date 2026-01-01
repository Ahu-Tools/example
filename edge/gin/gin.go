package gin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/Ahu-Tools/example/edge/gin/v1"
	// @ahum: imports
)

type Server struct {
	srv *http.Server
}

func New() *Server {
	return &Server{}
}

func RegisterRoutes(r *gin.Engine) {
	v1.RegisterVersion(r.Group("v1"))
	// @ahum: versions
}

func (s *Server) Configure() {
	router := gin.Default()

	// Register Gin routes
	RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	host := viper.GetString("edges.gin.server.host")
	port := viper.GetInt("edges.gin.server.port")
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Printf("Gin Server listening on %s", s.srv.Addr)

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting Gin server on %s: %v", s.srv.Addr, err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Shutting down Gin server on %s...", s.srv.Addr)
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down Gin server: %v", err)
	}
}
