package clients

import (
	"komari/database/dbcore"
	"komari/database/models"
	"komari/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 删除指定 UUID 的客户端配置
func DeleteClientConfig(clientUuid string) error {
	db := dbcore.GetDBInstance()
	err := db.Delete(&models.ClientConfig{ClientUUID: clientUuid}).Error
	if err != nil {
		return err
	}
	return nil
}

// 更新或插入客户端基本信息
func UpdateOrInsertBasicInfo(cbi ClientBasicInfo) error {
	db := dbcore.GetDBInstance()
	err := db.Save(&cbi).Error
	if err != nil {
		return err
	}
	return nil
}

// 更新客户端设置
func UpdateClientConfig(config models.ClientConfig) error {
	db := dbcore.GetDBInstance()
	err := db.Save(&config).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateClientByUUID 更新指定 UUID 的客户端配置
func UpdateClientByUUID(config models.Client) error {
	db := dbcore.GetDBInstance()
	result := db.Model(&models.Client{}).Where("uuid = ?", config.UUID).Updates(config)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetClientUUIDByToken 根据 Token 获取客户端 UUID
func GetClientUUIDByToken(token string) (uuid string, err error) {
	db := dbcore.GetDBInstance()
	var client models.Client
	err = db.Where("token = ?", token).First(&client).Error
	if err != nil {
		return "", err
	}
	return client.UUID, nil
}

// CreateClient 创建新客户端
func CreateClient(config models.ClientConfig) (clientUUID, token string, err error) {
	db := dbcore.GetDBInstance()
	token = utils.GenerateToken()
	clientUUID = uuid.New().String()

	client := models.Client{
		UUID:      clientUUID,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.Create(&client).Error
	if err != nil {
		return "", "", err
	}
	err = db.Create(&config).Error
	if err != nil {
		return "", "", err
	}
	return clientUUID, token, nil
}

// GetAllClients 获取所有客户端配置
func GetAllClients() (clients []models.Client, err error) {
	db := dbcore.GetDBInstance()
	err = db.Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// GetClientConfig 获取指定 UUID 的客户端配置
func GetClientConfig(uuid string) (client models.Client, err error) {
	db := dbcore.GetDBInstance()
	err = db.Where("uuid = ?", uuid).First(&client).Error
	if err != nil {
		return client, err
	}
	return client, nil
}

// ClientBasicInfo 客户端基本信息（假设的结构体，需根据实际定义调整）
type ClientBasicInfo struct {
	CPU       CPUReport       `json:"cpu"`
	GPU       GPUReport       `json:"gpu"`
	IpAddress IPAddressReport `json:"ip"`
	OS        string          `json:"os"`
}

// GetClientBasicInfo 获取指定 UUID 的客户端基本信息
func GetClientBasicInfo(uuid string) (client ClientBasicInfo, err error) {
	db := dbcore.GetDBInstance()
	var clientInfo models.ClientInfo
	err = db.Where("client_uuid = ?", uuid).First(&clientInfo).Error
	if err != nil {
		return client, err
	}

	client = ClientBasicInfo{
		CPU: CPUReport{
			Name:  clientInfo.CPUNAME,
			Arch:  clientInfo.CPUARCH,
			Cores: clientInfo.CPUCORES,
		},
		GPU: GPUReport{
			Name: clientInfo.GPUNAME,
		},
		OS: clientInfo.OS,
		// IpAddress: 未在数据库中找到对应字段，需确认
	}

	return client, nil
}
