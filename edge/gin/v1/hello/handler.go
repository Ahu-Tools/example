package hello

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// @ahum: imports
)

type Handler struct {
	//Add your chain here
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

func (h Handler) World(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}

// @ahum: handlers
