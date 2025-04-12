package admin

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	log.Println("WARNING: GetSettings is not implemented")
	c.JSON(200, gin.H{})
}
func EditSettings(c *gin.Context) {
	log.Println("WARNING: EditSettings is not implemented")
	c.JSON(200, gin.H{})
}
