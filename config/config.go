package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	PostgresConfig    PostgresConfig `toml:"postgres"`
	LnxConfig         LnxConfig      `toml:"lnx"`
	Boards            []BoardConfig  `toml:"boards"`
	CspConfig         string         `toml:"csp"`
	TemplateConfig    TemplateConfig `toml:"template"`
	UseCatalogVariant bool           `toml:"use_catalog_variant"`
	Port              int            `toml:"port"`
	ForceGzip         bool           `toml:"force_gzip"`
}

type BoardConfig struct {
	Name               string `toml:"name"`
	Description        string `toml:"description"`
	EnableLatex        bool   `toml:"enable_latex"`
	EnableTegaki       bool   `toml:"enable_tegaki"`
	EnableCountryFlags bool   `toml:"enable_country_flags"`
	EnablePolFlags     bool   `toml:"enable_pol_flags"`
	EnableMlpFlags     bool   `toml:"enable_mlp_flags"`
}

type PostgresConfig struct {
	ConnectionString string `toml:"connection_string"`
}

type LnxConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type TemplateConfig struct {
	ImagesUrl     string  `toml:"images_url"`
	ThumbnailsUrl string  `toml:"thumbnails_url"`
	OekakiUrl     string  `toml:"oekaki_url"`
	FaviconUrl    string  `toml:"favicon_url"`
	BaseCssUrl    string  `toml:"base_css_url"`
	JsUrl         string  `toml:"js_url"`
	TegakiJsUrl   string  `toml:"tegaki_js_url"`
	TegakiCssUrl  string  `toml:"tegaki_css_url"`
	FlagsJsUrl    string  `toml:"flags_js_url"`
	PolFlagsJsUrl string  `toml:"pol_flags_js_url"`
	MlpFlagsJsUrl string  `toml:"mlp_flags_js_url"`
	Themes        []Theme `toml:"themes"`
	DefaultTheme  Theme   `toml:"default_theme"`
}

type Theme struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

func LoadConfig() Config {
	configFile := os.Getenv("REISEN_CONFIG")

	if configFile == "" {
		configFile = "./config.toml"
	}

	f, err := os.Open(configFile)

	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	var conf Config

	if _, err := toml.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalln(err)
	}

	return conf
}

func (c *Config) IsLatexEnabled(board string) bool {
	for _, b := range c.Boards {
		if b.Name == board {
			return b.EnableLatex
		}
	}

	return false
}

func (c *Config) IsTegakiEnabled(board string) bool {
	for _, b := range c.Boards {
		if b.Name == board {
			return b.EnableTegaki
		}
	}

	return false
}

func (c *Config) AreCountryFlagsEnabled(board string) bool {
	for _, b := range c.Boards {
		if b.Name == board {
			return b.EnableCountryFlags
		}
	}

	return false
}

func (c *Config) ArePolFlagsEnabled(board string) bool {
	for _, b := range c.Boards {
		if b.Name == board {
			return b.EnablePolFlags
		}
	}

	return false
}

func (c *Config) AreMlpFlagsEnabled(board string) bool {
	for _, b := range c.Boards {
		if b.Name == board {
			return b.EnableMlpFlags
		}
	}

	return false
}
