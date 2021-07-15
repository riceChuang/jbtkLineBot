package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	BeautyUrl    string `yaml:"beautyurl"`
	MaxBeautyLen int32  `yaml:"maxbeautylen"`
	DcardUrl     string `yaml:"dcardurl"`
	MaxDcardLen  int32  `yaml:"maxdcardlen"`
	JokerUrl     string `yaml:"jokerurl"`
	MaxJokerLen  int32  `yaml:"maxjokerlen"`
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
