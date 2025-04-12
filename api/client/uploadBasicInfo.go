package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"komari/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadBasicInfo(c *gin.Context) {
	db := database.GetSQLiteInstance()
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid request body"})
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid JSON"})
		return
	}
	if data["token"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Token is required"})
		return
	}
	token := data["token"].(string)

	clientUUID, err := database.GetClientUUIDByToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": fmt.Sprintf("%v", err)})
		return
	}
	// Validate the required fields in the data
	requiredFields := []string{"cpu", "os"}
	for _, field := range requiredFields {
		if data[field] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": fmt.Sprintf("%s is required", field)})
			return
		}
	}

	// Validate nested fields in "cpu"
	cpuFields := []string{"name", "arch", "cores"}
	cpuData, ok := data["cpu"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid CPU data"})
		return
	}
	for _, field := range cpuFields {
		if cpuData[field] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": fmt.Sprintf("CPU %s is required", field)})
			return
		}
	}
	db.Exec(`INSERT OR REPLACE into ClientsInfo (ClientUUID,CPUNAME,CPUARCH,CPUCORES,OS,GPUNAME) values (?,?,?,?,?,?)`,
		clientUUID,
		data["cpu"].(map[string]interface{})["name"],
		data["cpu"].(map[string]interface{})["arch"],
		data["cpu"].(map[string]interface{})["cores"],
		data["os"],
		nil)

	// log.Println(string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body for further use
	c.JSON(200, gin.H{"status": "success"})
}
