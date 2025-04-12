package main

import (
	"komari/cmd"
	"komari/database"
)

func main() {

	database.InitDatabase()

	cmd.Execute()

	//TODO: Auto Clean 7 days

}
