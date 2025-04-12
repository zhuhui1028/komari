package database

import "github.com/google/uuid"

type ClientConfig struct {
	UUID        string `json:"uuid"`
	Token       string `json:"token"`
	CPU         bool   `json:"cpu"`
	GPU         bool   `json:"gpu"`
	RAM         bool   `json:"ram"`
	SWAP        bool   `json:"swap"`
	LOAD        bool   `json:"load"`
	UPTIME      bool   `json:"uptime"`
	TEMP        bool   `json:"temp"`
	OS          bool   `json:"os"`
	DISK        bool   `json:"disk"`
	NET         bool   `json:"net"`
	PROCESS     bool   `json:"process"`
	CONNECTIONS bool   `json:"connections"`
	Interval    int    `json:"interval"`
}

func UpdateClientByUUID(config ClientConfig) error {
	db := GetSQLiteInstance()
	_, err := db.Exec(`
		UPDATE Clients SET CPU = ?, GPU = ?, RAM = ?, SWAP = ?, LOAD = ?, UPTIME = ?, TEMP = ?, OS = ?, DISK = ?, NET = ?, PROCESS = ?, Connections = ?, Interval = ?, Token = ? WHERE UUID = ?;
	`, config.CPU, config.GPU, config.RAM, config.SWAP, config.LOAD, config.UPTIME, config.TEMP, config.OS, config.DISK, config.NET, config.PROCESS, config.CONNECTIONS, config.Interval, config.Token, config.UUID)

	if err != nil {
		return err
	}
	return nil
}
func GetClientUUIDByToken(token string) (uuid string, err error) {
	db := GetSQLiteInstance()

	var clientUUID string
	err = db.QueryRow("SELECT UUID FROM clients WHERE token = ?", token).Scan(&clientUUID)
	if err != nil {
		return "", err
	}

	return clientUUID, nil
}

func CreateClient(config ClientConfig) (clientUUID, token string, err error) {
	db := GetSQLiteInstance()
	token = generateToken()
	clientUUID = uuid.New().String()
	_, err = db.Exec(`
		INSERT INTO Clients (UUID, TOKEN, CPU, GPU, RAM, SWAP, LOAD, UPTIME, TEMP, OS, DISK, NET, PROCESS, Connections, Interval) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`, clientUUID, token, config.CPU, config.GPU, config.RAM, config.SWAP, config.LOAD, config.UPTIME, config.TEMP, config.OS, config.DISK, config.NET, config.PROCESS, config.CONNECTIONS, config.Interval)

	if err != nil {
		return "", "", err
	}
	return clientUUID, token, nil
}

// This function returns all clients CONFIG from the database
func GetAllClients() (clients []ClientConfig, err error) {
	db := GetSQLiteInstance()
	rows, err := db.Query(`
		SELECT UUID, TOKEN, CPU, GPU, RAM, SWAP, LOAD, UPTIME, TEMP, OS, DISK, NET, PROCESS, Connections, Interval FROM Clients;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var client ClientConfig
		err := rows.Scan(&client.UUID, &client.Token, &client.CPU, &client.GPU, &client.RAM, &client.SWAP, &client.LOAD, &client.UPTIME, &client.TEMP, &client.OS, &client.DISK, &client.NET, &client.PROCESS, &client.CONNECTIONS, &client.Interval)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return clients, nil
}
func GetClientConfig(uuid string) (client ClientConfig, err error) {
	db := GetSQLiteInstance()
	rows, err := db.Query(`
		SELECT UUID, TOKEN, CPU, GPU, RAM, SWAP, LOAD, UPTIME, TEMP, OS, DISK, NET, PROCESS, Connections, Interval FROM Clients WHERE UUID = ?;
	`, uuid)
	if err != nil {
		return client, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&client.UUID, &client.Token, &client.CPU, &client.GPU, &client.RAM, &client.SWAP, &client.LOAD, &client.UPTIME, &client.TEMP, &client.OS, &client.DISK, &client.NET, &client.PROCESS, &client.CONNECTIONS, &client.Interval)
		if err != nil {
			return client, err
		}
	}
	return client, nil
}

type ClientBasicInfo struct {
	CPU       CPU_Report       `json:"cpu"`
	GPU       GPU_Report       `json:"gpu"`
	IpAddress IPAddress_Report `json:"ip"`
	OS        string           `json:"os"`
}

func GetClientBasicInfo(uuid string) (client ClientBasicInfo, err error) {
	db := GetSQLiteInstance()
	row := db.QueryRow(`
		SELECT CPUNAME, CPUARCH, CPUCORES, OS ,GPUNAME FROM ClientsInfo WHERE ClientUUID = ?;
	`, uuid)
	err = row.Scan(&client.CPU.Name, &client.CPU.Arch, &client.CPU.Cores, &client.OS, &client.GPU.Name)
	if err != nil {
		return client, err
	}
	return client, nil
}
