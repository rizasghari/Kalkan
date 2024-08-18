package cfg

import (
	"fmt"
	"strings"

	"github.com/rizasghari/kalkan/internal/models"
	"github.com/spf13/viper"
)

type configuration struct {
	Server struct {
		Host string
		Port string
	}
	Origins []models.Origin
}

var Config *configuration

func NewConfiguration() (*configuration, error) {
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./../configs")
	viper.AddConfigPath("./../../configs")
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
