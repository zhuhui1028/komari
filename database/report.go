package database

import (
	"fmt"
	"strings"
)

type Report struct {
	UUID        string             `json:"uuid"`
	Token       string             `json:"token"`
	CPU         CPU_Report         `json:"cpu"`
	OS          string             `json:"os"`
	Ram         Ram_Report         `json:"ram"`
	Swap        Ram_Report         `json:"swap"`
	Load        Load_Report        `json:"load"`
	Disk        Disk_Report        `json:"disk"`
	Network     Network_Report     `json:"network"`
	Connections Connections_Report `json:"connections"`
	Uptime      int64              `json:"uptime"`
	Process     int                `json:"process"`
	IPAddress   IPAddress_Report   `json:"ip"`
	Message     string             `json:"message"`
}
type IPAddress_Report struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}
type CPU_Report struct {
	Name  string  `json:"name"`
	Cores int     `json:"cores"`
	Arch  string  `json:"arch"`
	Usage float64 `json:"usage"`
}
type GPU_Report struct {
	Name  string  `json:"name"`
	Usage float64 `json:"usage"`
}
type Ram_Report struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}
type Load_Report struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}
type Disk_Report struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}
type Network_Report struct {
	Up        int64 `json:"up"`
	Down      int64 `json:"down"`
	TotalUp   int64 `json:"totalUp"`
	TotalDown int64 `json:"totalDown"`
}
type Connections_Report struct {
	TCP int `json:"tcp"`
	UDP int `json:"udp"`
}

func SaveReport(data map[string]interface{}) (err error) {
	report := ParseReport(data)
	token := report.Token
	clientUUID, err := GetClientUUIDByToken(token)
	if err != nil {
		return err
	}
	return SaveClientReport(clientUUID, report)
}

func SaveClientReport(clientUUID string, report Report) (err error) {
	db := GetSQLiteInstance()

	// Prepare the SQL insert statement
	query := `
        INSERT INTO History (
            Client, Time, CPU, GPU, RAM, RAM_TOTAL, SWAP, SWAP_TOTAL, LOAD, TEMP,
            DISK, DISK_TOTAL, NET_IN, NET_OUT, NET_TOTAL_UP, NET_TOTAL_DOWN, PROCESS, Connections, Connections_UDP
        ) VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	// Execute the insert statement using fields from the Report struct
	_, err = db.Exec(query,
		clientUUID,
		report.CPU.Usage,         // CPU
		nil,                      // GPU (not provided in Report)
		report.Ram.Used,          // RAM
		report.Ram.Total,         // RAM_TOTAL
		report.Swap.Used,         // SWAP
		report.Swap.Total,        // SWAP_TOTAL
		report.Load.Load1,        // LOAD (using Load1 as primary load metric)
		nil,                      // TEMP (not provided in Report)
		report.Disk.Used,         // DISK
		report.Disk.Total,        // DISK_TOTAL
		report.Network.Down,      // NET_IN
		report.Network.Up,        // NET_OUT
		report.Network.TotalUp,   // NET_TOTAL_UP
		report.Network.TotalDown, // NET_TOTAL_DOWN
		report.Process,           // PROCESS
		report.Connections.TCP,   // Connections
		report.Connections.UDP,   // Connections_UDP
	)

	if err != nil {
		return fmt.Errorf("failed to save report: %v", err)
	}

	return nil
}

func ParseReport(data map[string]interface{}) (report Report) {
	report = Report{
		Token: getString(data, "token"),
		CPU: CPU_Report{
			Name:  getString(data, "cpu.name"),
			Cores: getInt(data, "cpu.cores"),
			Arch:  getString(data, "cpu.arch"),
			Usage: getFloat(data, "cpu.usage"),
		},
		OS: getString(data, "os"),
		Ram: Ram_Report{
			Total: getInt64(data, "ram.total"),
			Used:  getInt64(data, "ram.used"),
		},
		Swap: Ram_Report{
			Total: getInt64(data, "swap.total"),
			Used:  getInt64(data, "swap.used"),
		},
		Load: Load_Report{
			Load1:  getFloat(data, "load.load1"),
			Load5:  getFloat(data, "load.load5"),
			Load15: getFloat(data, "load.load15"),
		},
		Disk: Disk_Report{
			Total: getInt64(data, "disk.total"),
			Used:  getInt64(data, "disk.used"),
		},
		Network: Network_Report{
			Up:        getInt64(data, "network.up"),
			Down:      getInt64(data, "network.down"),
			TotalUp:   getInt64(data, "network.totalUp"),
			TotalDown: getInt64(data, "network.totalDown"),
		},
		Connections: Connections_Report{
			TCP: getInt(data, "connections.tcp"),
			UDP: getInt(data, "connections.udp"),
		},
		Uptime:  getInt64(data, "uptime"),
		Process: getInt(data, "process"),
		Message: getString(data, "message"),
	}
	return
}

// getString 获取嵌套字符串值
func getString(data map[string]interface{}, key string) string {
	val := getNestedValue(data, key)
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// getInt 获取嵌套整数值
func getInt(data map[string]interface{}, key string) int {
	val := getNestedValue(data, key)
	if num, ok := val.(float64); ok {
		return int(num)
	}
	if num, ok := val.(int); ok {
		return num
	}
	return 0
}

// getInt64 获取嵌套 int64 值
func getInt64(data map[string]interface{}, key string) int64 {
	val := getNestedValue(data, key)
	if num, ok := val.(float64); ok {
		return int64(num)
	}
	if num, ok := val.(int64); ok {
		return num
	}
	return 0
}

// getFloat 获取嵌套浮点数值
func getFloat(data map[string]interface{}, key string) float64 {
	val := getNestedValue(data, key)
	if num, ok := val.(float64); ok {
		return num
	}
	return 0.0
}

// getNestedValue 通用函数，用于解析嵌套键
func getNestedValue(data map[string]interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	current := data

	// 遍历所有键，逐步深入嵌套
	for i, k := range keys[:len(keys)-1] {
		val, ok := current[k]
		if !ok {
			return nil
		}
		// 确保当前值是 map 类型
		nextMap, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}
		current = nextMap

		// 如果是最后一个键，返回值
		if i == len(keys)-2 {
			return current[keys[len(keys)-1]]
		}
	}
	return nil
}
