package v1

import (
	"net/http"

	"github.com/Ahu-Tools/example/edge/connect/gen/hello/v1/hellov1connect"
	// @ahum: imports
)

func RegisterVersion(mux *http.ServeMux) {
	edge := NewEdge()
	path, handler := hellov1connect.NewServiceHandler(edge)
	mux.Handle(path, handler)
	// @ahum: edges
}
