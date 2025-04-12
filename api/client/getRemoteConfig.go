package client

import (
	"database/sql"
	"fmt"
	"komari/database"
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
	db := database.GetSQLiteInstance()

	clientUUID, err := database.GetClientUUIDByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}

	var config RemoteConfig
	query := `
		SELECT CPU, GPU, RAM, SWAP, LOAD, UPTIME, TEMP, OS, DISK, NET, PROCESS, Connections, Interval
		FROM Clients
		WHERE UUID = ?
	`
	row := db.QueryRow(query, clientUUID)

	err = row.Scan(
		&config.Cpu,
		&config.Gpu,
		&config.Ram,
		&config.Swap,
		&config.Load,
		&config.Uptime,
		&config.Temperature,
		&config.Os,
		&config.Disk,
		&config.Network,
		&config.Process,
		&config.Connections,
		&config.Interval,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "No data found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": fmt.Errorf("error querying client: %v", err)})
		return
	}

	c.JSON(http.StatusOK, config)

}
