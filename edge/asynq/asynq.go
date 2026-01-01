package asynq

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	// @ahum: imports
)

type Handler func(context.Context, *asynq.Task) error

var handlers = make(map[string]Handler)

type Server struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func New() *Server {
	return &Server{}
}

func RegisterTasks(mux *asynq.ServeMux) {
	for pattern, handler := range handlers {
		mux.HandleFunc(pattern, handler)
	}
}

func (s *Server) Configure() {
	s.mux = asynq.NewServeMux()

	RegisterTasks(s.mux)

	concurrency := viper.GetInt("edges.asynq.concurrency")
	rawQueues := viper.GetStringMap("edges.asynq.queues")
	var queues = make(map[string]int)
	for k, v := range rawQueues {
		queues[k] = v.(int)
	}

	s.srv = asynq.NewServerFromRedisClient(
		getRedis(),
		asynq.Config{
			Concurrency: concurrency,
			Queues:      queues,
		},
	)
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Printf("Asynq task server is going to run...")
		if err := s.srv.Run(s.mux); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting asynq task server: %e", err)
		}
	}()

	<-ctx.Done()

	log.Printf("Shutting down asynq server...")
	s.srv.Shutdown()
}

func getRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("edges.asynq.redis.host"), viper.GetInt("edges.asynq.redis.port")),
		Username: viper.GetString("edges.asynq.redis.username"),
		Password: viper.GetString("edges.asynq.redis.password"),
		DB:       viper.GetInt("edges.asynq.redis.db"),
	})

	return rdb
}

func RegisterHandler(version, pattern string, handler Handler) {
	handlers[version+":"+pattern] = handler
}
