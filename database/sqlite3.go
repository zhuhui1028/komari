package database

import (
	"database/sql"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"komari/cmd/flags"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
)

var (
	instance *sql.DB
	once     sync.Once
)

const (
	constantSalt = "hjv9NGgXDqjI6m"
)

func InitDatabase() {
	// 检查数据库文件是否存在
	if _, err := os.Stat(flags.DatabaseFile); os.IsNotExist(err) {
		log.Printf("Database file %q does not exist, creating...", flags.DatabaseFile)

		// 确保数据库文件目录存在
		dbDir := filepath.Dir(flags.DatabaseFile)
		if dbDir != "" { // 防止空路径导致错误
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				log.Fatalf("Failed to create directory %q for database file: %v", dbDir, err)
			}
		}

		// 创建数据库文件
		file, err := os.Create(flags.DatabaseFile)
		if err != nil {
			log.Fatalf("Failed to create SQLite3 database file %q: %v", flags.DatabaseFile, err)
		}
		if err := file.Close(); err != nil {
			log.Fatalf("Failed to close database file %q: %v", flags.DatabaseFile, err)
		}

		// 初始化 SQLite 数据库
		CreateSqliteDb()

		// 创建默认管理员账户
		username, passwd, err := CreateDefaultAdminAccount()
		if err != nil {
			log.Fatalf("Failed to create default admin account: %v", err)
		}
		log.Printf("Default admin account created with username: %s and password: %s", username, passwd)
	} else if err != nil {
		log.Fatalf("Failed to check database file %q: %v", flags.DatabaseFile, err)
	}
}

// GetSQLiteInstance returns a singleton instance of the SQLite3 database connection.
func GetSQLiteInstance() *sql.DB {

	once.Do(func() {
		var err error
		instance, err = sql.Open("sqlite3", flags.DatabaseFile)
		if err != nil {
			log.Fatalf("Failed to connect to SQLite3 database: %v", err)
		}
	})

	return instance
}

func CreateSqliteDb() {

	db := GetSQLiteInstance()

	log.Println("Creating SQLite3 database and tables...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Clients (
			UUID TEXT PRIMARY KEY,
			TOKEN TEXT UNIQUE,
			CPU BOOLEAN DEFAULT TRUE,
			GPU BOOLEAN DEFAULT TRUE,
			RAM BOOLEAN DEFAULT TRUE,
			SWAP BOOLEAN DEFAULT TRUE,
			LOAD BOOLEAN DEFAULT TRUE,
			UPTIME BOOLEAN DEFAULT TRUE,
			TEMP BOOLEAN DEFAULT TRUE,
			OS BOOLEAN DEFAULT TRUE,
			DISK BOOLEAN DEFAULT TRUE,
			NET BOOLEAN DEFAULT TRUE,
			PROCESS BOOLEAN DEFAULT TRUE,
			Connections BOOLEAN DEFAULT TRUE,
			Interval INTEGER DEFAULT 3
		);

		CREATE TABLE IF NOT EXISTS Users (
			UUID TEXT PRIMARY KEY,
			Username TEXT UNIQUE,
			Passwd TEXT
		);

		CREATE TABLE IF NOT EXISTS Config (
			Sitename TEXT
		);

		CREATE TABLE IF NOT EXISTS Custom (
			CustomCSS TEXT,
			CustomJS TEXT
		);

		CREATE TABLE IF NOT EXISTS History (
			Client TEXT,
			Time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CPU FLOAT,
			GPU FLOAT,
			RAM BIGINT,
			RAM_TOTAL BIGINT,
			SWAP BIGINT,
			SWAP_TOTAL BIGINT,
			LOAD FLOAT,
			TEMP FLOAT,
			DISK BIGINT,
			DISK_TOTAL BIGINT,
			NET_IN BIGINT,
			NET_OUT BIGINT,
			NET_TOTAL_UP BIGINT,
			NET_TOTAL_DOWN BIGINT,
			PROCESS INTEGER,
			Connections INTEGER,
			Connections_UDP INTEGER,
			FOREIGN KEY(Client) REFERENCES Clients(UUID) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS ClientsInfo (
			ClientUUID TEXT PRIMARY KEY,
			CPUNAME TEXT,
			CPUARCH TEXT,
			CPUCORES INTEGER,
			OS TEXT,
			GPUNAME TEXT,
			FOREIGN KEY(ClientUUID) REFERENCES Clients(UUID) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	// Sqlite do not support INTERVAL, use DATETIME('now', '+1 hour')
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Sessions (
    		UUID TEXT,
    		Session TEXT PRIMARY KEY,
    		Expires TIMESTAMP NOT NULL DEFAULT (DATETIME('now', '+1 hour')),
    		FOREIGN KEY(UUID) REFERENCES Users(UUID) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create table Sessions: %v", err)
	}
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, err := r.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

func generatePassword() string {
	return generateRandomString(12)
}

func generateToken() string {
	return generateRandomString(16) // Generate a 32-character random token
}
