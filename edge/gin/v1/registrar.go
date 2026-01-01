package v1

import (
	"github.com/Ahu-Tools/example/edge/gin/v1/hello"
	"github.com/gin-gonic/gin"
	// @ahum: imports
)

func RegisterVersion(r *gin.RouterGroup) {
	hello.RegisterRoutes(r.Group("hello"))
	// @ahum: entities
}
