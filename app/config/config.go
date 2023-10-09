package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Config represents the composition of yml settings.
type Config struct {
	Environment string

	Circle struct {
		MaxAmountPerUser int64
		MaxVoters        int
	}

	Aws struct {
		Auth struct {
			ClientId         string
			UserPoolId       string
			AwsDefaultRegion string
			ClientSecret     string
		}
	}

	Hosts struct {
		Vec string
		Svc struct {
			Sso string
		}
	}

	Db struct {
		Host      string
		Port      uint16
		Name      string
		User      string
		Password  string
		Migration bool
		Test      struct {
			Host     string
			Port     uint16
			Name     string
			User     string
			Password string
		}
	}

	Redis struct {
		Host     string
		Port     uint16
		Username string
		Db       uint16
		Timeout  uint16
		Password string
	}

	Security struct {
		Cors struct {
			Origins []string
		}
		Secrets struct {
			Key string
		}
	}
}

const (
	EnvironmentDev   = "development"
	EnvironmentProd  = "production"
	defaultFileName  = "config.service"
	secretFileName   = "secret.service"
	overrideFileName = "config.service.override"
)

func NewConfig(configPath string) *Config {
	c := &Config{}
	c.load(configPath)
	return c
}

// Load the configuration.
// The loaded configuration depends on the set environment
// variable ENVIRONMENT. If this variable is not set,
// the configuration will be loaded as development.
// Please follow the convention of naming the configuration files.
func (c *Config) load(configPath string) {
	c.readDefaultConfig(configPath)
	c.readSecretConfig(configPath)
	c.readOverrideConfig(configPath)
	c.checkEnvironment()
}

// readDefaultConfig reads the default configuration from the given
// config path. This configuration is required.
func (c *Config) readDefaultConfig(configPath string) {
	c.readConfig(configPath, defaultFileName)
}

// readSecretConfig reads the secret configuration from the given
// config path. This configuration is required.
func (c *Config) readSecretConfig(configPath string) {
	c.readConfig(configPath, secretFileName)
}

// readOverrideConfig reads the overwritten configuration from the given
// config path. This configuration is optional.
func (c *Config) readOverrideConfig(configPath string) {
	configDir := filepath.Dir(configPath)

	if _, err := os.Stat(configDir + "/" + overrideFileName + ".yml"); os.IsNotExist(err) {
		return
	}

	c.readConfig(configPath, overrideFileName)
}

// checkEnvironment against the set environment variable "ENVIRONMENT".
// If set, the environment will be set accordingly.
func (c *Config) checkEnvironment() {
	env := os.Getenv("ENVIRONMENT")

	if env == EnvironmentProd {
		c.Environment = EnvironmentProd
		return
	}

	c.Environment = EnvironmentDev
}

func (c *Config) readConfig(configPath string, configFileType string) {
	viperConfig := viper.New()

	viperConfig.SetConfigName(configFileType)
	viperConfig.SetConfigType("yml")
	viperConfig.AddConfigPath(filepath.Dir(configPath))

	if err := viperConfig.ReadInConfig(); err != nil {
		fmt.Printf("failed to read %s configuration. error: %s", configFileType, err)
		os.Exit(2)
	}

	err := viperConfig.Unmarshal(c)

	if err != nil {
		fmt.Printf("unable to decode Config. error: %s", err)
		os.Exit(2)
	}
}
