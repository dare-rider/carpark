package main

import (
	"flag"
	"fmt"
	"github.com/dare-rider/carpark/config"
	"github.com/dare-rider/carpark/constant"
	"github.com/dare-rider/carpark/db"
	"github.com/dare-rider/carpark/router"
	"net/http"
)

const (
	defaultConfigPath = "config/config.yml"
)

func main() {
	// Set all the flags with their defaults
	var configFilePath string
	flag.StringVar(&configFilePath, "config", defaultConfigPath, "absolute path to the configuration file")
	flag.Parse()

	// Load all dependencies and routes
	generalConfig := config.LoadConfig(configFilePath)
	constant.InitConstants(generalConfig.MiscConfig)
	db.InitMysqlDb(generalConfig.DbConfig)
	db.InitMigrations(generalConfig.DbConfig)
	routes := router.InitRoutes(generalConfig, db.MysqlConn())

	// Started server.
	http.ListenAndServe(":3005", routes)
	fmt.Println("Terminated")
}
