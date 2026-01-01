package edge

import (
	"context"
	"github.com/Ahu-Tools/example/edge/asynq"
	"github.com/Ahu-Tools/example/edge/connect"
	"github.com/Ahu-Tools/example/edge/gin"
	"sync"
	// @ahum: imports
)

type Edge interface {
	Configure()
	Start(context.Context, *sync.WaitGroup)
}

func Start(ctx context.Context, wg *sync.WaitGroup) {
	edges := []Edge{
		gin.New(),
		connect.New(),
		asynq.New(),
		// @ahum: edges
	}

	wg.Add(len(edges))
	for _, edge := range edges {
		edge.Configure()
		go edge.Start(ctx, wg)
	}

}
