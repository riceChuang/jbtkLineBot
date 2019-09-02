package types

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	BeautyUrl string `yaml:"beautyurl"`
	DcardUrl  string `yaml:"dcardurl"`
	JokerUrl  string `yaml:"jokerurl"`
}

var configData *Configuration

func InitialConfigPkg() {
	var config Configuration
	viper.SetConfigName("app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fail get config file: %s", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("error unmarshal : %s", err))
	}
	fmt.Println(config)

	configData = &config
}

func GetConfig() *Configuration {
	return configData
}
