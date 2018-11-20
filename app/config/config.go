package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/viper"
)

const (
	LOCAL      = "local"
	DEVELOP    = "develop"
	STAGE      = "stage"
	PRODUCTION = "production"
)

var Config *viper.Viper
var currentEnv string

func init() {
	Config = viper.New()
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Config.AutomaticEnv()

	if environment, ok := os.LookupEnv("GO_ENVIRONMENT"); !ok {
		currentEnv = DEVELOP
	} else {
		currentEnv = environment
	}
	Config.SetConfigName(currentEnv)
	gopath, _ := os.LookupEnv("GOPATH")
	Config.SetConfigFile(gopath + "/src/github.com/wenance/wequeue-management_api/app/config/" + currentEnv + ".yml")
	if currentEnv == LOCAL {
		Config.Set("aws_credentials", credentials.NewStaticCredentials("foo", "bar", ""))
	}	
	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %#v", err))
	}
	fmt.Printf("Current environment: %s \n", currentEnv)
}

func Get(key string) string {
	return Config.GetString(key)
}

func GetObject(key string) interface{} {
	return Config.Get(key)
}

func GetCurrentEnvironment() *string {
	return &currentEnv
}

func SetCurrentEnvironment(environment string) {
	currentEnv = environment
}
