package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

func SetNewrelic(router *gin.Engine) {
	if os.Getenv("NEW_RELIC_LICENSE_KEY") == "" || os.Getenv("NEW_RELIC_APP_NAME") == "" {
		return
	}
	cfg := newrelic.NewConfig(os.Getenv("NEW_RELIC_APP_NAME"), os.Getenv("NEW_RELIC_LICENSE_KEY"))
	app, err := newrelic.NewApplication(cfg)
	if err != nil {
		fmt.Printf("Error setting newrelic agent : %s", err.Error())
		panic("Error seting Newrelic agent")
	}
	router.Use(nrgin.Middleware(app))
}
