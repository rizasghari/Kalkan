package cfg

import (
	"fmt"
	"strings"

	"github.com/rizasghari/kalkan/internal/types"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server struct {
		Port string
	}
	Origins []types.Origin
	RL      struct {
		Enabled   bool
		Timeframe int
		Allowed   int
		Block     int
	}
	Redis struct {
		Url      string
		Password string
		DB       int
	}
}

var Config *Configuration

func NewConfiguration() (*Configuration, error) {
	viper.AddConfigPath("./internal/cfg")
	viper.AddConfigPath("$HOME")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config file: %s", err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	return Config, nil
}
