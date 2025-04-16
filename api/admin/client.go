package admin

import (
	"net/http"

	"github.com/akizon77/komari/database/clients"
	"github.com/akizon77/komari/database/history"
	"github.com/akizon77/komari/database/models"

	"github.com/gin-gonic/gin"
)

func AddClient(c *gin.Context) {
	var config models.ClientConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	uuid, token, err := clients.CreateClient(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "uuid": uuid, "token": token})
}

func EditClient(c *gin.Context) {
	var config models.ClientConfig
	var client models.Client

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	if client.UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "UUID is required"})
		return
	}
	if client.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Token is required"})
		return
	}
	if config.Interval <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Interval must be greater than to 0"})
		return
	}
	err := clients.UpdateClientConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func RemoveClient(c *gin.Context) {

	var req struct {
		UUID string `json:"uuid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "Invalid or missing UUID",
		})
	}
	err := clients.DeleteClientConfig(req.UUID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  "Failed to delete client" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
func ClearHistory(c *gin.Context) {
	if err := history.DeleteAll(); err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  "Failed to delete history" + err.Error(),
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

	clientBasicInfo, err := clients.GetClientBasicInfo(uuid)
	if err.Error() != "sql: no rows in result set" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	config, err := clients.GetClientConfig(uuid)
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
	cls, err := clients.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	result := []map[string]interface{}{}
	for i := range cls {
		clientBasicInfo, err := clients.GetClientBasicInfo(cls[i].UUID)
		if err.Error() != "sql: no rows in result set" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
		config, err := clients.GetClientConfig(cls[i].UUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
		result = append(result, map[string]interface{}{
			"uuid":   cls[i].UUID,
			"info":   clientBasicInfo,
			"config": config,
		})
	}
	c.JSON(http.StatusOK, result)
}
