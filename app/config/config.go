package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/viper"
	"os"
	"strings"
)
const (
	LOCAL = "local"
	DEVELOP = "develop"
	STAGE = "stage"
	PRODUCTION = "production"
)
var Config *viper.Viper

func init(){
	Config = viper.New()
	Config.SetConfigType("yaml")
	Config.AddConfigPath("./app/config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if os.Getenv("GO_ENVIRONMENT") == "" {
		Config.AddConfigPath("../app/config")
	}
	if os.Getenv("GO_ENVIRONMENT") == "" || os.Getenv("GO_ENVIRONMENT") == LOCAL {
		Config.SetConfigName(LOCAL)
		Config.Set("aws_credentials", credentials.NewStaticCredentials("foo", "bar", ""))
	}

	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Print(Config.GetString("engines.AWS.sns"))
}

func Get(key string) string{
	return Config.GetString(key)
}

func GetObject(key string) interface{}{
	return Config.Get(key)
}

