package database

import (
	"errors"
	"time"
)

type Sessions struct {
	UUID    string
	Session string
	Expires time.Time
}

func GetAllSessions() (sessions []Sessions, err error) {
	db := GetSQLiteInstance()
	rows, err := db.Query(`
		SELECT UUID, SESSION, EXPIRES FROM Sessions;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Sessions
		err = rows.Scan(&session.UUID, &session.Session, &session.Expires)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
func CreateSession(uuid string, expires int) (string, error) {
	session := generateRandomString(32)
	db := GetSQLiteInstance()
	_, err := db.Exec(`
		INSERT INTO Sessions (UUID, SESSION, EXPIRES) VALUES (?,?,?);
	`, uuid, session, time.Now().Add(time.Duration(expires)*time.Second))
	if err != nil {
		return "", err
	}
	return session, nil
}

func GetSession(session string) (uuid string, err error) {
	db := GetSQLiteInstance()
	row := db.QueryRow(`
		SELECT UUID, EXPIRES FROM Sessions WHERE SESSION = ?;
	`, session)

	var expires time.Time
	err = row.Scan(&uuid, &expires)
	if err != nil {
		return "", err
	}

	if time.Now().After(expires) {
		// Session expired, delete it
		_ = DeleteSession(session)
		return "", errors.New("session expired")
	}

	return uuid, nil
}

func DeleteSession(session string) (err error) {
	db := GetSQLiteInstance()
	_, err = db.Exec(`
		DELETE FROM Sessions WHERE SESSION = ?;
	`, session)
	if err != nil {
		return err
	}
	return nil
}
