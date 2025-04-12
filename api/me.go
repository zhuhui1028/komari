package api

import (
	"komari/database"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	userName := "Guest"
	db := database.GetSQLiteInstance()
	session, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(200, gin.H{"username": userName})
		return
	}
	row := db.QueryRow(`
			SELECT Username FROM Users WHERE UUID = (SELECT UUID FROM Sessions WHERE SESSION = ?);
		`, session)
	err = row.Scan(&userName)
	if err != nil {
		c.JSON(200, gin.H{"username": userName})
		return
	}
	c.JSON(200, gin.H{"username": userName})

}
