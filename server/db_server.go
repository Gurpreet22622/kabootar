package server

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

var dbhandler *sql.DB

func InitDatabase(config *viper.Viper) *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up Viper to read from environment variables
	viper.AutomaticEnv()

	// Set default values or get values from environment variables
	connectionString := viper.GetString("DATABASE_URL")
	//maxIdleConnections := viper.GetInt("DB_MAX_IDLE_CONNECTIONS")
	//maxOpenConnections := viper.GetInt("DB_MAX_OPEN_CONNECTIONS")
	//connectionMaxLifetime := viper.GetDuration("DB_CONNECTION_MAX_LIFETIME")
	driverName := "postgres"

	if connectionString == "" {
		log.Fatal("Database connection string is missing")
	}
	dbHandler, err := sql.Open(driverName, connectionString)
	if err != nil {
		log.Fatal("Error while initializing database: ", err)
	}
	//dbHandler.SetMaxIdleConns(maxIdleConnections)
	//dbHandler.SetMaxOpenConns(maxOpenConnections)
	//dbHandler.SetConnMaxLifetime(connectionMaxLifetime)

	// connectionString := config.GetString("database.connection_string")
	// maxIdleConnections := config.GetInt("database.max_idle_connectons")
	// maxOpenConnections := config.GetInt("database.max_open_connections")
	// connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	// driverName := config.GetString("database.driver_name")
	// if connectionString == "" {
	// 	log.Fatal("Database connection string is missing")
	// }
	// dbHandler, err := sql.Open(driverName, connectionString)
	// if err != nil {
	// 	log.Fatal("Error while initializing database: ", err)
	// }
	//dbHandler.SetMaxIdleConns(maxIdleConnections)
	//dbHandler.SetMaxOpenConns(maxOpenConnections)
	//dbHandler.SetConnMaxLifetime(connectionMaxLifetime)
	err = dbHandler.Ping()
	if err != nil {
		dbHandler.Close()
		log.Fatal("Error while validating database: ", err)
	}
	dbhandler = dbHandler
	return dbHandler
}
