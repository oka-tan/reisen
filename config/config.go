package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	PostgresConfig    PostgresConfig
	LnxConfig         LnxConfig
	Boards            []BoardConfig
	CspConfig         string
	TemplateConfig    TemplateConfig
	Hosts             []string
	Production        bool
	UseCatalogVariant bool
}

type BoardConfig struct {
	Name        string
	Description string
}

type PostgresConfig struct {
	ConnectionString string
}

type LnxConfig struct {
	Host string
	Port int
}

type TemplateConfig struct {
	ImagesUrl     string
	ThumbnailsUrl string
	OekakiUrl     string
	FaviconUrl    string
	BaseCssUrl    string
	JsUrl         string
	TegakiJsUrl   string
	TegakiCssUrl  string
	Themes        []Theme
	DefaultTheme  Theme
}

type Theme struct {
	Name string
	Url  string
}

func LoadConfig() Config {
	configFile := os.Getenv("REISEN_CONFIG")

	if configFile == "" {
		configFile = "./config.json"
	}

	f, err := os.Open(configFile)

	if err != nil {
		log.Fatalln(err)
	}

	var conf Config

	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalln(err)
	}

	return conf
}
