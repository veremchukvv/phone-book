package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr string
	Database   Database
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Schema   string
}

func LoadConfig(configFileName string) (*Config, error) {
	log := logrus.WithField("Method", "LoadConfig")
	v := viper.New()

	if configFileName != "" {
		v.SetConfigName(configFileName)
	}

	v.SetEnvPrefix("PHONE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AddConfigPath(".")

	err := v.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, err
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
		log.Println("No config detected, using default settings")
	}

	loadDefaultSettingsFor(v)
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func loadDefaultSettingsFor(v *viper.Viper) {
	// Порт, который слушает сервис
	v.SetDefault("ListenAddr", ":9090")

	// настройки подключения к бд
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", "54322")
	v.SetDefault("database.user", "phone")
	v.SetDefault("database.password", "phone")
	v.SetDefault("database.name", "phonedb")
	v.SetDefault("database.schema", "")
}

func (d *Database) ToDataSourceName() string {
	schema := d.Schema
	if schema == "" {
		schema = "public"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", "postgres", d.User, d.Password, d.Host, d.Port, d.Name, schema)
}
