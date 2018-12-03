package main

import (
	"fmt"

	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/server"
	_ "github.com/wenance/wequeue-management_api/app/validation"
	_ "github.com/wenance/wequeue-management_api/docs"
)

// @title WeQueue Management API
// @version 0.1
// @description This is the Wenance event-bus Management Api
// @host hhttp://localhost:8080
// @BasePath /
// @contact.name Javier Iseruk
// @contact.url http://www.swagger.io/support
// @contact.email javier.iseruk@wenance.com
func main() {
	server.Init()
	fmt.Printf("Service ready, environment %s", *config.GetCurrentEnvironment())
}
