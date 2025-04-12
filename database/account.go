package database

import (
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/google/uuid"
)

type User struct {
	UUID     string
	Username string
	Passwd   string
}

/*
CREATE TABLE IF NOT EXISTS Users (

	UUID TEXT PRIMARY KEY,
	Username TEXT UNIQUE,
	Passwd TEXT

);
*/
// CheckPassword checks if the password is correct
//
// Returns the UUID of the user if the password is correct, otherwise returns an empty string and false
func CheckPassword(username, passwd string) (uuid string, success bool) {
	db := GetSQLiteInstance()
	var hashedPassword string
	err := db.QueryRow("SELECT UUID, Passwd FROM Users WHERE Username = ?", username).Scan(&uuid, &hashedPassword)
	if err != nil {
		return "", false
	}
	if hashPasswd(passwd) != hashedPassword {
		return "", false
	}
	return uuid, true
}
func ForceResetPassword(username, passwd string) (err error) {
	db := GetSQLiteInstance()

	_, err = db.Exec("UPDATE accounts SET passwd = ? WHERE username = ?", hashPasswd(passwd), username)
	if err != nil {
		return err
	}
	return nil
}

func hashPasswd(passwd string) string {
	saltedPassword := passwd + constantSalt
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))
	hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hashedPassword
}

func CreateDefaultAdminAccount() (username, passwd string, err error) {
	db := GetSQLiteInstance()
	username = os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}

	passwd = os.Getenv("ADMIN_PASSWORD")
	if passwd == "" {
		passwd = generatePassword()
	}

	hashedPassword := hashPasswd(passwd)

	userUUID := uuid.New().String()
	_, err = db.Exec(`
        INSERT INTO Users (UUID, Username, Passwd) VALUES (?, ?, ?);
    `, userUUID, username, hashedPassword)
	if err != nil {
		return "", "", err
	}

	return username, passwd, nil
}
