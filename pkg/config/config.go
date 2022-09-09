// Copyright 2022 TCDZENGIN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type EnvironmentConfiguration struct {
	Env string `default:"Development"`
}

type DatabaseConfiguration struct {
	Name     string `default:"upwork"`
	Username string `default:"user"`
	Password string `default:"password"`
	Host     string `default:"localhost"`
	Port     string `default:"5432"`
}

type Configuration struct {
	Environment EnvironmentConfiguration
	Database    DatabaseConfiguration
}

const (
	configFileType      string = "yaml"
	configName          string = "config"
	configFileExtension string = ".yml"
	productionEnv       string = "Production"
)

var Config *Configuration = New("")

func New(path string) *Configuration {

	if path == "" {
		path = "./"
	}

	viper.SetConfigType(configFileType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(path)

	viper.BindEnv("environment.env", "B_ENV")
	viper.SetDefault("environment.env", "Development")

	viper.BindEnv("database.name", "DB_NAME")
	viper.SetDefault("database.name", "upwork")
	viper.BindEnv("database.username", "DB_USERNAME")
	viper.SetDefault("database.username", "postgres")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.SetDefault("database.password", "password")
	viper.BindEnv("database.host", "DB_HOST")
	viper.SetDefault("database.host", "localhost")
	viper.BindEnv("database.port", "DB_PORT")
	viper.SetDefault("database.port", 5432)

	configFilePath := filepath.Join(path, configName) + configFileExtension
	if err := readConfiguration(configFilePath); err != nil {
		return nil
	}

	viper.AutomaticEnv()
	var cfg *Configuration
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		if cfg.Environment.Env != productionEnv {
			log.Println("Config file changed:", e.Name)
		}
	})

	viper.WatchConfig()

	return cfg
}

// read configuration from file
func readConfiguration(filePath string) error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// if file does not exist, simply create one
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			os.Create(filePath)
		} else {
			return err
		}
		// let's write defaults
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}
