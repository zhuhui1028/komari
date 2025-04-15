package api

import (
	"komari/database/accounts"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	userName := "Guest"
	session, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(200, gin.H{"username": userName})
		return
	}
	userName, err = accounts.GetSession(session)
	if err != nil {
		userName = "Guest"
	}
	c.JSON(200, gin.H{"username": userName})

}
