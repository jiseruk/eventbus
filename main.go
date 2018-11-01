package main

import (
	"github.com/wenance/wequeue-management_api/app/server"
	_ "github.com/wenance/wequeue-management_api/app/validation"
	_ "github.com/wenance/wequeue-management_api/docs"
)

// @title WeQueue Management API
// @version 0.1
// @description This is the Wenance event-bus Management Api
// @host localhost:8080
// @BasePath /
func main() {
	server.Init()
}
