package v1

import (
	"context"

	"connectrpc.com/connect"
	hellov1 "github.com/Ahu-Tools/example/edge/connect/gen/hello/v1"
	// @ahum: imports
)

type Edge struct {
}

func NewEdge() *Edge {
	return &Edge{}
}

func (e *Edge) Health(c context.Context, req *connect.Request[hellov1.HealthRequest]) (*connect.Response[hellov1.HealthResponse], error) {
	res := connect.NewResponse(&hellov1.HealthResponse{
		Message: "UP",
	})

	return res, nil
}

func (e *Edge) World(c context.Context, req *connect.Request[hellov1.WorldRequest]) (*connect.Response[hellov1.WorldResponse], error) {
	res := connect.NewResponse(&hellov1.WorldResponse{
		Message: "Congratulation, World is up!",
	})

	return res, nil
}

// @ahum: methods
