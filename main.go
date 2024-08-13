package main

import (
	"database/sql"
	"kabootar/config"
	"kabootar/server"
	"log"

	_ "github.com/lib/pq"
)

var dbHandler *sql.DB

func main() {
	log.Println("Starting Kabootar")
	log.Println("Initializing configuration")
	config := config.InitConfig("kabootar")
	log.Println("Initializing Database")
	dbHandler = server.InitDatabase(config)
	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}
