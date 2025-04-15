package main

import (
	"komari/cmd"
	"komari/database/accounts"
	"komari/database/dbcore"
	"komari/database/history"
	"log"
	"time"
)

func main() {

	if !dbcore.InitDatabase() {
		user, passwd, err := accounts.CreateDefaultAdminAccount()
		if err != nil {
			panic(err)
		}
		log.Println("Default admin account created. Username:", user, ", Password:", passwd)
	}

	// Delete old history
	go func() {
		ticker := time.NewTicker(time.Hour * 1)

		select {
		case <-ticker.C:
			history.DeleteRecordBefore(time.Now().Add(-time.Hour * 24 * 7))
		}
	}()

	cmd.Execute()
}
