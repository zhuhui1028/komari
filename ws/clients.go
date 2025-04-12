package ws

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetClients(c *gin.Context) {
	log.Println("WARNING: GetSettings is not implemented")
	c.JSON(200, gin.H{})
}
