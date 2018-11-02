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
	Config.SetConfigType("yaml")
	Config.AddConfigPath("./app/config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.MergeInConfig()
	if os.Getenv("GO_ENVIRONMENT") == "" {
		Config.AddConfigPath("../app/config")
	}
	if os.Getenv("GO_ENVIRONMENT") == "" || os.Getenv("GO_ENVIRONMENT") == LOCAL {
		Config.SetConfigName(LOCAL)
		Config.Set("aws_credentials", credentials.NewStaticCredentials("foo", "bar", ""))
	}
	if environment, ok := os.LookupEnv("GO_ENVIRONMENT"); !ok {
		currentEnv = LOCAL
	} else {
		currentEnv = environment
	}
	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Print(Config.GetString("engines.AWS.sns"))
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
