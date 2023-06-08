package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	ModeWork struct {
		IsDebug string `yaml:"isDebug" env-description:"Debug mode" env-default:"yes"`
	} `yaml:"workMode"`
	Server struct {
		Host string `yaml:"host" env-description:"Server host" env-default:"localhost"`
		Port string `yaml:"port" env-description:"Server port" env-default:"8000"`
	} `yaml:"server"`
	UrlRepo struct {
		SitesFile string `yaml:"file" env-description:"File with URL" env-default:"../sites.txt"`
		Refresh   int    `yaml:"refresh" env-description:"Get refresh time (seconds)" env-default:"60"`
		Timeout   int    `yaml:"timeout" env-description:"Get request timeout (seconds)" env-default:"40"`
	} `yaml:"url_repo"`
}

var path = "../config.yml"

func NewConfig() Config {

	log.Println("\t\tRead application configuration...")
	config := &Config{}

	if err := cleanenv.ReadConfig(path, config); err != nil {
		help, _ := cleanenv.GetDescription(config, nil)
		log.Println(help)
		log.Fatal(fmt.Sprintf("%s", err))
	}
	log.Println("\t\tGet configuration - OK!")

	return *config
}
