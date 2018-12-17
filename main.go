package main

import (
	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/server"
	_ "github.com/wenance/wequeue-management_api/app/validation"
	"github.com/wenance/wequeue-management_api/docs"
)

// @title WeQueue Management API
// @version 0.1
// @description Bondi (event-bus Management Api)
// @BasePath /
// @contact.name Javier Iseruk
// @contact.url http://www.swagger.io/support
// @contact.email javier.iseruk@wenance.com
func main() {
	docs.SwaggerInfo.Host = config.Get("bondi.endpoint")
	server.Init()
}
