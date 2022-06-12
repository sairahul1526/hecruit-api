package main

import (
	"math/rand"
	"time"

	API "hecruit-backend/api"
	CRON "hecruit-backend/api/cron"
	CONFIG "hecruit-backend/config"
	DATABASE "hecruit-backend/database"
)

func main() {

	rand.Seed(time.Now().UnixNano()) // seed for random generator

	CONFIG.LoadConfig()
	DATABASE.ConnectDatabase()

	go CRON.SendEmailsContinously()

	API.StartServer()

}
