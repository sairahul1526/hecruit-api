package main

import (
	"math/rand"
	"time"

	API "hecruit-backend/api"
	CONFIG "hecruit-backend/config"
	DATABASE "hecruit-backend/database"
)

func main() {

	rand.Seed(time.Now().UnixNano()) // seed for random generator

	CONFIG.LoadConfig()
	DATABASE.ConnectDatabase()

	API.StartServer()

}
