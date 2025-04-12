package admin

import (
	"komari/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddClient(c *gin.Context) {
	var config database.ClientConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	uuid, token, err := database.CreateClient(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "uuid": uuid, "token": token})
}

func EditClient(c *gin.Context) {
	var config database.ClientConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if config.UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "UUID is required"})
		return
	}
	if config.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Token is required"})
		return
	}
	if config.Interval <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Interval must be greater than to 0"})
		return
	}
	err := database.UpdateClientByUUID(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func RemoveClient(c *gin.Context) {
	db := database.GetSQLiteInstance()
	var req struct {
		UUID string `json:"uuid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "Invalid or missing UUID",
		})
	}
	_, err := db.Exec("DELETE FROM Clients WHERE UUID = ?", req.UUID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  "Failed to remove client" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
func ClearHistory(c *gin.Context) {
	db := database.GetSQLiteInstance()
	var req struct {
		UUID string `json:"uuid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "Invalid or missing UUID",
		})
	}
	_, err := db.Exec("DELETE FROM History WHERE Client = ?", req.UUID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  "Failed to clear history" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func GetClient(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "Invalid or missing UUID",
		})
		return
	}
	result := map[string]interface{}{}

	clientBasicInfo, err := database.GetClientBasicInfo(uuid)
	if err.Error() != "sql: no rows in result set" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	config, err := database.GetClientConfig(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	result = map[string]interface{}{
		"uuid":   uuid,
		"info":   clientBasicInfo,
		"config": config,
	}
	c.JSON(200, result)
}
func ListClients(c *gin.Context) {
	clients, err := database.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	result := []map[string]interface{}{}
	for i := range clients {
		clientBasicInfo, err := database.GetClientBasicInfo(clients[i].UUID)
		if err.Error() != "sql: no rows in result set" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
		result = append(result, map[string]interface{}{
			"uuid":   clients[i].UUID,
			"info":   clientBasicInfo,
			"config": clients[i],
		})
	}
	c.JSON(http.StatusOK, result)
}
