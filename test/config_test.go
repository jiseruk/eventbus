package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenance/wequeue-management_api/app/config"
)

func setup() {
	config.Init()
	defer os.Setenv("GO_ENVIRONMENT", "")

}
func TestConfigOptions(t *testing.T) {

	t.Run("It should set develop as default environment", func(t *testing.T) {
		setup()
		assert.Equal(t, config.DEVELOP, *config.GetCurrentEnvironment())
	})

	t.Run("It should set environment from GO_ENVIRONMENT variable if it's defined", func(t *testing.T) {
		os.Setenv("GO_ENVIRONMENT", config.LOCAL)
		setup()
		assert.Equal(t, config.LOCAL, *config.GetCurrentEnvironment())
	})

	t.Run("It should return the config value for the default env if GO_ENVIRONMENT is unset", func(t *testing.T) {
		setup()
		assert.Equal(t, "arn:aws:iam::719849485599:role/AmazonECSTask-wequeue-service-dev",
			config.Get("engines.AWS.lambda.executionRole"))
	})

	t.Run("Env variable should override yaml config", func(t *testing.T) {
		os.Setenv("ENGINES_AWS_LAMBDA_EXECUTIONROLE", "arn:env:role")
		setup()
		assert.Equal(t, "arn:env:role", config.Get("engines.AWS.lambda.executionRole"))
	})

	t.Run("it should return []string if the key has multiple values", func(t *testing.T) {
		setup()
		assert.Equal(t, []string{"subnet-0776e02a4f48a376d", "subnet-056133ad47ef8e2f6", "subnet-0a0cc57e1b32a6439"},
			config.Config.GetStringSlice("engines.AWS.lambda.subnetIds"))
	})
}
