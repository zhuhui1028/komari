package client

import (
	"database/sql"
	"fmt"
	"komari/database/clients"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RemoteConfig struct {
	Cpu         bool `json:"cpu"`
	Gpu         bool `json:"gpu"`
	Ram         bool `json:"ram"`
	Swap        bool `json:"swap"`
	Load        bool `json:"load"`
	Uptime      bool `json:"uptime"`
	Temperature bool `json:"temperature"`
	Os          bool `json:"os"`
	Disk        bool `json:"disk"`
	Network     bool `json:"network"`
	Process     bool `json:"process"`
	Interval    int  `json:"interval"`
	Connections bool `json:"connections"`
}

func GetRemoteConfig(c *gin.Context) {
	token := c.Query("token")

	clientUUID, err := clients.GetClientUUIDByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}
	config, err := clients.GetClientConfig(clientUUID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "No data found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": fmt.Errorf("error querying client: %v", err)})
		return
	}

	c.JSON(http.StatusOK, config)

}
