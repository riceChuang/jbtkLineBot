package types

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	BeautyUrl string `yaml:"beautyurl"`
	DcardUrl  string `yaml:"dcardurl"`
	JokerUrl  string `yaml:"jokerurl"`
}

func New(fileName string) *Configuration {
	flag.Parse()

	c := Configuration{}

	//read and parse config file
	rootDirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("config: file error: %s", err.Error())
	}

	configPath := filepath.Join(rootDirPath, fileName)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		// config exists
		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("config: file error: %s", err.Error())
		}

		err = _yaml.Unmarshal(file, &c)
		if err != nil {
			log.Fatal("config: config error:", err)
		}
	}

	return &c
}
