package admin

import (
	"database/sql"
	"komari/database"
	"log"

	"github.com/gin-gonic/gin"
)

// Custom 结构体定义，首字母大写以便导出
type Custom struct {
	CustomCss string `json:"customCss"`
	CustomJs  string `json:"customJs"`
}

// GetCustom 获取自定义配置
func GetCustom(c *gin.Context) {
	db := database.GetSQLiteInstance()
	var custom Custom

	// 查询单行数据
	err := db.QueryRow("SELECT CustomCss, CustomJs FROM Custom LIMIT 1").Scan(&custom.CustomCss, &custom.CustomJs)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有数据，返回默认值或空对象
			c.JSON(200, Custom{})
			return
		}
		log.Printf("Failed to query custom data: %v", err)
		c.JSON(500, gin.H{"status": "error", "error": "Internal Server Error"})
		return
	}

	c.JSON(200, custom)
}

// EditCustom 更新自定义配置
func EditCustom(c *gin.Context) {
	var req Custom
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(400, gin.H{
			"status": "error",
			"error":  "Invalid request body",
		})
		return
	}

	db := database.GetSQLiteInstance()

	// 检查是否存在记录
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Custom").Scan(&count)
	if err != nil {
		log.Printf("Failed to check custom table: %v", err)
		c.JSON(500, gin.H{"status": "error", "error": "Internal Server Error"})
		return
	}

	if count == 0 {
		// 如果表为空，插入新记录
		_, err = db.Exec("INSERT INTO Custom (CustomCss, CustomJs) VALUES (?, ?)", req.CustomCss, req.CustomJs)
	} else {
		// 否则更新现有记录（假设只有一行，使用 LIMIT 1 确保安全）
		_, err = db.Exec("UPDATE Custom SET CustomCss = ?, CustomJs = ? LIMIT 1", req.CustomCss, req.CustomJs)
	}

	if err != nil {
		log.Printf("Failed to update/insert custom data: %v", err)
		c.JSON(500, gin.H{"status": "error", "error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
